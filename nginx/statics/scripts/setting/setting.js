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
        // StartWatch();
    } catch (ex) {
        console.error(ex);
    }
}

// 初期化
Init();

// ユーザー検索
const pupup_search_area = document.getElementById("pupup_search_area");
const search_user_button = document.getElementById("search_user_button");
search_user_button.addEventListener("click",async function (evt) {
    // ポップアップ表示
    pupup_search_area.classList.add("is-show");
});

function createSearchResultContainer(name, uid) {
    // ベースのdivを作成
    const resultDiv = document.createElement('div');
    resultDiv.className = "search_result";

    // アイコン要素を作成
    const iconElement = document.createElement('img');
    iconElement.className = "search_result_icon result";
    iconElement.src = GetIcon(uid); // アイコンのURLを設定

    // 名前要素を作成
    const nameElement = document.createElement('p');
    nameElement.className = "search_result_name result";
    nameElement.textContent = name; // 名前を設定

    // ボタン要素を作成
    const buttonElement = document.createElement('button');
    buttonElement.className = "send_request_button result";
    buttonElement.id = uid; // IDを設定
    buttonElement.textContent = "送信"; // ボタンのテキストを設定

    // ベースのdivに要素を追加
    resultDiv.appendChild(iconElement);
    resultDiv.appendChild(nameElement);
    resultDiv.appendChild(buttonElement);

    return resultDiv; // 作成したdivを返す
}

// 検索エリア
const user_searh_result = document.getElementById("user_searh_result");
const search_username_value = document.getElementById("search_username_value");
const user_search_button = document.getElementById("user_search_button");
user_search_button.addEventListener("click",async function (evt) {
    // メッセージを表示
    user_searh_result.textContent = "情報を取得しています";

    try {
        // 検索を実行する
        console.log(search_username_value.value);

        // JWT取得
        const jwtToken = await GetJwt();

        const req = await fetch("/app/friend/search",{
            method: "POST",
            headers: {
                "Authorization" : "Bearer " + jwtToken,
                "Content-type" : "application/json"
            },
            body: JSON.stringify({
                "name" : search_username_value.value
            })
        });

        // 結果を取得
        const result = await req.json();

        // 既存の結果を削除
        user_searh_result.innerHTML = "";

        // 結果を表示
        result.forEach(user => {
            const showData = createSearchResultContainer(user["name"],user["userid"]);
            user_searh_result.appendChild(showData);
        });
    } catch (ex) {
        user_searh_result.textContent = "エラーが発生しました";
        console.error(ex);
    }
})