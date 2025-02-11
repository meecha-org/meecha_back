package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

// websocket の接続
type WsConn struct {
	userid string
	wsconn *websocket.Conn
}

func (conn *WsConn) SetUserID(userid string) (error) {
	conn.userid = userid
	return nil
}

func (conn *WsConn) Read() (ReadMessage,error) {
	// メッセージを読みこむ
	_, msg, err := conn.wsconn.ReadMessage()
	if err != nil {
		return ReadMessage{},err
	}

	decodeMsg := ReadMessage{}
	// json を戻す
	err = json.Unmarshal(msg,&decodeMsg)
	
	// エラー処理
	if err != nil {
		return ReadMessage{},err
	}

	return decodeMsg,nil
}

func (conn *WsConn) Write(msg WriteMessage) (error) {
	// メッセージを読みこむ
	err := conn.wsconn.WriteJSON(msg)
	if err != nil {
		return err
	}
	
	return nil
}

func (conn *WsConn) Close() (error) {
	// 切断
	return conn.wsconn.Close()
}

func InitConn(wsconn *websocket.Conn) *WsConn {
	return &WsConn{
		wsconn: wsconn,
	}
}
