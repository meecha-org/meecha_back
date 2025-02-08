package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func CreateCollection(app *pocketbase.PocketBase) {
	app.OnServe().BindFunc(func(evt *core.ServeEvent) error {
		// ラベル作成
		labelCollection, err := InitLabel(app)

		// エラー処理
		if err == nil {
			// ラベル作成
			err = CreateLabel(app, labelCollection, "admin")

			// エラー処理
			if err != nil {
				return err
			}

			// ラベル作成
			err = CreateLabel(app, labelCollection, "owner")

			// エラー処理
			if err != nil {
				return err
			}

			// ラベル作成
			err = CreateLabel(app, labelCollection, "subscriber")

			// エラー処理
			if err != nil {
				return err
			}

			// 設定
			err = SettingProvider(app, labelCollection.Id)

			// エラー処理
			if err != nil {
				return err
			}
		}

		return evt.Next()
	})
}

func SettingProvider(app *pocketbase.PocketBase, LabelID string) error {
	// ユーザーのコレクションを取得
	collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")

	// エラー処理
	if err != nil {
		log.Println(err)
		return err
	}

	// Oauth2 だけ有効化する
	collection.OTP.Enabled = false
	collection.PasswordAuth.Enabled = false
	collection.OAuth2.Enabled = true
	collection.MFA.Enabled = false

	// アラート無効化
	collection.AuthAlert.Enabled = false

	collection.Fields.Add(&core.RelationField{
		Name:          "labels",
		CascadeDelete: false,
		CollectionId:  LabelID,
		MaxSelect:     100,
	})

	// 更新
	err = app.Save(collection)

	// エラー処理
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func InitLabel(app *pocketbase.PocketBase) (*core.Collection, error) {
	collection := core.NewBaseCollection("Label")

	// 閲覧ルール設定
	collection.ViewRule = types.Pointer("@request.auth.id != ''")
	collection.ListRule = nil   //管理者限定
	collection.CreateRule = nil //管理者限定
	collection.UpdateRule = nil //管理者限定
	collection.DeleteRule = nil //管理者限定
	collection.Fields.Add(&core.TextField{
		Name:     "name",
		Required: true,
	})

	// index 追加
	collection.AddIndex("UNIQUE_LABEL", true, "name", "")

	return collection, app.Save(collection)
}

func CreateLabel(app *pocketbase.PocketBase, collection *core.Collection, name string) error {
	// レコード作成
	record := core.NewRecord(collection)

	// 名前設定
	record.Set("name", name)

	return app.Save(record)
}
