async function Init() {
    try {
        // 認証情報取得
        const authData = await GetSession();

        // 認証済みの場合
        window.location.href = AuthedURL;
    } catch (ex) {
        console.error(ex);
    }
}

// 初期化
Init();

// ログイン関数
async function DoLogin(provider) {
    try {
        // ログイン実行
        const authData = await Login(provider);

        console.log(authData);

        // リロード
        window.location.reload();
    } catch (ex) {
        console.error(ex);

        // 失敗した時アラート出す
        alert("ログインに失敗しました");
    }
}