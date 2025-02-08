package main

import (
	"log"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	// "github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	// "github.com/pocketbase/pocketbase/models"
)

// ユーザーデータ (トークンで返すやつ)
type UserData struct {
	UserID      string    // ユーザーID
	UserName    string    // ユーザー名
	Email       string    // メールアドレス
	Labels      []string  // ユーザーについたラベル
	ProviderUID string    // 認証プロバイダのユーザーID
	Provider    string    // 認証プロバイダ
	Created     time.Time // 作成日
	Updated     time.Time // 更新日
}

// IDを指定したら帰ってくるユーザーデータ
type UserInfo struct {
	UserName string   // ユーザー名
	Labels   []string // ユーザーについたラベル
}

func main() {
	app := pocketbase.New()

	// ユーザー認証
	app.OnServe().BindFunc(func(evt *core.ServeEvent) error {
		// api group
		// ユーザーを認証する関数
		evt.Router.POST("/jwt", func(evt *core.RequestEvent) error {
			// ユーザー取得
			info, err := evt.RequestInfo()

			// エラー処理
			if err != nil {
				app.Logger().Error("error",err)
				return evt.Error(http.StatusInternalServerError, "something went wrong", map[string]validation.Error{
					"title": validation.NewError("message", "failed to get auth"),
				})
			}

			authRecord := info.Auth

			// 認証プロバイダの情報取得
			records,err := app.FindAllExternalAuthsByRecord(authRecord)

			// エラー処理
			if err != nil {
				app.Logger().Error("error",err)
				return evt.Error(http.StatusInternalServerError, "something went wrong", map[string]validation.Error{
					"title": validation.NewError("message", "failed to get auth"),
				})
			}

			// ラベルのIDを取得する
			labels := authRecord.Get("labels").([]string)

			// ラベルの文字列取得
			return_labels := GetLabels(app,labels)
			
			// 認証プロバイダの情報を取得
			provider := records[0]

			// 返すデータ
			return_data := UserData{
				UserID:      authRecord.Id,
				UserName:    authRecord.GetString("name"),
				Email:       authRecord.GetString("email"),
				Labels:      return_labels,
				ProviderUID: provider.ProviderId(),
				Provider:    provider.Provider(),
				Created:     authRecord.GetDateTime("created").Time(),
				Updated:     authRecord.GetDateTime("updated").Time(),
			}

			return evt.JSON(http.StatusOK, map[string]UserData{"result": return_data})
			
		}).Bind(apis.RequireAuth("users"))

		evt.Router.GET("/icon/{userid}",func(evt *core.RequestEvent) error {
			// 画像を取得するコレクション
			targetCollection := "_pb_users_auth_"

			// ユーザーID取得
			userid := evt.Request.PathValue("userid")
			
			// ユーザーを取得する
			user, err := app.FindRecordById(targetCollection, userid)

			// エラー処理
			if err != nil {
				app.Logger().Error("error",err)
				return evt.JSON(http.StatusInternalServerError, map[string]string{
					"result": "Failed to get user",
				})
			}

			// アバターのURL
			avatar := app.Settings().Meta.AppURL + "/api/files/" + user.Collection().Id + "/" + user.Id + "/" + user.GetString("avatar")

			return evt.Redirect(http.StatusTemporaryRedirect, avatar)
		})

		evt.Router.GET("/info/{userid}",func(evt *core.RequestEvent) error {
			// 画像を取得するコレクション
			targetCollection := "_pb_users_auth_"

			// ユーザーID取得
			userid := evt.Request.PathValue("userid")
			
			// ユーザーを取得する
			user, err := app.FindRecordById(targetCollection, userid)

			// エラー処理
			if err != nil {
				app.Logger().Error("error",err)
				return evt.JSON(http.StatusInternalServerError, map[string]string{
					"result": "Failed to get user",
				})
			}

			// ラベルのIDを取得する
			labels := user.Get("labels").([]string)

			// ラベルの文字列取得
			return_labels := GetLabels(app,labels)

			return evt.JSON(http.StatusOK,map[string]UserInfo{
				"result" : {
					UserName: user.GetString("name"),
					Labels:   return_labels,
				},
			})
		})

		return evt.Next()
	})

	// 設定フックする
	InstallHook(app)

	// コレクション作成
	CreateCollection(app)

	// 起動
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func GetLabels(app *pocketbase.PocketBase,labels []string) ([]string) {
	// 返すラベル
	return_labels := []string{}

	// ラベルを回す
	for _, val := range labels {
		// ラベルを取得する
		label, err := app.FindRecordById("Label", val)

		// エラー処理
		if err != nil {
			log.Println(err)
			continue
		}

		// ラベルを追加
		return_labels = append(return_labels, label.GetString("name"))
	}

	return return_labels
}