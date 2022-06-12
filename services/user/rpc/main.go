package main

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
	discover "wonderful-hand-common/discover-services"
	"wonderful-hand-user/rpc/internal/config"
	_ "wonderful-hand-user/rpc/internal/dao"
	"wonderful-hand-user/rpc/internal/service"
	"wonderful-hand-user/rpc/user"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer()
	ln, err := net.Listen("tcp", cfg.Network.Addr)

	user.RegisterUserServiceServer(srv, service.MewUserService(cfg))

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Etcd.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = discover.Register(
		context.Background(),
		client,
		cfg.Server.Name,
		cfg.Network.Host,
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Server %s listened on %s have started\n", cfg.Server.Name, cfg.Network.Addr)
	log.Fatalln(srv.Serve(ln))
}
