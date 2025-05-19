package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"google.golang.org/grpc"
)

var (
	// アプリ
	gApp = &pocketbase.PocketBase{}

	// APIキー
	apiKey = ""
)

func RunServer(app *pocketbase.PocketBase) {
	log.Print("Bind Grpc Server : " + os.Getenv("GRPC_ADDR"))

	// グローバル変数に設定
	gApp = app
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
		return &SearchResult{},errors.New("failed to valid key")
	}

	users := []struct {
		Id    string `db:"id" json:"id"`
		Name  string `db:"name" json:"name"`
	}{}

	// 検索
	err := gApp.DB().
		Select("id", "name","labels").
		From("users").
		AndWhere(dbx.Like("name", args.Name)).
		All(&users)
	
	// エラー処理
	if err != nil {
		return &SearchResult{},errors.New("failed to get user")
	}

	// ユーザーリスト
	rUsers := []*User{}	

	// ユーザーデータを回す
	for _, val := range users {

		// リストに追加
	    rUsers = append(rUsers, &User{
			UserID:        val.Id,
			UserName:      val.Name,
		})
	}
	
	return &SearchResult{
		Users: rUsers,
	},nil
}
