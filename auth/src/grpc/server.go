package grpc

import (
	"auth/models"
	"context"
	"errors"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

var (
	// APIキー
	apiKey = ""
)

func RunServer() {
	log.Print("Bind Grpc Server : " + os.Getenv("GRPC_ADDR"))

	// グローバル変数に設定
	apiKey = os.Getenv("GRPC_KEY")

	// 9000番ポートでクライアントからのリクエストを受け付けるようにする
	listen, err := net.Listen("tcp", os.Getenv("GRPC_ADDR"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// GRPC設定
	RegisterAuthServiceServer(grpcServer, &authServer{})

	// 以下でリッスンし続ける
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	log.Print("main end")
}

type authServer struct {
}

// SearchByLabel implements AuthServiceServer.
func (srpc *authServer) SearchByLabel(context.Context, *LabelSearchMessage) (*SearchResult, error) {
	return &SearchResult{
		Users: []*User{},
	}, nil
}

// SearchByName implements AuthServiceServer.
func (srpc *authServer) SearchByName(ctx context.Context, args *NameSearchMessage) (*SearchResult, error) {
	// キーが無効な場合
	if args.APIKEY != apiKey {
		return &SearchResult{}, errors.New("failed to valid key")
	}

	// ユーザーを検索
	users, err := models.SearchUserByName(args.Name)
	if err != nil {
		return &SearchResult{}, err
	}

	// ユーザーリスト
	rUsers := []*User{}

	// ユーザーデータを回す
	for _, val := range users {

		// リストに追加
		rUsers = append(rUsers, &User{
			UserID:   val.UserID,
			UserName: val.Name,
		})
	}

	return &SearchResult{
		Users: rUsers,
	}, nil
}
