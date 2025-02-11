// イベント伝える用のオブジェクト作成
const GpsElem = document.createElement(null);

// 位置情報が変わった時のイベント名
const OnLocationChangeName = "LocationChange";

// 位置情報が更新された時
function success(pos) {
    const crd = pos.coords;

    // イベント発火
    PostLocation(crd.latitude,crd.longitude);
}

let GPSData = {
    "longitude" : 0,
    "latitude" : 0,
};

// 位置情報を更新する関数
function PostLocation(latitude,longitude) {
    // GPSのデータ
    GPSData = {
        "longitude" : longitude,
        "latitude" : latitude,
    }

    const LocationEvent = new CustomEvent(OnLocationChangeName,{
        detail: GPSData
    });

    // イベント発火
    GpsElem.dispatchEvent(LocationEvent);
}

// 位置情報を取得する関数
function GetLocation() {
    return GPSData;
}

// エラーが起きた時
function error(err) {
    console.error(`ERROR(${err.code}): ${err.message}`);
}

// 位置情報の設定
options = {
    enableHighAccuracy: true,
    timeout: 5000,
    maximumAge: 0,
};

// GPSの監視ID
let watchid = "";

// 監視を開始する関数
function StartWatch() {
    // 監視開始
    console.log("GPS監視を開始します");

    // 監視を開始する
    watchid = navigator.geolocation.watchPosition(success, error, options);
};

// 監視を停止する関数
function StopWatch() {
    console.log("GPS監視を停止します");
    navigator.geolocation.clearWatch(watchid);
}