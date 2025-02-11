// websocket 接続
let wsconn = null;

function Connect() {
    wsconn = new WebSocket("wss://dev-meecha.mattuu.com/app/ws");
    
    // 接続が開いたとき
    wsconn.onopen = onOpen; 

    // メッセージを受け取ったとき
    wsconn.onmessage = (evt) => {
        console.log("onMessage");

        // イベント発火
        onMessage(evt.data);
    };

    // 切断されたとき
    wsconn.onclose = onClose;
}

// 接続されたとき
async function onOpen(evt) {
    // トークン取得
    const token = await GetJwt();

    SendMessage(AuthMessage,token);
}

// メッセージを受け取ったとき
function onMessage(text) {
    console.log(text);
}

// 切断されたとき
function onClose(evt) {
    console.log("切断されました ３秒後に再接続します");

    setTimeout(() => {
        // 接続
        Connect();
    }, 3000);
}

function Disconnect() {

}

// メッセージを送信する関数
function SendMessage(type,payload) {
    // 送信
    wsconn.send(JSON.stringify({
        type,
        payload
    }));
}