package services

import (
	"errors"
	"new-meecha/grpckit"
	"new-meecha/logger"
	"new-meecha/models"
	rediscache "new-meecha/redis-cache"
	"new-meecha/utils"
)

type SearchResult struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
}

func SearchByName(name string) ([]SearchResult, error) {
	utils.Println(name)

	// 検索する
	searchResult, err := grpckit.SearchUser("", name)

	// エラー処理
	if err != nil {
		return []SearchResult{}, err
	}

	// 結果を変換
	result := []SearchResult{}
	for _, val := range searchResult.Users {
		// 結果に追加
		result = append(result, SearchResult{
			UserID: val.UserID,
			Name:   val.Name,
		})
	}

	return result, nil
}

func SendFriendRequest(myId, targetId string) error {
	// 自分が関わるリクエスト取得
	_, err := models.GetFriend(models.FriendArgs{
		UserID:   myId,
		TargetID: targetId,
	})

	// エラー処理
	if err == nil {
		return errors.New("Already friends")
	}

	// リクエストを作成
	return models.CreateRequest(models.FriendRequestArgs{
		UserID:   myId,
		TargetID: targetId,
	})
}

type FriendRequest struct {
	RequestID  string `json:"id"`
	SenderID   string `json:"sender"`
	SenderName string `json:"senderName"`
	TargetID   string `json:"target"`
	TargetName string `json:"targetName"`
}

func GetSentRequest(userid string) ([]FriendRequest, error) {
	// 送信済みを取得
	GetRequests, err := models.GetSentRequest(userid)

	// エラー処理
	if err != nil {
		return []FriendRequest{}, err
	}

	// 自身のユーザー情報を取得
	user, err := grpckit.GetUser(userid)

	// エラー処理
	if err != nil {
		return []FriendRequest{}, err
	}

	// 返すリクエスト
	requests := []FriendRequest{}

	for _, request := range GetRequests {
		// 相手のユーザー情報を取得
		targetUser, err := grpckit.GetUser(request.TargetID)

		// エラー処理
		if err != nil {
			logger.PrintErr(err)
			continue
		}

		// リクエストを変換
		requests = append(requests, FriendRequest{
			RequestID:  request.RequestID,
			SenderID:   request.SenderID,
			SenderName: user.Name,
			TargetID:   request.TargetID,
			TargetName: targetUser.Name,
		})
	}

	return requests, nil
}

// 受信済みを取得
func GetRecvedRequest(userid string) ([]FriendRequest, error) {
	// 送信済みを取得
	GetRequests, err := models.GetRecvedRequest(userid)

	// エラー処理
	if err != nil {
		return []FriendRequest{}, err
	}

	// 自身のユーザー情報を取得
	user, err := grpckit.GetUser(userid)

	// エラー処理
	if err != nil {
		return []FriendRequest{}, err
	}

	// 返すリクエスト
	requests := []FriendRequest{}

	for _, request := range GetRequests {
		// 送信元のユーザー情報を取得
		senderUser, err := grpckit.GetUser(request.SenderID)

		// エラー処理
		if err != nil {
			logger.PrintErr(err)
			continue
		}

		// リクエストを変換
		requests = append(requests, FriendRequest{
			RequestID:  request.RequestID,
			SenderID:   request.SenderID,
			SenderName: senderUser.Name,
			TargetID:   request.TargetID,
			TargetName: user.Name,
		})
	}

	return requests, nil
}

// 承認する
func AcceptRequest(requestId string, myId string) error {
	// リクエストを取得
	request, err := models.GetRequestByID(requestId)

	// エラー処理
	if err != nil {
		return err
	}

	// ターゲットを検証する
	if request.TargetID != myId {
		return errors.New("invalid request")
	}

	// リクエストを承認
	err = models.AcceptRequest(requestId)

	// エラー処理
	if err != nil {
		return err
	}

	// キャッシュ更新 (自分)
	err = CacheFriend(myId)

	// エラー処理
	if err != nil {
		// キャッシュの保存に失敗した場合
		utils.Println(err)
	}

	// キャッシュ更新 (相手)
	err = CacheFriend(request.SenderID)

	// エラー処理
	if err != nil {
		// キャッシュの保存に失敗した場合
		utils.Println(err)
	}

	// リクエストを削除
	return models.RemoveRequest(requestId)
}

func RejectRequest(myid, requestId string) error {
	// リクエストを取得
	request, err := models.GetRequestByID(requestId)

	// エラー処理
	if err != nil {
		return err
	}

	// ターゲットを検証する
	if request.TargetID != myid {
		return errors.New("invalid request")
	}

	// リクエストを削除
	return models.RemoveRequest(requestId)
}

// フレンドのデータ
type Friend struct {
	Name   string `json:"name"` // ユーザー名
	UserID string `json:"id"`   // ユーザーID
}

func GetFriendList(userid string) ([]Friend, error) {
	// キャッシュから取得
	cached, err := rediscache.GetCacheFriend(userid)

	// 成功したとき
	if err == nil {
		returnFriends := []Friend{}

		// フレンドIDを回す
		for _, friend := range cached.FriendIds {
			// ユーザー情報を取得
			user, err := grpckit.GetUser(friend)

			// エラー処理
			if err != nil {
				logger.PrintErr(err)
				continue
			}

			// フレンド情報を追加
			returnFriends = append(returnFriends, Friend{
				Name:   user.Name,
				UserID: user.UserID,
			})
		}
		return returnFriends, nil
	}

	// フレンドリストを取得
	friends, err := models.GetFriendList(userid)

	// エラー処理
	if err != nil {
		return []Friend{}, err
	}

	// キャッシュに保存
	err = CacheFriend(userid)

	// エラー処理
	if err != nil {
		// キャッシュの保存に失敗した場合
		utils.Println(err)
	}

	// フレンド情報を返す
	returnFriends := []Friend{}

	// フレンドIDを回す
	for _, friend := range friends {
		// ユーザー情報を取得
		user, err := grpckit.GetUser(friend)

		// エラー処理
		if err != nil {
			logger.Println(err)
			continue
		}

		// フレンド情報を追加
		returnFriends = append(returnFriends, Friend{
			Name:   user.Name,
			UserID: user.UserID,
		})
	}

	return returnFriends, nil
}

// フレンド削除する関数
func RemoveFriend(userid string, targetid string) error {
	// フレンド取得
	friend, err := models.GetFriend(models.FriendArgs{
		UserID:   userid,
		TargetID: targetid,
	})

	// エラー処理
	if err != nil {
		return err
	}

	// フレンドを削除
	err = models.RemoveFriend(friend.FriendID)

	// エラー処理
	if err != nil {
		return err
	}

	// キャッシュ更新 (自分)
	err = CacheFriend(userid)

	// エラー処理
	if err != nil {
		// キャッシュの保存に失敗した場合
		utils.Println(err)
	}

	// キャッシュ更新 (相手)
	err = CacheFriend(targetid)

	// エラー処理
	if err != nil {
		// キャッシュの保存に失敗した場合
		utils.Println(err)
	}

	// フレンドを削除
	return nil
}

// ユーザーのフレンドをキャッシュする
func CacheFriend(userid string) error {
	// フレンドリストを取得
	friends, err := models.GetFriendList(userid)

	// エラー処理
	if err != nil {
		return err
	}

	// キャッシュ更新
	utils.Println("キャッシュ更新")
	utils.Println(userid)
	utils.Println(friends)

	// キャッシュに保存
	err = rediscache.AddCacheFriend(rediscache.CacheFriendArgs{
		UserID: userid,
		Data: rediscache.FriendCache{
			FriendIds: friends,
		},
	})

	// エラー処理
	if err != nil {
		return err
	}

	return nil
}

// 送信済みリクエストをキャンセルする
func CancelRequest(myid, requestId string) error {
	// リクエストを取得
	request, err := models.GetRequestByID(requestId)

	// エラー処理
	if err != nil {
		return err
	}

	// 自分が送ったリクエストか
	if request.SenderID != myid {
		return errors.New("invalid request")
	}

	// リクエストを削除
	return models.RemoveRequest(requestId)
}
