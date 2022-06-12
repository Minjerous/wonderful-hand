package service

import (
	"context"
	"time"
	"wonderful-hand-common/cryptx"
	"wonderful-hand-common/jwt"
	"wonderful-hand-user/rpc/internal/dao"
	"wonderful-hand-user/rpc/internal/models"
	"wonderful-hand-user/rpc/user"
)

func (s *UserService) UserLogin(
	_ context.Context,
	req *user.UserLoginRegisterRequest,
) (resp *user.UserLoginRegisterResponse, _ error) {
	resp = new(user.UserLoginRegisterResponse)
	user := models.User{}
	err := dao.DB.Where("`name` = ?", req.GetUsername()).First(&user).Error
	if err != nil {
		resp.StatusCode = StatusBadLogin
		resp.StatusMsg = "query db failed"
		return
	}
	if cryptx.EncryptSHA256(req.GetPassword()) == user.Password {
		resp.UserId = user.UID
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
	err = passwordNotMatched
	resp.StatusCode = StatusBadLogin
	resp.StatusMsg = "incorrect password"
	return
}
