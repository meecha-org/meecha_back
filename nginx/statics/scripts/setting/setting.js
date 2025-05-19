// ユーザー名
const user_name = document.getElementById("user_name");

// ユーザーアイコン
const user_icon = document.getElementById("user_icon");

// ログアウトボタン
const logout_link = document.getElementById("logout_link");
logout_link.addEventListener("click", async function (evt) {
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
        const userData = await auth.GetInfo();

        if (userData == null) {
            // ログインにリダイレクト
            window.location.href = './login.html';
            return;
        }

        // ユーザー情報を設定
        userInfo = userData;
        console.log(userInfo);

        // ユーザー名設定
        user_name.textContent = userInfo["name"];

        // アイコン設定
        user_icon.src = GetIcon(userInfo["user_id"]);

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
search_user_button.addEventListener("click", async function (evt) {
    // ポップアップ表示
    pupup_search_area.classList.add("is-show");
});

// リクエストを送信
async function SendRequest(uid) {
    // JWT取得
    const jwtToken = await GetJwt();

    // リクエスト送信
    const req = await fetch("/app/friend/request", {
        method: "POST",
        headers: {
            "Authorization": jwtToken,
            "Content-type": "application/json"
        },
        body: JSON.stringify({
            "userid": uid
        })
    });

    if (req.status != 200) {
        return false;
    }

    return true;
}

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

    buttonElement.addEventListener("click", async function (evt) {
        try {
            // リクエストを送信
            if (await SendRequest(uid)) {
                ShowNotify("リクエストを送信しました");
            } else {
                throw "リクエストの送信に失敗しました";
            }

            
        } catch (ex) {
            console.error(ex);
            ShowNotify("リクエストの送信に失敗しました");
        }
    });

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
user_search_button.addEventListener("click", async function (evt) {
    // メッセージを表示
    user_searh_result.textContent = "情報を取得しています";

    try {
        // 検索を実行する
        console.log(search_username_value.value);

        // JWT取得
        const jwtToken = await GetJwt();

        const req = await fetch("/app/friend/search", {
            method: "POST",
            headers: {
                "Authorization": jwtToken,
                "Content-type": "application/json"
            },
            body: JSON.stringify({
                "name": search_username_value.value
            })
        });

        // 結果を取得
        const result = await req.json();

        // 既存の結果を削除
        user_searh_result.innerHTML = "";

        // 結果を表示
        result.forEach(user => {
            const showData = createSearchResultContainer(user["name"], user["userid"]);
            user_searh_result.appendChild(showData);
        });
    } catch (ex) {
        user_searh_result.textContent = "エラーが発生しました";
        console.error(ex);
    }
})

// 送信済みエリア
const sended_log_area = document.getElementById("sended_log_area");
const sended_request_show_area = document.getElementById("sended_request_show_area");
const show_sended_request_area = document.getElementById("show_sended_request_area");
const sended_request_button = document.getElementById("sended_request_button");
sended_request_button.addEventListener("click",async function (evt) {
    sended_log_area.textContent = "データを取得しています";

    try {
        show_sended_request_area.classList.add("is-show");

        // JWT取得
        const jwtToken = await GetJwt();

        // リクエスト送信
        const req = await fetch("/app/friend/sentrequest",{
            method: "GET",
            headers: {
                "Authorization": jwtToken,
                "Content-type": "application/json"
            },
        });

        // エリアを削除
        sended_request_show_area.innerHTML = "";

        // json に変換
        const result = await req.json();
        result.forEach(async function(request) {
            // ユーザー情報を取得
            const req = await fetch("/auth/info/" + request["target"],{
                method: "GET"
            });

            // jsonにする
            const uinfo = await req.json();

            // フレンドリクエストを作成
            sended_request_show_area.appendChild(createSentRequestDataArea(GetIcon(request["target"]),uinfo["UserName"],request["id"]));
        });

        sended_log_area.textContent = "";
    } catch (ex) {
        console.error(ex);
        sended_log_area.textContent = "エラーが発生しました";
    }
});

// リクエストをキャンセルする
async function CancelRequest(requestId) {
    // JWT取得
    const jwtToken = await GetJwt();

    // リクエスト送信
    const req = await fetch("/app/friend/cancel",{
        method: "POST",
        headers: {
            "Authorization": jwtToken,
            "Content-type": "application/json"
        },
        body: JSON.stringify({
            "requestid": requestId
        })
    });

    if (req.status != 200) {
        return false;
    }

    return true;
}

// 送信済みリクエストを作成
function createSentRequestDataArea(userIconSrc, username, buttonId) {
    // リクエストデータエリアのdivを作成
    const requestDataArea = document.createElement('div');
    requestDataArea.className = "request_data_area";

    // ユーザーアイコンのimg要素を作成
    const iconElement = document.createElement('img');
    iconElement.src = userIconSrc; // アイコンのURLを設定
    iconElement.className = "user_icon";

    // ユーザー名のp要素を作成
    const nameElement = document.createElement('p');
    nameElement.className = "username_area";
    nameElement.textContent = username; // ユーザー名を設定

    // ボタン要素を作成
    const buttonElement = document.createElement('button');
    buttonElement.className = "friend_btn";
    buttonElement.id = buttonId; // ボタンのIDを設定
    buttonElement.textContent = "取り消し"; // ボタンのテキストを設定
    buttonElement.addEventListener("click", async function (evt) {
        try {
            // キャセルリクエスト
            if (await CancelRequest(buttonId)) {
                ShowNotify("キャンセルしました");
            } else {
                throw "キャンセルに失敗しました";
            }

            // リロード
            sended_request_button.click();
        } catch (ex) {
            console.error(ex);
            ShowNotify("キャンセルに失敗しました");
        }
    })

    // リクエストデータエリアに要素を追加
    requestDataArea.appendChild(iconElement);
    requestDataArea.appendChild(nameElement);
    requestDataArea.appendChild(buttonElement);

    return requestDataArea; // 作成したdivを返す
}

// 受信済みエリア
const get_request_button = document.getElementById("get_request_button");
const show_recved_request_area = document.getElementById("show_recved_request_area");
const recved_request_show_area = document.getElementById("recved_request_show_area");
const recved_log_area = document.getElementById("recved_log_area");
get_request_button.addEventListener("click",async function (evt) {
    show_recved_request_area.classList.add("is-show");

    recved_log_area.textContent = "データを取得しています";
    
    try {
        // JWT取得
        const jwtToken = await GetJwt();

        const req = await fetch("/app/friend/recvrequest",{
            method: "GET",
            headers: {
                "Authorization": jwtToken,
                "Content-type": "application/json"
            },
        });

        // エリアを削除
        recved_request_show_area.innerHTML = "";

        // json に変換
        const result = await req.json();
        result.forEach(async function(request) {
            // ユーザー情報を取得
            const req = await fetch("/auth/info/" + request["sender"],{
                method: "GET"
            });

            // jsonにする
            const uinfo = await req.json();

            // フレンドリクエストを作成
            recved_request_show_area.appendChild(createRecvedRequestDataArea(GetIcon(request["sender"]),uinfo["UserName"],request["id"]));
        });

        recved_log_area.textContent = "";
    } catch (ex) {
        console.error(ex);
        recved_log_area.textContent = "エラーが発生しました";
    }
});

// 受信済みリクエストを作成
function createRecvedRequestDataArea(userIconSrc, username, buttonId) {
    // リクエストデータエリアのdivを作成
    const requestDataArea = document.createElement('div');
    requestDataArea.className = "request_data_area";

    // ユーザーアイコンのimg要素を作成
    const iconElement = document.createElement('img');
    iconElement.src = userIconSrc; // アイコンのURLを設定
    iconElement.className = "user_icon";

    // ユーザー名のp要素を作成
    const nameElement = document.createElement('p');
    nameElement.className = "username_area";
    nameElement.textContent = username; // ユーザー名を設定

    // 承認ボタンを作成
    const buttonElement = document.createElement('button');
    buttonElement.className = "accept_request_btn";
    buttonElement.id = buttonId; // ボタンのIDを設定
    buttonElement.textContent = "承認"; // ボタンのテキストを設定
    buttonElement.addEventListener("click", async function (evt) {
        try {
            // リクエストを承認
            if (await AcceptRequest(buttonId)) {
                ShowNotify("承認しました");
            } else {
                throw "承認に失敗しました";
            }

            // リロード
            get_request_button.click();
        } catch (ex) {
            console.error(ex);
            ShowNotify("承認に失敗しました");
        }
    });

    // 拒否ボタンを作成
    const rejectButton = document.createElement('button');
    rejectButton.className = "reject_request_btn";
    rejectButton.id = buttonId; // ボタンのIDを設定
    rejectButton.textContent = "拒否"; // ボタンのテキストを設定
    rejectButton.addEventListener("click", async function (evt) {
        try {
            // メッセージを表示
            recved_log_area.textContent = "拒否しています";

            // リクエストを拒否
            if (await RejectRequest(buttonId)) {
                ShowNotify("拒否しました");
            } else {
                throw "拒否に失敗しました";
            }

            // リロード
            get_request_button.click();
        } catch (ex) {
            console.error(ex);
            ShowNotify("エラーが発生しました");
        }
    })

    // リクエストデータエリアに要素を追加
    requestDataArea.appendChild(iconElement);
    requestDataArea.appendChild(nameElement);
    requestDataArea.appendChild(buttonElement);
    requestDataArea.appendChild(rejectButton);

    return requestDataArea; // 作成したdivを返す
}

// フレンドリクエストを承認する
async function AcceptRequest(requestId) {
    // メッセージを表示
    recved_log_area.textContent = "承認しています";

    try {
        // JWT取得
        const jwtToken = await GetJwt();

        // リクエスト送信
        const req = await fetch("/app/friend/accept",{
            method: "POST",
            headers: {
                "Authorization": jwtToken,
                "Content-type": "application/json"
            },
            body: JSON.stringify({
                "requestid": requestId
            })
        });

        // 成功したか
        return req.status == 200;
    } catch (ex) {
        console.error(ex);
        recved_log_area.textContent = "エラーが発生しました";
    }
}

// フレンドリクエストを拒否する
async function RejectRequest(requestId) {
    try {
        // JWT取得
        const jwtToken = await GetJwt();

        // リクエスト送信
        const req = await fetch("/app/friend/reject",{
            method: "POST",
            headers: {
                "Authorization": jwtToken,
                "Content-type": "application/json"
            },
            body: JSON.stringify({
                "requestid": requestId
            })
        });

        // 成功したか
        return req.status == 200;
    } catch (ex) {
        console.error(ex);
        recved_log_area.textContent = "エラーが発生しました";
    }
}

// フレンド一覧
const get_friends_button = document.getElementById("get_friends_button");
const pupup_friends_show_area = document.getElementById("pupup_friends_show_area");
const friend_log_area = document.getElementById("friend_log_area");
const friends_show_area = document.getElementById("friends_show_area");
get_friends_button.addEventListener("click", async function (evt) {
    pupup_friends_show_area.classList.add("is-show");

    friend_log_area.textContent = "データを取得しています";

    // エリアを削除
    friends_show_area.innerHTML = "";

    try {
        // JWT取得
        const jwtToken = await GetJwt();

        // リクエスト送信
        const req = await fetch("/app/friend/list",{
            method: "GET",
            headers: {
                "Authorization": jwtToken,
                "Content-type": "application/json"
            },
        });

        // json に変換
        const result = await req.json();

        result.forEach(async function(userid) {
            // ユーザー情報を取得
            const req = await fetch("/auth/info/" + userid,{
                method: "GET"
            });

            // jsonにする
            const uinfo = await req.json();

            // フレンドリクエストを作成
            const friendDiv = document.createElement('div');
            friendDiv.className = "friend_data_area";
            friendDiv.id = userid;

            // アイコン
            const iconElement = document.createElement('img');
            iconElement.src = GetIcon(userid); // アイコンのURLを設定
            iconElement.className = "user_icon";

            // ユーザー名
            const nameElement = document.createElement('p');
            nameElement.className = "username_area";
            nameElement.textContent = uinfo["UserName"]; // ユーザー名を設定

            // ボタン要素を作成
            const buttonElement = document.createElement('button');
            buttonElement.className = "friend_btn";
            buttonElement.id = userid; // ボタンのIDを設定
            buttonElement.textContent = "取り消し"; // ボタンのテキストを設定
            buttonElement.addEventListener("click", async function (evt) {
                try {
                    // フレンドを削除
                    if (await DeleteFriend(userid)) {
                        ShowNotify("削除しました");
                    } else {
                        throw "削除に失敗しました";
                    }

                    // リロード
                    get_friends_button.click();
                } catch (ex) {
                    console.error(ex);
                    ShowNotify("エラーが発生しました");
                }
            });

            // div に追加
            friendDiv.appendChild(iconElement);
            friendDiv.appendChild(nameElement);
            friendDiv.appendChild(buttonElement);

            friends_show_area.appendChild(friendDiv);

            console.log(uinfo["UserName"]);
        });

        // メッセージを表示
        friend_log_area.textContent = "";
    } catch (ex) {
        console.error(ex);
        friend_log_area.textContent = "エラーが発生しました";
    }
});

// フレンド削除する
async function DeleteFriend(userid) {
    const jwtToken = await GetJwt();

    // リクエスト送信
    const req = await fetch("/app/friend/remove",{
        method: "POST",
        headers: {
            "Authorization": jwtToken,
            "Content-type": "application/json"
        },
        body: JSON.stringify({
            "userid": userid
        })
    });

    return req.status == 200;
}

function ShowNotify(text) {
    toastr.options = {
        "closeButton": false,
        "debug": false,
        "newestOnTop": false,
        "progressBar": true,
        "positionClass": "toast-top-center",
        "preventDuplicates": false,
        "onclick": null,
        "showDuration": "300",
        "hideDuration": "1000",
        "timeOut": "5000",
        "extendedTimeOut": "1000",
        "showEasing": "swing",
        "hideEasing": "linear",
        "showMethod": "fadeIn",
        "hideMethod": "fadeOut"
    }

    toastr["info"](text, "通知")
}