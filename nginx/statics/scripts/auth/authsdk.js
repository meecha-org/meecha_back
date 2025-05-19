// 認証初期化
const auth = new AuthBase('/auth/'); 

async function GetToken() {
    return auth.getToken();
}

function IsAuthed() {
    return auth.IsAuthed();
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
    const acToken = await GetToken();

    const req = await fetch("/auth/jwt",{
        method: "GET",
        headers : {
            "Authorization" : acToken,
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
    try {
        return await auth.OauthLogin(provider,function() {
            // リロードする
            window.location.reload();
        });
    } catch (ex) {
        console.error(ex);
    }
}

function GetIcon(userid) {
    return `/auth/icon/${userid}`;
}