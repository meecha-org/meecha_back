async function Init() {
    try {
        const userData = await auth.GetInfo();

        if (userData != null) {
            // 認証済みの場合
            window.location.href = AuthedURL;
            return;
        }

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
        await Login(provider);

    } catch (ex) {
        console.error(ex);

        // 失敗した時アラート出す
        alert("ログインに失敗しました");
    }
}