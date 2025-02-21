// 認証初期化
const pb = new PocketBase("/auth/");

// セッション取得
async function GetSession() {
    // 認証更新
    const authData = await pb.collection('users').authRefresh();

    console.log(authData);

    return authData;
}

function IsAuthed() {
    return pb.authStore.isValid;
}

async function GetJwt() {
    // セッションストレージから現在時刻を取得
    const tokenGenTime = window.sessionStorage.getItem("tokenGenTime");

    // 5分以内なら
    if (Date.now() - tokenGenTime < 5 * 60 * 1000) {
        // トークンを取得
        const token = window.sessionStorage.getItem("token");
        return token;
    }

    // セッション更新
    const authData = await GetSession();

    const req = await fetch("/auth/jwt",{
        method: "GET",
        headers : {
            "Authorization" : authData["token"],
        }
    });

    // トークンを取得
    const payload = await req.json();
    const token = payload["result"];

    // トークンをセッションストレージに保存
    window.sessionStorage.setItem("token", token);
    // 現在時刻をセッションストレージに保存
    window.sessionStorage.setItem("tokenGenTime", Date.now());

    return token;
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