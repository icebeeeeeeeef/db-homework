package integration

import (
	"context"
	"net"
	"testing"
	intrv1 "webook/api/proto/gen/intr/v1"
	"webook/interactive/domain"
	"webook/interactive/grpc"

	"github.com/stretchr/testify/assert"
	grpcLib "google.golang.org/grpc"
)

// 简单的gRPC测试，不依赖外部服务
func TestGRPCServerBasic(t *testing.T) {
	// 创建模拟的service
	mockSvc := &mockInteractiveService{}

	// 创建监听器
	lis, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	defer lis.Close()

	// 创建gRPC服务器
	server := grpcLib.NewServer()
	defer server.Stop()

	// 注册服务
	grpcServer := grpc.NewInteractiveServiceServer(mockSvc)
	intrv1.RegisterIntrServiceServer(server, grpcServer)

	// 启动服务器
	serverReady := make(chan bool)
	go func() {
		serverReady <- true
		_ = server.Serve(lis)
	}()

	// 等待服务器启动
	<-serverReady

	// 创建客户端连接
	conn, err := grpcLib.Dial(lis.Addr().String(), grpcLib.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := intrv1.NewIntrServiceClient(conn)

	// 测试Like请求
	likeReq := &intrv1.LikeRequest{
		Biz: "article",
		Id:  123,
		Uid: 456,
	}
	likeResp, err := client.Like(context.Background(), likeReq)
	assert.NoError(t, err)
	assert.NotNil(t, likeResp)
	assert.True(t, likeResp.Success)

	// 验证请求格式
	assert.Equal(t, "article", likeReq.Biz)
	assert.Equal(t, int64(123), likeReq.Id)
	assert.Equal(t, int64(456), likeReq.Uid)

	// 验证响应格式
	assert.True(t, likeResp.Success)
}

// 模拟的InteractiveService实现
type mockInteractiveService struct{}

func (m *mockInteractiveService) Like(ctx context.Context, biz string, id, uid int64) error {
	return nil
}

func (m *mockInteractiveService) CancelLike(ctx context.Context, biz string, id, uid int64) error {
	return nil
}

func (m *mockInteractiveService) Collect(ctx context.Context, biz string, id, uid int64) error {
	return nil
}

func (m *mockInteractiveService) Get(ctx context.Context, biz string, id, uid int64) (domain.Interactive, error) {
	return domain.Interactive{
		Biz:        biz,
		BizId:      id,
		Readcnt:    10,
		Likecnt:    5,
		Collectcnt: 3,
		Liked:      true,
		Collected:  false,
	}, nil
}

func (m *mockInteractiveService) IncrReadIfPresent(ctx context.Context, biz string, id int64) error {
	return nil
}

func (m *mockInteractiveService) GetByIds(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error) {
	result := make(map[int64]domain.Interactive)
	for _, id := range ids {
		result[id] = domain.Interactive{
			Biz:        biz,
			BizId:      id,
			Readcnt:    1,
			Likecnt:    0,
			Collectcnt: 0,
			Liked:      false,
			Collected:  false,
		}
	}
	return result, nil
}
