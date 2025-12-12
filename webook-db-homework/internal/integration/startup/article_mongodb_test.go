//go:build e2e

package startup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"webook/internal/domain"
	"webook/internal/web"

	ijwt "webook/internal/web/jwt"

	articledao "webook/internal/repository/dao/article"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleHandlerTestSuite struct {
	suite.Suite
	server   *gin.Engine
	client   *mongo.Client
	database *mongo.Database
	col      *mongo.Collection
	livecol  *mongo.Collection
	node     *snowflake.Node
}

func TestArticle(t *testing.T) {
	suite.Run(t, new(ArticleHandlerTestSuite))
}

// InitTestMongoDB 初始化测试用的 MongoDB 连接
func InitTestMongoDB() (*mongo.Client, *mongo.Database, *mongo.Collection, *mongo.Collection, *snowflake.Node, error) {
	// 连接 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 尝试多种连接方式
	connectionStrings := []string{
		"mongodb://root:123123@localhost:27017", // 根据容器配置
		"mongodb://localhost:27017",
		"mongodb://admin:123456@localhost:27017",
		"mongodb://root:123456@localhost:27017",
		"mongodb://admin:admin@localhost:27017",
		"mongodb://root:root@localhost:27017",
	}

	var client *mongo.Client
	var err error

	for _, uri := range connectionStrings {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err == nil {
			// 测试连接
			err = client.Ping(ctx, nil)
			if err == nil {
				fmt.Printf("MongoDB 连接成功，使用连接字符串: %s\n", uri)
				break
			}
		}
	}
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// 获取数据库和集合
	database := client.Database("webook_test")
	col := database.Collection("articles")               // 制作库
	livecol := database.Collection("published_articles") // 线上库

	// 初始化雪花算法节点
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return client, database, col, livecol, node, nil
}

func (s *ArticleHandlerTestSuite) SetupSuite() {
	s.server = gin.Default()

	// 初始化 MongoDB 连接
	var err error
	s.client, s.database, s.col, s.livecol, s.node, err = InitTestMongoDB()
	if err != nil {
		s.T().Skipf("MongoDB 连接失败，跳过所有测试: %v", err)
		return
	}

	s.server.Use(func(context *gin.Context) {
		// 直接设置好
		context.Set("claims", &ijwt.UserClaims{
			UserId: 123,
		})
		context.Next()
	})
	/*
	   // 手动初始化 ArticleHandler，使用 MongoDB DAO
	   articleDAO := articledao.NewMongoDBArticleDAO(s.client, s.node, s.database)
	   articleRepository := articlerepo.NewCachedArticleRepository(articleDAO, nil, nil, nil) // MongoDB 不需要其他 DAO
	   logger := InitLogger()                                                                 // 使用本地的 InitLogger
	   articleService := service.NewArticleService(articleRepository, logger)
	   articleHandler := web.NewArticleHandler(articleService, logger)
	   articleHandler.RegisterRoutes(s.server)
	*/
}

func (s *ArticleHandlerTestSuite) TearDownTest() {
	// 清理 MongoDB 集合
	if s.col == nil || s.livecol == nil {
		return // MongoDB 连接失败，跳过清理
	}

	ctx := context.Background()

	// 删除所有文档
	_, err := s.col.DeleteMany(ctx, bson.M{})
	if err != nil {
		s.T().Logf("清理 articles 集合失败: %v", err)
	}

	_, err = s.livecol.DeleteMany(ctx, bson.M{})
	if err != nil {
		s.T().Logf("清理 live_articles 集合失败: %v", err)
	}
}

func (s *ArticleHandlerTestSuite) TestArticle() {
	s.T().Log("test article")
}

func (s *ArticleHandlerTestSuite) TestMongoDBConnection() {
	t := s.T()

	// 测试 MongoDB 连接
	ctx := context.Background()

	// 插入测试数据
	testArticle := articledao.Article{
		ID:        999,
		Title:     "测试标题",
		Content:   "测试内容",
		AuthorID:  123,
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
		Status:    0,
	}

	_, err := s.col.InsertOne(ctx, testArticle)
	if err != nil {
		t.Skipf("MongoDB 连接失败，跳过测试: %v", err)
		return
	}

	// 查询测试数据
	filter := bson.M{"id": 999}
	var result articledao.Article
	err = s.col.FindOne(ctx, filter).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, "测试标题", result.Title)
	assert.Equal(t, "测试内容", result.Content)
	assert.Equal(t, int64(123), result.AuthorID)

	// 清理测试数据
	_, err = s.col.DeleteOne(ctx, filter)
	assert.NoError(t, err)

	t.Log("MongoDB 连接和操作测试通过")
}

func (s *ArticleHandlerTestSuite) TestArticleHandler_Edit() {
	t := s.T()

	testCases := []struct {
		name string
		// 要提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after func(t *testing.T)
		// 构造请求，直接使用 req
		// 也就是说，我们放弃测试 Bind 的异常分支
		req Article

		// 预期响应
		wantCode   int
		wantResult web.Result[int64]
	}{
		{
			name: "新建帖子",
			before: func(t *testing.T) {
				// 什么也不需要做
			},
			after: func(t *testing.T) {
				// 验证一下数据
				ctx := context.Background()
				filter := bson.M{"author_id": 123}
				var art articledao.Article
				err := s.col.FindOne(ctx, filter).Decode(&art)
				assert.NoError(t, err)
				assert.True(t, art.CreatedAt > 0)
				assert.True(t, art.UpdatedAt > 0)
				assert.True(t, art.ID > 0)
				art.ID = 0
				// 重置了这些值，因为无法比较
				art.UpdatedAt = 0
				art.CreatedAt = 0
				assert.Equal(t, articledao.Article{
					ID:       0,
					Title:    "hello，你好",
					Content:  "随便试试",
					AuthorID: 123,
					Status:   uint8(domain.ArticleStatusDraft), // 新建的文章状态应该是 0
				}, art)
			},
			req: Article{
				Title:   "hello，你好",
				Content: "随便试试",
			},
			wantCode: 200,
			wantResult: web.Result[int64]{
				Code: 0,
				Msg:  "编辑成功",
				Data: 1,
			},
		},
		{
			// 这个是已经有了，然后修改之后再保存
			name: "更新帖子",
			before: func(t *testing.T) {
				// 模拟已经存在的帖子，并且是已经发布的帖子
				ctx := context.Background()
				article := articledao.Article{
					ID:        2,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					AuthorID:  123,
					Status:    uint8(domain.ArticleStatusDraft), // domain.ArticleStatusPublished
				}
				_, err := s.col.InsertOne(ctx, article)
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证一下数据
				ctx := context.Background()
				filter := bson.M{"id": 2}
				var art articledao.Article
				err := s.col.FindOne(ctx, filter).Decode(&art)
				assert.NoError(t, err)
				assert.True(t, art.UpdatedAt > 234)
				art.UpdatedAt = 0
				assert.Equal(t, articledao.Article{
					ID:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorID: 123,
					// 创建时间没变
					CreatedAt: 456,
					Status:    uint8(domain.ArticleStatusDraft), // 更新后状态应该是 2
				}, art)
			},
			req: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 200,
			wantResult: web.Result[int64]{
				Code: 0,
				Msg:  "编辑成功",
				Data: 2,
			},
		},

		{
			name: "更新别人的帖子",
			before: func(t *testing.T) {
				// 模拟已经存在的帖子
				ctx := context.Background()
				article := articledao.Article{
					ID:        3,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					// 注意。这个 AuthorID 我们设置为另外一个人的ID
					AuthorID: 789,
					Status:   uint8(domain.ArticleStatusDraft), // domain.ArticleStatusPublished
				}
				_, err := s.col.InsertOne(ctx, article)
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 更新应该是失败了，数据没有发生变化
				ctx := context.Background()
				filter := bson.M{"id": 3}
				var art articledao.Article
				err := s.col.FindOne(ctx, filter).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, articledao.Article{
					ID:        3,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					AuthorID:  789,
					Status:    uint8(domain.ArticleStatusDraft), // domain.ArticleStatusPublished
				}, art)
			},
			req: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 200,
			wantResult: web.Result[int64]{
				Code: 500,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			// 不能有 error
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/edit", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type",
				"application/json")
			recorder := httptest.NewRecorder()

			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			// 利用泛型来限定结果必须是 int64
			var result web.Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}

}

func (s *ArticleHandlerTestSuite) TestArticle_Publish() {
	t := s.T()

	testCases := []struct {
		name string
		// 要提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after func(t *testing.T)
		req   Article

		// 预期响应
		wantCode   int
		wantResult web.Result[int64]
	}{
		{
			name: "新建帖子并发表",
			before: func(t *testing.T) {
				// 什么也不需要做
			},
			after: func(t *testing.T) {
				// 验证一下数据
				ctx := context.Background()
				filter := bson.M{"author_id": 123}
				var art articledao.Article
				err := s.col.FindOne(ctx, filter).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, "hello，你好", art.Title)
				assert.Equal(t, "随便试试", art.Content)
				assert.Equal(t, int64(123), art.AuthorID)
				assert.True(t, art.CreatedAt > 0)
				assert.True(t, art.UpdatedAt > 0)
				var publishedArt articledao.ReaderArticle
				// 对于 ReaderArticle，需要使用嵌套查询条件
				readerFilter := bson.M{"author_id": 123}
				t.Logf("查询条件: %+v", readerFilter)
				err = s.livecol.FindOne(ctx, readerFilter).Decode(&publishedArt)
				if err != nil {
					t.Logf("查询线上库失败: %v", err)
					// 尝试查询所有文档看看是否有数据
					cursor, err := s.livecol.Find(ctx, bson.M{})
					if err == nil {
						var allDocs []bson.M
						cursor.All(ctx, &allDocs)
						t.Logf("线上库中的所有文档: %+v", allDocs)
					}
					// 尝试不同的查询条件
					t.Logf("尝试查询 author_id: 123")
					cursor2, err := s.livecol.Find(ctx, bson.M{"author_id": 123})
					if err == nil {
						var docs []bson.M
						cursor2.All(ctx, &docs)
						t.Logf("使用 author_id 查询结果: %+v", docs)
					}
				}
				assert.NoError(t, err)
				t.Logf("查询到的 publishedArt: %+v", publishedArt)
				assert.Equal(t, "hello，你好", publishedArt.Title)
				assert.Equal(t, "随便试试", publishedArt.Content)
				assert.Equal(t, int64(123), publishedArt.AuthorID)
				assert.True(t, publishedArt.CreatedAt > 0)
				assert.True(t, publishedArt.UpdatedAt > 0)
			},
			req: Article{
				Title:   "hello，你好",
				Content: "随便试试",
			},
			wantCode: 200,
			wantResult: web.Result[int64]{
				Code: 0,
				Msg:  "发表成功",
				Data: 1,
			},
		},

		{
			// 制作库有，但是线上库没有
			name: "更新帖子并新发表",
			before: func(t *testing.T) {
				// 模拟已经存在的帖子
				ctx := context.Background()
				article := articledao.Article{
					ID:        2,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					AuthorID:  123,
				}
				_, err := s.col.InsertOne(ctx, article)
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证一下数据
				ctx := context.Background()
				filter := bson.M{"id": 2}
				var art articledao.Article
				err := s.col.FindOne(ctx, filter).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, "新的标题", art.Title)
				assert.Equal(t, "新的内容", art.Content)
				assert.Equal(t, int64(123), art.AuthorID)
				// 创建时间没变
				assert.Equal(t, int64(456), art.CreatedAt)
				// 更新时间变了
				assert.True(t, art.UpdatedAt > 234)
				var publishedArt articledao.ReaderArticle
				readerFilter := bson.M{"id": 2}
				err = s.livecol.FindOne(ctx, readerFilter).Decode(&publishedArt)
				assert.NoError(t, err)
				assert.Equal(t, "新的标题", publishedArt.Title)
				assert.Equal(t, "新的内容", publishedArt.Content)
				assert.Equal(t, int64(123), publishedArt.AuthorID)
				assert.True(t, publishedArt.CreatedAt > 0)
				assert.True(t, publishedArt.UpdatedAt > 0)
			},
			req: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 200,
			wantResult: web.Result[int64]{
				Code: 0,
				Msg:  "发表成功",
				Data: 2,
			},
		},
		{
			name: "更新帖子，并且重新发表",
			before: func(t *testing.T) {
				ctx := context.Background()
				art := articledao.Article{
					ID:        3,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					AuthorID:  123,
				}
				_, err := s.col.InsertOne(ctx, art)
				assert.NoError(t, err)
				_, err = s.livecol.InsertOne(ctx, articledao.ReaderArticle(art))
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx := context.Background()
				filter := bson.M{"id": 3}
				var art articledao.Article
				err := s.col.FindOne(ctx, filter).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, "新的标题", art.Title)
				assert.Equal(t, "新的内容", art.Content)
				assert.Equal(t, int64(123), art.AuthorID)
				// 创建时间没变
				assert.Equal(t, int64(456), art.CreatedAt)
				// 更新时间变了
				assert.True(t, art.UpdatedAt > 234)

				var part articledao.ReaderArticle
				readerFilter := bson.M{"id": 3}
				err = s.livecol.FindOne(ctx, readerFilter).Decode(&part)
				assert.NoError(t, err)
				assert.Equal(t, "新的标题", part.Title)
				assert.Equal(t, "新的内容", part.Content)
				assert.Equal(t, int64(123), part.AuthorID)
				// 创建时间没变
				assert.Equal(t, int64(456), part.CreatedAt)
				// 更新时间变了
				assert.True(t, part.UpdatedAt > 234)
			},
			req: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 200,
			wantResult: web.Result[int64]{
				Code: 0,
				Msg:  "发表成功",
				Data: 3,
			},
		},
		{
			name: "更新别人的帖子，并且发表失败",
			before: func(t *testing.T) {
				ctx := context.Background()
				art := articledao.Article{
					ID:        4,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					// 注意。这个 AuthorID 我们设置为另外一个人的ID
					AuthorID: 789,
				}
				_, err := s.col.InsertOne(ctx, art)
				assert.NoError(t, err)
				_, err = s.livecol.InsertOne(ctx, articledao.ReaderArticle{
					ID:        4,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					AuthorID:  789,
				})
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 更新应该是失败了，数据没有发生变化
				ctx := context.Background()
				filter := bson.M{"id": 4}
				var art articledao.Article
				err := s.col.FindOne(ctx, filter).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, "我的标题", art.Title)
				assert.Equal(t, "我的内容", art.Content)
				assert.Equal(t, int64(456), art.CreatedAt)
				assert.Equal(t, int64(234), art.UpdatedAt)
				assert.Equal(t, int64(789), art.AuthorID)

				var part articledao.ReaderArticle
				// 数据没有变化
				readerFilter := bson.M{"id": 4}
				err = s.livecol.FindOne(ctx, readerFilter).Decode(&part)
				assert.NoError(t, err)
				assert.Equal(t, "我的标题", part.Title)
				assert.Equal(t, "我的内容", part.Content)
				assert.Equal(t, int64(789), part.AuthorID)
				// 创建时间没变
				assert.Equal(t, int64(456), part.CreatedAt)
				// 更新时间变了
				assert.Equal(t, int64(234), part.UpdatedAt)
			},
			req: Article{
				Id:      4,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: 200,
			wantResult: web.Result[int64]{
				Code: 500,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			// 不能有 error
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/publish", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type",
				"application/json")
			recorder := httptest.NewRecorder()

			s.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			// 反序列化为结果
			// 利用泛型来限定结果必须是 int64
			var result web.Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
