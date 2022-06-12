package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
	"wonderful-hand-common/rest/errdef"
	"wonderful-hand-user/api/router/helper"
)

func template[T helper.RequestModel, E helper.ResponseModel](
	ctx *gin.Context, req T,
	backendFunc func(ctx context.Context, req T) (resp E, err errdef.Err)) {
	if err := req.Read(ctx); err != errdef.Nil {
		helper.WriteErr(ctx, err)
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	resp, err := backendFunc(c, req)
	helper.Write(ctx, resp, err)
}

func (r *router) handleRegister(ctx *gin.Context) {
	req := &helper.RegisterLoginReq{}
	template(ctx, req, r.backend.Register)
}

func (r *router) handleLogin(ctx *gin.Context) {
	req := &helper.RegisterLoginReq{}
	template(ctx, req, r.backend.Login)
}
