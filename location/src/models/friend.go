package models

type Friend struct {
	FriendID string `gorm:"unique"`         //フレンドID
	SenderID string `gorm:"primaryKey"`     //送信者ID
	TargetID string `gorm:"primaryKey"`     //相手のID
	Created  int64  `gorm:"autoCreateTime"` // 作成時間としてUNIX秒を使用する
}

// フレンド一覧
func GetFriendList(userid string) ([]string, error) {
	// フレンドリストを取得
	friends := []Friend{}

	// フレンドリストを取得
	result := dbconn.Where(&Friend{
		SenderID: userid,
	}).Or(&Friend{
		TargetID: userid,
	}).Find(&friends)

	// エラー処理
	if result.Error != nil {
		return []string{}, result.Error
	}

	// 返すリスト
	return_friends := []string{}
	for _, friend := range friends {
		if friend.SenderID == userid {
			return_friends = append(return_friends, friend.TargetID)
		} else {
			return_friends = append(return_friends, friend.SenderID)
		}
	}

	return return_friends, nil
}