package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	user2 "wonderful-hand-user/rpc/user"
)

func TestLogin(t *testing.T) {
	conn, err := grpc.Dial(":8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	rpc := user2.NewUserServiceClient(conn)
	resp, err := rpc.UserLogin(context.Background(), &user2.UserLoginRegisterRequest{
		Username: "iGxnon",
		Password: "114514",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(resp)
}
