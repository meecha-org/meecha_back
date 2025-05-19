// 認証初期化
const auth = new AuthBase('/auth/'); 

async function GetToken() {
    return auth.getToken();
}

function IsAuthed() {
    return auth.IsAuthed();
}

async function GetJwt() {
    // セッション更新
    const acToken = await GetToken();

    return acToken;
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
    return `/auth/assets/${userid}.png`;
}