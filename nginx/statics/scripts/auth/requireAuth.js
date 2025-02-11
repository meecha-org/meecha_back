async function CheckAuth() {
    try {
        // 認証されているか
        if (!IsAuthed()) {
            // 認証されていない時リダイレクト
            window.location.href = LoginPageURL;
        }

        console.log("認証済み");
    } catch (ex) {
        console.error(ex);
    }
}

CheckAuth();