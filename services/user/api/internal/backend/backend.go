package backend

import (
	"wonderful-hand-user/api/router"
	"wonderful-hand-user/rpc/user"
)

type RpcClis struct {
	UserSrvCli user.UserServiceClient
}

type backend struct {
	router.UnimplementedBackend
	RpcClis RpcClis
}

func New(clis RpcClis) router.Backend {
	return &backend{RpcClis: clis}
}
