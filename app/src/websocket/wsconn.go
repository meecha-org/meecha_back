package websocket

import (
	"new-meecha/middlewares"
	"new-meecha/utils"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}

	// websocket の接続
	wsConns = map[string]*WsConn{}
)

func HandleWs(ctx echo.Context) error {
	// 接続をアップグレードする
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}

	// 接続をラップする
	wsconn := InitConn(ws)

	// 読み込む
	msg, err := wsconn.Read()

	// エラー処理
	if err != nil {
		utils.Println(err)

		// 接続を閉じる
		CloseErr := wsconn.Close()

		//エラー処理
		if CloseErr != nil {
			utils.Println(CloseErr)
			return nil
		}

		return nil
	}

	// トークンを検証する
	claim, err := middlewares.ValidToken(msg.Payload)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return nil
	}

	// 認証済みメッセージ
	utils.Println("Websocket Authed")

	// 辞書に保存
	wsConns[claim.UserID] = wsconn
	// 設定
	wsconn.SetUserID(claim.UserID)

	// 認証済みにする
	err = wsconn.Write(WriteMessage{
		Type: AuthComplete,
		Payload: "",
	})

	// エラー処理
	if err != nil {
		utils.Println(err)
		return nil
	}

	// 別スレッドで起動
	go HandleMessage(wsconn)

	return nil
}

func HandleMessage(wsconn *WsConn) {
	for {
		// メッセージを読む
		msg,err := wsconn.Read()

		// エラー処理
		if err != nil {
			utils.Println(err)
			break
		}

		utils.Println(msg)
	}
}
