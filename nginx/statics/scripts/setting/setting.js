// ユーザー名
const user_name = document.getElementById("user_name");

// ユーザーアイコン
const user_icon = document.getElementById("user_icon");

// ログアウトボタン
const logout_link = document.getElementById("logout_link");
logout_link.addEventListener("click",async function (evt) {
    // イベントキャンセル
    evt.preventDefault();

    // ログアウト
    await Logout(true);
});

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

        // ユーザー名設定
        user_name.textContent = userInfo["name"];

        // アイコン設定
        user_icon.src = GetIcon(userInfo["id"]);

        // // 位置情報の監視を開始
        StartWatch();
    } catch (ex) {
        console.error(ex);
    }
}

// 初期化
Init();