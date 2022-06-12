package service

import (
	"context"
	"time"
	"wonderful-hand-common/jwt"
	"wonderful-hand-user/rpc/user"
)

func (s UserService) UserTokenRefresh(
	_ context.Context,
	req *user.UserTokenRefreshRequest,
) (resp *user.UserTokenRefreshResponse, _ error) {
	defer func() {
		if err := recover(); err != nil {
			resp.StatusCode = StatusBadVerify
			resp.StatusMsg = "bad token"
		}
	}()
	resp = new(user.UserTokenRefreshResponse)
	err, token := s.tokenVer.VerifyWithDefaultKey(req.GetRefreshToken())
	if err != nil {
		resp.StatusCode = StatusBadVerify
		resp.StatusMsg = "refresh token is invalid"
		return
	}
	if token.Claims.(jwt.MapClaims)["is_access"].(bool) {
		resp.StatusCode = StatusBadRefresh
		resp.StatusMsg = "refresh token cannot be access token"
		return
	}
	resp.AccessToken, err = s.tokenGen.Sign(tokenClaims{
		StandClaims: jwt.StandClaims{
			ExpiresAt: time.Now().Unix() + s.cfg.Auth.AccessExpire,
			IssuedAt:  time.Now().Unix(),
			Issuer:    s.cfg.Server.Name,
		},
		IsAccess: true,
	})

	if err != nil {
		resp.StatusCode = StatusBadLogin
		resp.StatusMsg = "sign token failed"
	}

	resp.RefreshToken, err = s.tokenGen.Sign(tokenClaims{
		StandClaims: jwt.StandClaims{
			ExpiresAt: time.Now().Unix() + (s.cfg.Auth.AccessExpire * 5),
			IssuedAt:  time.Now().Unix(),
			Issuer:    s.cfg.Server.Name,
		},
		IsAccess: false,
	})

	if err != nil {
		resp.StatusCode = StatusBadLogin
		resp.StatusMsg = "sign token failed"
	}

	return
}
