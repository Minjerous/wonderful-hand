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

func (s UserService) UserRegister(
	_ context.Context,
	req *user.UserLoginRegisterRequest,
) (resp *user.UserLoginRegisterResponse, _ error) {
	resp = new(user.UserLoginRegisterResponse)
	user := models.User{}
	err := dao.DB.Where("`name` = ?", req.GetUsername()).First(&user).Error
	if err == nil {
		resp.StatusCode = StatusBadRegister
		resp.StatusMsg = "account already created"
		return
	}

	user.Name = req.Username
	user.NickName = req.Username
	user.Password = cryptx.EncryptSHA256(req.Password)
	err = dao.DB.Create(&user).Error
	if err != nil {
		resp.StatusCode = StatusBadRegister
		resp.StatusMsg = "create failed"
		return
	}

	var uid int64
	err = dao.DB.Model(&user).Where("`name`=?", user.Name).Select("uid").Find(&uid).Error
	if err != nil {
		resp.StatusCode = StatusBadRegister
		resp.StatusMsg = "query db failed"
		return
	}

	resp.UserId = uid
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
