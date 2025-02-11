package grpc

import (
	"log"
	"new-meecha/utils"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	gconn  AuthServiceClient = nil
	apiKey                   = ""
)

func Init() {
	utils.Println(os.Getenv("GRPC_SERVER"))

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(os.Getenv("GRPC_SERVER"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	// クライアント
	grpcConn := NewAuthServiceClient(conn)

	// グローバル変数に設定
	gconn = grpcConn
	apiKey = os.Getenv("GRPC_KEY")
}

func SearchByName(name string) ([]*User, error) {
	// 検索する
	result, err := gconn.SearchByName(context.Background(), &NameSearchMessage{
		APIKEY: apiKey,
		Name:   name,
	})

	// エラー処理
	if err != nil {
		return []*User{}, err
	}

	return result.Users, nil
}
