//go:build e2e && oldarticle

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/domain"
	"webook/internal/integration/startup"
	"webook/internal/service"
	"webook/internal/web"

	ijwt "webook/internal/web/jwt"

	articlerepo "webook/internal/repository/article"
	articledao "webook/internal/repository/dao/article"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ArticleHandlerTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func TestArticle(t *testing.T) {
	suite.Run(t, new(ArticleHandlerTestSuite))
}

func (s *ArticleHandlerTestSuite) SetupSuite() {
	s.server = gin.Default()
	s.db = startup.InitTestDB()
	s.server.Use(func(context *gin.Context) {
		// 直接设置好
		context.Set("claims", &ijwt.UserClaims{
			UserId: 123,
		})
		context.Next()
	})

	// 手动初始化 ArticleHandler，避免依赖 wire 的复杂依赖注入
	articleDAO := articledao.NewArticleDAO(s.db)
	articleAuthorDAO := articledao.NewArticleAuthorDAO(s.db)
	articleReaderDAO := articledao.NewArticleReaderDAO(s.db)
	articleRepository := articlerepo.NewCachedArticleRepository(articleDAO, articleAuthorDAO, articleReaderDAO, s.db)
	articleService := service.NewArticleService(articleRepository, startup.InitLogger())
	articleHandler := web.NewArticleHandler(articleService, startup.InitLogger())
	articleHandler.RegisterRoutes(s.server)
}

func (s *ArticleHandlerTestSuite) TearDownTest() {
	err := s.db.Exec("TRUNCATE TABLE `articles`").Error
	assert.NoError(s.T(), err)
	s.db.Exec("TRUNCATE TABLE `reader_articles`")
}

func (s *ArticleHandlerTestSuite) TestArticle() {
	s.T().Log("test article")
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
				var art articledao.Article
				s.db.Where("author_id = ?", 123).First(&art)
				assert.True(t, art.CreatedAt > 0)
				assert.True(t, art.UpdatedAt > 0)
				// 重置了这些值，因为无法比较
				art.UpdatedAt = 0
				art.CreatedAt = 0
				assert.Equal(t, articledao.Article{
					ID:       1,
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
				s.db.Create(&articledao.Article{
					ID:        2,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					AuthorID:  123,
					Status:    uint8(domain.ArticleStatusDraft), // domain.ArticleStatusPublished
				})
			},
			after: func(t *testing.T) {
				// 验证一下数据
				var art articledao.Article
				s.db.Where("id = ?", 2).First(&art)
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
				s.db.Create(&articledao.Article{
					ID:        3,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					// 注意。这个 AuthorID 我们设置为另外一个人的ID
					AuthorID: 789,
					Status:   uint8(domain.ArticleStatusDraft), // domain.ArticleStatusPublished
				})
			},
			after: func(t *testing.T) {
				// 更新应该是失败了，数据没有发生变化
				var art articledao.Article
				s.db.Where("id = ?", 3).First(&art)
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
				var art articledao.Article
				s.db.Where("author_id = ?", 123).First(&art)
				assert.Equal(t, "hello，你好", art.Title)
				assert.Equal(t, "随便试试", art.Content)
				assert.Equal(t, int64(123), art.AuthorID)
				assert.True(t, art.CreatedAt > 0)
				assert.True(t, art.UpdatedAt > 0)
				var publishedArt articledao.ReaderArticle
				s.db.Where("author_id = ?", 123).First(&publishedArt)
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
				s.db.Create(&articledao.Article{
					ID:        2,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					AuthorID:  123,
				})
			},
			after: func(t *testing.T) {
				// 验证一下数据
				var art articledao.Article
				s.db.Where("id = ?", 2).First(&art)
				assert.Equal(t, "新的标题", art.Title)
				assert.Equal(t, "新的内容", art.Content)
				assert.Equal(t, int64(123), art.AuthorID)
				// 创建时间没变
				assert.Equal(t, int64(456), art.CreatedAt)
				// 更新时间变了
				assert.True(t, art.UpdatedAt > 234)
				var publishedArt articledao.ReaderArticle
				s.db.Where("id = ?", 2).First(&publishedArt)
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
				art := articledao.Article{
					ID:        3,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					AuthorID:  123,
				}
				s.db.Create(&art)
				readerArt := articledao.ReaderArticle(art)
				s.db.Create(&readerArt)
			},
			after: func(t *testing.T) {
				var art articledao.Article
				s.db.Where("id = ?", 3).First(&art)
				assert.Equal(t, "新的标题", art.Title)
				assert.Equal(t, "新的内容", art.Content)
				assert.Equal(t, int64(123), art.AuthorID)
				// 创建时间没变
				assert.Equal(t, int64(456), art.CreatedAt)
				// 更新时间变了
				assert.True(t, art.UpdatedAt > 234)

				var part articledao.ReaderArticle
				s.db.Where("id = ?", 3).First(&part)
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
				art := articledao.Article{
					ID:        4,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					// 注意。这个 AuthorID 我们设置为另外一个人的ID
					AuthorID: 789,
				}
				s.db.Create(&art)
				s.db.Create(&articledao.ReaderArticle{
					ID:        4,
					Title:     "我的标题",
					Content:   "我的内容",
					CreatedAt: 456,
					UpdatedAt: 234,
					AuthorID:  789,
					Status:    uint8(domain.ArticleStatusDraft),
				})
			},

			after: func(t *testing.T) {
				// 更新应该是失败了，数据没有发生变化
				var art articledao.Article
				s.db.Where("id = ?", 4).First(&art)
				assert.Equal(t, "我的标题", art.Title)
				assert.Equal(t, "我的内容", art.Content)
				assert.Equal(t, int64(456), art.CreatedAt)
				assert.Equal(t, int64(234), art.UpdatedAt)
				assert.Equal(t, int64(789), art.AuthorID)

				var part articledao.ReaderArticle
				// 数据没有变化
				s.db.Where("id = ?", 4).First(&part)
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
