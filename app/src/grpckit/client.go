package grpckit

import (
	"context"
	"os"

	"google.golang.org/grpc"
)

var (
	client AuthBaseServiceClient
)

func Init() error {
	// GRPC クライアント初期化
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(os.Getenv("GRPC_ADDR"), grpc.WithInsecure())
	if err != nil {
		return err
	}

	// クライアントを作成する
	client = NewAuthBaseServiceClient(conn)

	return nil
}

func GetUser(userID string) (*User, error) {
	return client.GetUser(context.Background(), &GetUserRequest{UserID: userID})
}

func SearchUser(email string, name string) (*SearchResult, error) {
	return client.SearchUser(context.Background(), &SearchRequest{
		Email: email,
		Name:  name,
	})
}

func GetLabel(labelName string) (*Label, error) {
	return client.GetLabel(context.Background(), &GetLabelRequest{LabelName: labelName})
}