package main

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"time"
	discover "wonderful-hand-common/discover-services"
	"wonderful-hand-game/internal/config"
	"wonderful-hand-game/rpc/game"
	"wonderful-hand-game/rpc/service"
	"wonderful-hand-game/server"
	"wonderful-hand-user/rpc/user"
)

func main() {
	c, err := config.ReadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Etcd.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}

	resolver.Register(discover.NewBuilder(client))

	userConn, err := grpc.Dial(fmt.Sprintf("%s:///%s", discover.Scheme, c.GRPC.User),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // todo 使用证书
	defer userConn.Close()

	userRPC := user.NewUserServiceClient(userConn)

	clis := &server.RPCClis{UserRPC: userRPC}

	if err != nil {
		log.Fatalln(err)
	}

	s := server.New(&c, clis)
	go startGRPC(c, s)
	s.Run()
}

func startGRPC(cfg config.Config, s *server.Server) {
	srv := grpc.NewServer()
	ln, err := net.Listen("tcp", cfg.Network.GRPCAddr)
	if err != nil {
		log.Fatalln(err)
	}
	game.RegisterGameServiceServer(srv, service.New(s))

	log.Printf("RPC Server %s listened on %s have started\n", cfg.Server.Name, cfg.Network.GRPCAddr)
	log.Fatalln(srv.Serve(ln))
}
