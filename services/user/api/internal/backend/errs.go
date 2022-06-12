package backend

import (
	"net/http"
	"wonderful-hand-common/rest/errdef"
	"wonderful-hand-user/api/router/helper"
)

var (
	TimeoutErr  = errdef.New(http.StatusRequestTimeout, helper.CodeTimeout, "timeout")
	InternalErr = errdef.New(http.StatusInternalServerError, helper.CodeInternalErr, "internal error")
)
