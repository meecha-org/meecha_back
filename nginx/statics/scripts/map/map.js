// マップを生成
var map = L.map('show_map').setView([51.505, -0.09], 13);

// レイヤを設定
// L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
//     attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
// }).addTo(map);
var tileLayer = L.tileLayer('https://mt1.google.com/vt/lyrs=r&x={x}&y={y}&z={z}', {
    attribution: "<a href='https://developers.google.com/maps/documentation' target='_blank'>Google Map</a>"
});
tileLayer.addTo(map);

// マップの位置を設定
function SetMapView(latitude, longitude, zoom = 15) {
    // 設定
    map.setView([latitude, longitude], zoom);
}

// ズームを取得
function GetZoom() {
    return map._zoom;
}

// ピンのjson
let mapPins = {};

// 新くピンを追加する関数
function NewPin(PinID, IconURL, latitude, longitude) {
    //ピンを画像に変える    
    var pinIcon = L.icon({
        iconUrl: IconURL,
        iconSize: [40, 40],
        iconAnchor: [37, 75],
        popupAnchor: [0, -70],
        className: "UserIcon",
    });

    const pin = L.marker([latitude, longitude], { icon: pinIcon }).addTo(map);

    // jsonに保存
    mapPins[PinID] = {
        Pin: pin,
        latitude: latitude,
        longitude: longitude,
        IconURL: IconURL
    };
}

// ピンを移動する関数
function MovePin(PinID, latitude, longitude) {
    const targetPin = mapPins[PinID];

    // ピンを削除する
    map.removeLayer(targetPin.Pin);

    // 新くピンを撃つ
    NewPin(PinID,targetPin.IconURL,latitude,longitude);
}

// ピンが存在するか
function ExistPin(PinID) {
    return PinID in mapPins;
}

// ピンを削除する
function RemovePin(PinID) {
    const targetPin = mapPins[PinID];
    map.removeLayer(targetPin.Pin);
    delete mapPins[PinID];
}