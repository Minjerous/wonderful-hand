package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wonderful-hand-common/rest/errdef"
)

func Write(ctx *gin.Context, resp any, err errdef.Err) {
	if !errdef.IsNil(err) {
		r := DefaultResp{}
		r.StatusCode = err.StatusCode
		r.StatusMsg = err.Description
		ctx.JSON(err.HttpCode, r)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func WriteErr(ctx *gin.Context, err errdef.Err) {
	resp := DefaultResp{}
	resp.StatusCode = err.StatusCode
	resp.StatusMsg = err.Description
	ctx.JSON(err.HttpCode, resp)
	ctx.Abort()
	return
}
