package services

import (
	"new-meecha/grpc"
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
