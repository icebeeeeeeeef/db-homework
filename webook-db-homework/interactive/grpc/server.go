package grpc

import (
	"context"
	intrv1 "webook/api/proto/gen/intr/v1"
	"webook/interactive/domain"
	"webook/interactive/service"

	"google.golang.org/grpc"
)

// 在这里把grpc包装成server
type InteractiveServiceServer struct {
	intrv1.UnimplementedIntrServiceServer
	svc service.InteractiveService
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{
		svc: svc,
	}
}

func (i *InteractiveServiceServer) Register(server *grpc.Server) {
	intrv1.RegisterIntrServiceServer(server, i)
}

func (i *InteractiveServiceServer) Like(ctx context.Context, req *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	err := i.svc.Like(ctx, req.GetBiz(), req.GetId(), req.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.LikeResponse{
		Success: true,
	}, nil
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, req *intrv1.CancelLikeRequest) (*intrv1.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx, req.GetBiz(), req.GetId(), req.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.CancelLikeResponse{
		Success: true,
	}, nil
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, req *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	err := i.svc.Collect(ctx, req.GetBiz(), req.GetId(), req.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.CollectResponse{
		Success: true,
	}, nil
}

func (i *InteractiveServiceServer) Get(ctx context.Context, req *intrv1.GetRequest) (*intrv1.GetResponse, error) {
	inter, err := i.svc.Get(ctx, req.GetBiz(), req.GetId(), req.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Interactive: i.toDTO(inter),
	}, nil
}

func (i *InteractiveServiceServer) IncrReadIfPresent(ctx context.Context, req *intrv1.IncrReadIfPresentRequest) (*intrv1.IncrReadIfPresentResponse, error) {
	err := i.svc.IncrReadIfPresent(ctx, req.GetBiz(), req.GetId())
	if err != nil {
		return nil, err
	}
	return &intrv1.IncrReadIfPresentResponse{
		Success: true,
	}, nil
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, req *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	inters, err := i.svc.GetByIds(ctx, req.GetBiz(), req.GetIds())
	if err != nil {
		return nil, err
	}
	mp := make(map[int64]*intrv1.Interactive)
	for _, inter := range inters {
		mp[inter.BizId] = i.toDTO(inter)
	}
	return &intrv1.GetByIdsResponse{
		Interactive: mp,
	}, nil
}

func (i *InteractiveServiceServer) toDTO(inter domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        inter.Biz,
		BizId:      inter.BizId,
		Liked:      inter.Liked,
		Collected:  inter.Collected,
		Readcnt:    inter.Readcnt,
		Likecnt:    inter.Likecnt,
		Collectcnt: inter.Collectcnt,
	}
}
