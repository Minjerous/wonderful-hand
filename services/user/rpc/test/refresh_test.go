package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	user2 "wonderful-hand-user/rpc/user"
)

func TestRefresh(t *testing.T) {
	conn, err := grpc.Dial(":8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	rpc := user2.NewUserServiceClient(conn)
	resp, err := rpc.UserTokenRefresh(context.Background(), &user2.UserTokenRefreshRequest{
		// RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTUwNTQzMjQsImlhdCI6MTY1NDk2NzkyNCwiaXNzIjoidXNlci1ycGMiLCJpc19hY2Nlc3MiOnRydWV9.Bvqw9QzBayRTFoPHTUGxkDf2ZGLsUBMNVXp4LSMy5dc",
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTU0MDAyMTUsImlhdCI6MTY1NDk2ODIxNSwiaXNzIjoidXNlci1ycGMiLCJpc19hY2Nlc3MiOmZhbHNlfQ.wme9i5wBxnqHwMFMhVRjQbEgVnkSklt_9Ui19ZW9Y3w",
	})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(resp)
}
