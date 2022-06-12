package service

import (
	"context"
	"wonderful-hand-common/jwt"
	"wonderful-hand-user/rpc/user"
)

func (s UserService) UserTokenVerify(
	_ context.Context,
	req *user.UserTokenVerifyRequest,
) (resp *user.UserTokenVerifyResponse, _ error) {
	defer func() {
		if err := recover(); err != nil {
			resp.StatusCode = StatusBadVerify
			resp.StatusMsg = "bad token"
		}
	}()
	resp = new(user.UserTokenVerifyResponse)
	err, token := s.tokenVer.VerifyWithDefaultKey(req.GetAccessToken())
	if err != nil {
		resp.StatusCode = StatusBadVerify
		resp.StatusMsg = "invalid jwt"
		return
	}
	if !token.Claims.(jwt.MapClaims)["is_access"].(bool) {
		resp.StatusCode = StatusBadVerify
		resp.StatusMsg = "access token cannot be refresh token"
	}
	return
}
