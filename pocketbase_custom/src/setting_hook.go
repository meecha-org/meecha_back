package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/osutils"
)


func InstallHook(app *pocketbase.PocketBase) {
	// 初期化イベント
	app.OnServe().BindFunc(func(evt *core.ServeEvent) error {
		// 初期化ファンクション差し替え
		evt.InstallerFunc = func(app core.App, systemSuperuser *core.Record, baseURL string) error {
			token, err := systemSuperuser.NewStaticAuthToken(30 * time.Minute)
			if err != nil {
				return err
			}

			// launch url (ignore errors and always print a help text as fallback)
			url := fmt.Sprintf("%s/_/#/pbinstal/%s", strings.TrimRight(os.Getenv("AUTHROOT"), "/"), token)
			_ = osutils.LaunchURL(url)
			color.Magenta("\n(!) Launch the URL below in the browser if it hasn't been open already to create your first superuser account:")
			color.New(color.Bold).Add(color.FgCyan).Println(url)
			color.New(color.FgHiBlack, color.Italic).Printf("(you can also create your first superuser by running: %s superuser upsert EMAIL PASS)\n\n", os.Args[0])

			return nil
		}

		// 初期化実行
		return evt.Next()
	})

	// 起動まえにURLを差し替える
	app.OnBootstrap().BindFunc(func(evt *core.BootstrapEvent) error {
		// 初期化する
		err := evt.Next()

		// エラー処理
		if err != nil {
			return err
		}

		// 設定取得
		settings := evt.App.Settings()

		// アプリのURL変更
		settings.Meta.AppURL = os.Getenv("AUTHROOT")

		// 設定保存
		err = evt.App.Save(settings)

		// エラー処理
		if err != nil {
			return err
		}

		return evt.Next()
	})
}
