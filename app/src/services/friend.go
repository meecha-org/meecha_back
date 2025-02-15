package services

import (
	"errors"
	"new-meecha/grpc"
	"new-meecha/models"
	"new-meecha/utils"
)

type SearchResult struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
}

func SearchByName(name string) ([]SearchResult,error) {
	utils.Println(name)

	// 検索する
	users,err := grpc.SearchByName(name)

	// エラー処理
	if err != nil {
		return []SearchResult{},err
	}

	// 結果を変換
	result := []SearchResult{}
	for _, val := range users {
		// 結果に追加
		result = append(result, SearchResult{
			UserID: val.UserID,
			Name:   val.UserName,
		})
	}

	return result,nil
}

func SendFriendRequest(myId,targetId string) error {
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
	RequestID string `json:"id"`
	SenderID string	`json:"sender"`
	TargetID string `json:"target"`
}

func GetSentRequest(userid string) ([]FriendRequest,error) {
	// 送信済みを取得
	GetRequests,err := models.GetSentRequest(userid)

	// エラー処理
	if err != nil {
		return []FriendRequest{},err
	}

	// 返すリクエスト
	requests := []FriendRequest{}

	for _, request := range GetRequests {
		// リクエストを変換
		requests = append(requests, FriendRequest{
			RequestID: request.RequestID,
			SenderID:  request.SenderID,
			TargetID:  request.TargetID,
		})
	}

	return requests,nil
}

// 受信済みを取得
func GetRecvedRequest(userid string) ([]FriendRequest,error) {
	// 送信済みを取得
	GetRequests,err := models.GetRecvedRequest(userid)

	// エラー処理
	if err != nil {
		return []FriendRequest{},err
	}

	// 返すリクエスト
	requests := []FriendRequest{}

	for _, request := range GetRequests {
		// リクエストを変換
		requests = append(requests, FriendRequest{
			RequestID: request.RequestID,
			SenderID:  request.SenderID,
			TargetID:  request.TargetID,
		})
	}

	return requests,nil
}

// 承認する
func AcceptRequest(requestId string,myId string) error {
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

	// リクエストを削除
	return models.RemoveRequest(requestId)
}

func RejectRequest(myid,requestId string) error {
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

func GetFriendList(userid string) ([]string,error) {
	// フレンドリストを取得
	friends,err := models.GetFriendList(userid)

	// エラー処理
	if err != nil {
		return []string{},err
	}

	return friends,nil
}

// フレンド削除する関数
func RemoveFriend(userid string,targetid string) error {
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
	return models.RemoveFriend(friend.FriendID)
}