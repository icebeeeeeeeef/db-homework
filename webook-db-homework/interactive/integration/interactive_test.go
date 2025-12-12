package integration

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
	intrv1 "webook/api/proto/gen/intr/v1"
	"webook/interactive/grpc"
	"webook/interactive/integration/startup"
	"webook/interactive/repository"
	"webook/interactive/repository/cache"
	"webook/interactive/repository/dao"
	"webook/interactive/service"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	grpcLib "google.golang.org/grpc"
	"gorm.io/gorm"
)

type InteractiveTestSuite struct {
	suite.Suite
	db     *gorm.DB
	svc    service.InteractiveService
	repo   repository.InteractiveRepository
	server *grpcLib.Server
	client intrv1.IntrServiceClient
	lis    net.Listener
}

// 一个简单的内存假实现，替代 Redis 缓存，避免集成测试依赖外部 Redis
// 使用真实 Redis 缓存进行集成验证

func TestInteractive(t *testing.T) {
	suite.Run(t, new(InteractiveTestSuite))
}

func (s *InteractiveTestSuite) SetupSuite() {
	// 设置数据库配置
	viper.Set("db.mysql.dsn", "root:123123@tcp(172.28.240.236:13316)/webook_db")

	// 设置Redis配置
	viper.Set("db.redis.addr", "172.28.240.236:6379")

	s.db = startup.InitDB(startup.InitLogger())
	// 保障表结构
	_ = s.db.AutoMigrate(
		&dao.Interactive{},
		&dao.UserLikeSomething{},
		&dao.Collection{},
		&dao.UserCollectSomething{},
	)

	// 组装 repo 与 svc（使用真实 Redis 缓存）
	d := dao.NewInteractiveDAO(s.db)
	rdb := startup.InitRedis()
	c := cache.NewInteractiveCache(rdb)
	s.repo = repository.NewInteractiveRepository(d, c)

	s.svc = startup.InitInteractiveService()

	// 启动gRPC服务器
	s.startGRPCServer()
}

func (s *InteractiveTestSuite) startGRPCServer() {
	// 创建监听器
	var err error
	s.lis, err = net.Listen("tcp", ":0") // 使用随机端口
	assert.NoError(s.T(), err)

	// 创建gRPC服务器
	s.server = grpcLib.NewServer()

	// 注册服务
	grpcServer := grpc.NewInteractiveServiceServer(s.svc)
	intrv1.RegisterIntrServiceServer(s.server, grpcServer)

	// 启动服务器
	serverReady := make(chan bool)
	go func() {
		serverReady <- true
		_ = s.server.Serve(s.lis)
	}()

	// 等待服务器启动
	<-serverReady

	// 创建客户端连接
	conn, err := grpcLib.Dial(s.lis.Addr().String(), grpcLib.WithInsecure())
	assert.NoError(s.T(), err)
	s.client = intrv1.NewIntrServiceClient(conn)
}

func (s *InteractiveTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Stop()
	}
	if s.lis != nil {
		s.lis.Close()
	}
}

func (s *InteractiveTestSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// 清理数据库
	err := s.db.Exec("TRUNCATE TABLE `interactives`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `user_like_somethings`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `user_collect_somethings`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `collections`").Error
	assert.NoError(s.T(), err)

	// 清理 Redis 缓存
	rdb := startup.InitRedis()
	err = rdb.FlushDB(ctx).Err()
	assert.NoError(s.T(), err)
}

func (s *InteractiveTestSuite) TestIncrRead() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	biz := "article"
	id := int64(10001)

	// 测试gRPC IncrReadIfPresent
	incrReq := &intrv1.IncrReadIfPresentRequest{
		Biz: biz,
		Id:  id,
	}
	incrResp, err := s.client.IncrReadIfPresent(ctx, incrReq)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), incrResp)
	assert.True(s.T(), incrResp.Success)

	// 验证数据库中的数据
	var inter dao.Interactive
	err = s.db.WithContext(ctx).Where("biz = ? AND biz_id = ?", biz, id).First(&inter).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), inter.Readcnt)
	assert.Equal(s.T(), biz, inter.Biz)
	assert.Equal(s.T(), id, inter.BizId)
}

func (s *InteractiveTestSuite) TestLikeAndCancelLike() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	biz := "article"
	id := int64(10001)
	uid := int64(123)

	// 点赞
	likeReq := &intrv1.LikeRequest{
		Biz: biz,
		Id:  id,
		Uid: uid,
	}
	likeResp, err := s.client.Like(ctx, likeReq)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), likeResp)
	assert.True(s.T(), likeResp.Success)

	// 验证数据库中的点赞数据
	var inter dao.Interactive
	err = s.db.WithContext(ctx).Where("biz = ? AND biz_id = ?", biz, id).First(&inter).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), inter.Likecnt)
	assert.Equal(s.T(), biz, inter.Biz)
	assert.Equal(s.T(), id, inter.BizId)

	// 验证用户点赞记录
	var likeRecord dao.UserLikeSomething
	err = s.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).First(&likeRecord).Error
	assert.NoError(s.T(), err)
	assert.True(s.T(), likeRecord.Status)

	// 取消点赞
	cancelLikeReq := &intrv1.CancelLikeRequest{
		Biz: biz,
		Id:  id,
		Uid: uid,
	}
	cancelLikeResp, err := s.client.CancelLike(ctx, cancelLikeReq)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cancelLikeResp)
	assert.True(s.T(), cancelLikeResp.Success)

	// 验证取消点赞后的数据
	inter = dao.Interactive{}
	err = s.db.WithContext(ctx).Where("biz = ? AND biz_id = ?", biz, id).First(&inter).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(0), inter.Likecnt)

	// 验证用户点赞记录状态
	likeRecord = dao.UserLikeSomething{}
	err = s.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).First(&likeRecord).Error
	assert.NoError(s.T(), err)
	assert.False(s.T(), likeRecord.Status)
}

func (s *InteractiveTestSuite) TestCollect() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	biz := "article"
	id := int64(20001)
	uid := int64(321)

	// 执行收藏
	collectReq := &intrv1.CollectRequest{
		Biz: biz,
		Id:  id,
		Uid: uid,
	}
	collectResp, err := s.client.Collect(ctx, collectReq)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), collectResp)
	assert.True(s.T(), collectResp.Success)

	// 验证数据库中的收藏数据
	var inter dao.Interactive
	err = s.db.WithContext(ctx).Where("biz = ? AND biz_id = ?", biz, id).First(&inter).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), inter.Collectcnt)
	assert.Equal(s.T(), biz, inter.Biz)
	assert.Equal(s.T(), id, inter.BizId)

	// 验证用户收藏记录
	var collectRecord dao.UserCollectSomething
	err = s.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).First(&collectRecord).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), biz, collectRecord.Biz)
	assert.Equal(s.T(), id, collectRecord.BizId)
	assert.Equal(s.T(), uid, collectRecord.UID)
}

func (s *InteractiveTestSuite) TestGetAndStatusFlags() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	biz := "article"
	id := int64(30001)
	uid := int64(777)

	// 预置数据
	err := s.db.WithContext(ctx).Create(&dao.Interactive{
		BizId:      id,
		Biz:        biz,
		Readcnt:    5,
		Likecnt:    2,
		Collectcnt: 1,
	}).Error
	assert.NoError(s.T(), err)

	// 预置点赞状态
	err = s.db.WithContext(ctx).Create(&dao.UserLikeSomething{
		BizId:  id,
		Biz:    biz,
		UID:    uid,
		Status: true,
	}).Error
	assert.NoError(s.T(), err)

	// 通过 gRPC Get 获取信息
	getReq := &intrv1.GetRequest{
		Biz: biz,
		Id:  id,
		Uid: uid,
	}
	resp, err := s.client.Get(ctx, getReq)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotNil(s.T(), resp.Interactive)

	got := resp.Interactive
	assert.Equal(s.T(), biz, got.Biz)
	assert.Equal(s.T(), id, got.BizId)
	assert.Equal(s.T(), int64(5), got.Readcnt)
	assert.Equal(s.T(), int64(2), got.Likecnt)
	assert.Equal(s.T(), int64(1), got.Collectcnt)
	assert.Equal(s.T(), true, got.Liked)
	assert.Equal(s.T(), false, got.Collected)
}

func (s *InteractiveTestSuite) TestGetByIdsAndBatchIncRead() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	biz := "article"
	ids := []int64{40001, 40002}

	// 预置两条互动数据
	for _, id := range ids {
		err := s.db.WithContext(ctx).Create(&dao.Interactive{
			BizId:      id,
			Biz:        biz,
			Readcnt:    0,
			Likecnt:    0,
			Collectcnt: 0,
		}).Error
		assert.NoError(s.T(), err)
	}

	// 测试批量增加阅读数
	err := s.repo.BatchIncRead(ctx, []string{biz, biz}, ids)
	assert.NoError(s.T(), err)

	// 验证数据库中的数据
	for _, id := range ids {
		var inter dao.Interactive
		err := s.db.WithContext(ctx).Where("biz = ? AND biz_id = ?", biz, id).First(&inter).Error
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), int64(1), inter.Readcnt)
	}

	// 验证数据库中的数据
	var count int64
	err = s.db.WithContext(ctx).Model(&dao.Interactive{}).Where("biz = ?", biz).Count(&count).Error
	assert.NoError(s.T(), err)
	fmt.Printf("数据库中biz=%s的记录数: %d\n", biz, count)

	// 测试gRPC GetByIds
	getByIdsReq := &intrv1.GetByIdsRequest{
		Biz: biz,
		Ids: ids,
	}
	resp, err := s.client.GetByIds(ctx, getByIdsReq)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotNil(s.T(), resp.Interactive)

	fmt.Printf("GetByIds返回的map长度: %d\n", len(resp.Interactive))
	for k, v := range resp.Interactive {
		fmt.Printf("Key: %d, Value: %+v\n", k, v)
	}

	assert.Len(s.T(), resp.Interactive, 2)

	// 验证每个Interactive对象
	for _, id := range ids {
		inter, exists := resp.Interactive[id]
		assert.True(s.T(), exists)
		assert.NotNil(s.T(), inter)
		assert.Equal(s.T(), biz, inter.Biz)
		assert.Equal(s.T(), id, inter.BizId)
		assert.Equal(s.T(), int64(1), inter.Readcnt)
	}
}
