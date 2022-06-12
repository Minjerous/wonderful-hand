package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wonderful-hand-common/rest/errdef"
)

var (
	_ RequestModel = (*RegisterLoginReq)(nil)
	_ RequestModel = (*TokenReq)(nil)
)

type RequestModel interface {
	Read(ctx *gin.Context) errdef.Err
}

type RegisterLoginReq struct {
	Username string
	Password string
}

func (r *RegisterLoginReq) check() errdef.Err {
	if len(r.Username) > 32 {
		return errdef.Errorf(http.StatusNotAcceptable,
			CodeUnacceptedParam,
			"unexpect username, length should be less than 32")
	}
	if len(r.Password) > 32 {
		return errdef.Errorf(http.StatusNotAcceptable,
			CodeUnacceptedParam,
			"unexpect password, length should be less than 32")
	}
	if len(r.Username)*len(r.Password) == 0 {
		return errdef.Errorf(http.StatusNotAcceptable,
			CodeUnacceptedParam,
			"expect username and password")
	}
	return errdef.Nil
}

func (r *RegisterLoginReq) Read(ctx *gin.Context) errdef.Err {
	r.Username = ctx.Query("username")
	r.Password = ctx.Query("password")
	return r.check()
}

type TokenReq struct {
	Token string // 这是发送请求的用户的 Token 不是 Uid 用户的 Token，用于返回 resp 中 is_follow
}

func (u *TokenReq) Read(ctx *gin.Context) errdef.Err {
	u.Token = ctx.Query("token")
	return errdef.Nil
}
