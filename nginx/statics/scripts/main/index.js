// ユーザー情報
let userInfo = {};

// 初期化済みか
let IsGeoInit = false;

async function Init() {
    try {
        // 認証情報取得
        const uinfo = await GetSession();

        // ユーザー情報を設定
        userInfo = uinfo["record"];
        console.log(userInfo);

        // 位置情報の監視を開始
        StartWatch();

        // websocket 接続
        Connect();

        // タイマー設定
        setInterval(() => {
            // 位置情報取得
            const data = GetLocation();

            console.log(data.latitude,data.longitude);

            // メッセージ送信
            SendMessage(GeoMessage,JSON.stringify({
                latitude: data.latitude,
                longitude: data.longitude
            }));
        }, 3000);
    } catch (ex) {
        console.error(ex);
    }
}

// 初期化
Init();


// 位置情報のイベントを受け取る
GpsElem.addEventListener(OnLocationChangeName,function(evt) {
    const data = evt.detail;

    console.log(data.latitude,data.longitude);

    if (!IsGeoInit) {
        console.log("ピンを指します");

        // 自分の場所にピンを刺す
        NewPin(userInfo["id"],GetIcon(userInfo["id"]),data.latitude,data.longitude);

        // マップを動かす
        GpsPin.click();

        // 初回の場合
        IsGeoInit = true;
        return;
    }

    console.log("ピンを動かします");
    // 初回以外の場合ピンを移動する
    MovePin(userInfo["id"],data.latitude,data.longitude);
});

// 現在地用のピンを取得
const GpsPin = document.getElementById("pin");

GpsPin.addEventListener("click",function(evt){
    const data = GetLocation();

    console.log(data.latitude,data.longitude);

    // マップ設定
    SetMapView(data.latitude,data.longitude,GetZoom());
})