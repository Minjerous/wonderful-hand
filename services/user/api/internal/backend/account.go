package backend

import (
	"context"
	"wonderful-hand-common/rest/errdef"
	"wonderful-hand-user/api/router/helper"
	"wonderful-hand-user/rpc/user"
)

func (b *backend) Register(
	ctx context.Context,
	req *helper.RegisterLoginReq,
) (resp helper.RegisterLoginResp, err errdef.Err) {
	err = errdef.Nil
	rp, e := b.RpcClis.UserSrvCli.UserRegister(ctx, &user.UserLoginRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if e != nil {
		err = InternalErr
		if ctx.Err() != nil {
			err = TimeoutErr
		}
		return
	}
	if rp.StatusCode != 0 {
		resp.StatusCode = int(rp.StatusCode)
		resp.StatusMsg = rp.GetStatusMsg()
		return
	}
	resp.UserId = rp.UserId
	resp.AccessToken = rp.AccessToken
	resp.RefreshToken = rp.RefreshToken
	return
}
func (b *backend) Login(
	ctx context.Context,
	req *helper.RegisterLoginReq,
) (resp helper.RegisterLoginResp, err errdef.Err) {
	err = errdef.Nil
	rp, e := b.RpcClis.UserSrvCli.UserLogin(ctx, &user.UserLoginRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if e != nil {
		err = InternalErr
		if ctx.Err() != nil {
			err = TimeoutErr
		}
		return
	}
	if rp.StatusCode != 0 {
		err = errdef.Errorf(400, helper.CodeWrongParam, "login failed")
		return
	}
	resp.UserId = rp.UserId
	resp.AccessToken = rp.AccessToken
	resp.RefreshToken = rp.RefreshToken
	return
}
