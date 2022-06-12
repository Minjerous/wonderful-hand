package router

import (
	"context"
	"wonderful-hand-common/rest/errdef"
	"wonderful-hand-user/api/router/helper"
)

type backendFunc[T helper.RequestModel, E helper.ResponseModel] func(ctx context.Context, req T) (resp E, err errdef.Err)

var (
	ub UnimplementedBackend
	_  backendFunc[*helper.RegisterLoginReq, helper.RegisterLoginResp] = ub.Register
	_  backendFunc[*helper.RegisterLoginReq, helper.RegisterLoginResp] = ub.Login
)

type Backend interface {
	Register(ctx context.Context, req *helper.RegisterLoginReq) (resp helper.RegisterLoginResp, err errdef.Err)
	Login(ctx context.Context, req *helper.RegisterLoginReq) (resp helper.RegisterLoginResp, err errdef.Err)
}

var _ Backend = (*UnimplementedBackend)(nil)

type UnimplementedBackend struct{}

func (u UnimplementedBackend) Register(_ context.Context, _ *helper.RegisterLoginReq) (_ helper.RegisterLoginResp, _ errdef.Err) {
	return
}

func (u UnimplementedBackend) Login(_ context.Context, _ *helper.RegisterLoginReq) (_ helper.RegisterLoginResp, _ errdef.Err) {
	return
}
