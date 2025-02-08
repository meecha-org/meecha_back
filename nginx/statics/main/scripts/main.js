async function Init() {
    try {
        // 認証情報取得
        const authData = await GetSession();

        console.log(authData["record"]);

        // ユーザーデータ取得
        const UserData = authData["record"];

        // 画像設定
        document.getElementById("usericon").src = GetIcon(UserData["id"]);
    } catch (ex) {
        console.error(ex);

        // 認証していない場合ログインに飛ばす
        window.location.href = LoginURL;
    }
}

// 初期化
Init();

async function DoLogout() {
    try {
        await Logout(true);
    } catch (ex) {
        console.error(ex);

        alert("ログアウトに失敗しました");
    }
}

async function TestJwt() {
    // セッション更新
    const authData = await GetSession();

    console.log(authData["token"]);

    const req = await fetch("/auth/jwt",{
        method: "GET",
        headers : {
            "Authorization" : authData["token"],
        }
    })

    const payload = await req.json();
    console.log(payload["result"]);
}