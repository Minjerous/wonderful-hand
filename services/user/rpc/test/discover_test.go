package test

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"testing"
	"time"
	discover "wonderful-hand-common/discover-services"
	user2 "wonderful-hand-user/rpc/user"
)

func TestDiscover(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379", "127.0.0.1:2479", "127.0.0.1:2579"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}

	resolver.Register(discover.NewBuilder(client))

	userConn, err := grpc.Dial(fmt.Sprintf("%s:///%s", discover.Scheme, "user-rpc"),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // todo 使用证书
	defer userConn.Close()

	userRPC := user2.NewUserServiceClient(userConn)

	resp, err := userRPC.UserLogin(context.Background(), &user2.UserLoginRegisterRequest{
		Username: "iGxnon",
		Password: "114514",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(resp)
}
