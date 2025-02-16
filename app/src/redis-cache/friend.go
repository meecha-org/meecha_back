package rediscache

import (
	"context"
	"new-meecha/utils"
	"time"

	// "new-meecha/utils"

	"github.com/vmihailenco/msgpack/v5"
)

type CacheFriendArgs struct {
	UserID string //相手のユーザーID
	Data FriendCache //フレンドのID一覧
}
type FriendCache struct {
	FriendIds []string	//フレンドのID一覧
}

// キャッシュを追加
func AddCacheFriend(args CacheFriendArgs) error {
	// バックグラウンドコンテキスト
	ctx := context.Background()

	// msgpack に変換
	idsbin, err := msgpack.Marshal(args.Data)
    if err != nil {
        return err
    }

	utils.Println(args.UserID)
	utils.Println(args.Data)

	// redis に保存
	err = conn.Set(ctx,args.UserID,idsbin,time.Second * 300).Err()

	// エラー処理
	if err != nil {
		return err
	}

	return nil
}

// キャッシュから取得
func GetCacheFriend(userid string) (FriendCache,error) {
	// バックグラウンドコンテキスト
	ctx := context.Background()

	// redis に保存	
	result,err := conn.Get(ctx,userid).Result()	

	// エラー処理
	if err != nil {
		return FriendCache{},err
	}

	// 返すIDS
	returnIds := FriendCache{}

	// msgpack から変換
	if err := msgpack.Unmarshal([]byte(result), &returnIds); err != nil {
		return FriendCache{},err
	}

	return returnIds, nil
}


// キャッシュを削除
func DelCacheFriend(args CacheFriendArgs) error {
	// バックグラウンドコンテキスト
	ctx := context.Background()

	// redis に保存
	_,err := conn.Del(ctx,args.UserID).Result()

	// エラー処理
	if err != nil {
		return err
	}

	return nil
}

// キャッシュに存在するか
func ExistCacheFriend(args CacheFriendArgs) bool {
	// バックグラウンドコンテキスト
	ctx := context.Background()

	// redis に保存
	_,err := conn.Exists(ctx,args.UserID).Result()

	// エラー処理
	if err != nil {
		return false
	}

	return true
}