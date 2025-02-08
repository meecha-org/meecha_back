// 認証初期化
const pb = new PocketBase("/auth/");

// セッション取得
async function GetSession() {
    // 認証更新
    const authData = await pb.collection('users').authRefresh();

    console.log(authData);

    return authData;
}


function Logout(do_reload) {
    try {
        // "logout" the last authenticated model
        pb.authStore.clear();
    } catch (ex) {
        console.error(ex);
    }

    // リロードが必要なとき
    if (do_reload) {
        // リロードする
        window.location.reload();
    }
}

async function Login(provider) {
    const authData = await pb.collection('users').authWithOAuth2({ provider: provider });

    // after the above you can also access the auth data from the authStore
    console.log(pb.authStore.isValid);
    console.log(pb.authStore.token);
    console.log(pb.authStore.model.id);
}

function GetIcon(userid) {
    return `/auth/icon/${userid}`;
}