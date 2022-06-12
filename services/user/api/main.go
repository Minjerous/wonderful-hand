package main

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
	discover "wonderful-hand-common/discover-services"
	"wonderful-hand-common/rest"
	"wonderful-hand-user/api/internal/backend"
	"wonderful-hand-user/api/internal/config"
	"wonderful-hand-user/api/router"
	"wonderful-hand-user/rpc/user"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Printf("parse config failed, using default instead, cause %v\n", err)
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Etcd.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}

	resolver.Register(discover.NewBuilder(client))

	userConn, err := grpc.Dial(fmt.Sprintf("%s:///%s", discover.Scheme, cfg.GRPC.User),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // todo 使用证书
	defer userConn.Close()

	userRPC := user.NewUserServiceClient(userConn)

	clis := backend.RpcClis{UserSrvCli: userRPC}
	rest.SetRouter(router.New(backend.New(clis)))
	rest.SetName(cfg.Server.Name)
	rest.AddMiddleware(rest.MiddlewareCors{})
	rest.SetAddr(cfg.Server.Address)

	if cfg.Server.HTTPCert.Enable {
		log.Fatalln(rest.ServerTLS(cfg.Server.HTTPCert.CertFilePath, cfg.Server.HTTPCert.KeyFilePath))
	}
	log.Fatalln(rest.Serve())
}
