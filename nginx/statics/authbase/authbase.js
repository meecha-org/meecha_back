class AuthBase {
    // 各種Auth URL
    static DiscordAuthURL = "/oauth/discord";
    static GoogleAuthURL = "/oauth/google";
    static GithubAuthURL = "/oauth/github";
    static MicrosoftAuthURL = "/oauth/microsoftonline";

    constructor(url) {
        // 最後が/ なら削除
        if (url.endsWith("/")) {
            url = url.slice(0, -1);
        }

        // ベースURL
        this.baseURL = url;
    }

    async login(username, password) {
        const response = await fetch(this.baseURL + '/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });

        if (response.ok) {
            const data = await response.json();
            return data.token;
        } else {
            throw new Error('Login failed');
        }
    }

    // ログアウトする関数
    async Logout() {
        // ログアウトする (リフレッシュトークンでログアウトする)
        const req = await fetch("/auth/logout", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                "Authorization" : localStorage.getItem("token")
            }
        });

        // 失敗した場合
        if (!req.ok) {
            return new Error("Logout failed");
        }

        // ログアウトしたらアクセストークンも削除
        window.sessionStorage.removeItem("actoken");
        window.sessionStorage.removeItem("actime");

        // ローカルストレージからもトークンを削除
        localStorage.removeItem("token");

        return;
    }

    async OauthLogin(provider,LoginCallback) {
        // ポップアップ
        if (provider == "discord") {
            this.openPopup(this.baseURL + AuthBase.DiscordAuthURL);
        }

        if (provider == "google") {
            this.openPopup(this.baseURL + AuthBase.GoogleAuthURL);
        }

        if (provider == "github") {
            this.openPopup(this.baseURL + AuthBase.GithubAuthURL);
        }

        if (provider == "microsoftonline") {
            this.openPopup(this.baseURL + AuthBase.MicrosoftAuthURL);
        }

        // コールバック
        this.LoginCallback = LoginCallback;
    }

    async getToken() {
        // 5分以内なら更新しない
        if (Date.now() - window.sessionStorage.getItem("actime") < 5 * 60 * 1000) {
            console.log("Cache Token");
            return window.sessionStorage.getItem("actoken");
        }

        // トークンを取得
        const req = await fetch(this.baseURL + '/token', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                "Authorization": localStorage.getItem("token")
            }
        });

        // トークンを取得
        if (req.ok) {
            const data = await req.json();

            console.log("Get Token");

            // トークンを保存する
            window.sessionStorage.setItem("actoken", data.token);
            // 現在の時間を保存する
            window.sessionStorage.setItem("actime", Date.now());

            return data.token;
        }

        return null;
    }

    openPopup(url) {
        window.open(url + "?popup=1", "popupWindow", "width=1200,height=800");

        window.addEventListener("message", (event) => {
            if (event.data == "Login-Success") {
                // コールバック
                this.LoginCallback();
            }
        });
    }

    async Post(url,headers,body) {
        // header にトークンを追加
        headers.Authorization = await this.getToken();

        // POST
        const req = await fetch(url, {
            method: 'POST',
            headers: headers,
            body: body
        });

        if (req.ok) {
            return await req.json();
        }

        return null;
    }

    async Get(url,headers) {
        // header にトークンを追加
        headers.Authorization = await this.getToken();

        // GET
        const req = await fetch(url, {
            method: 'GET',
            headers: headers,
        });

        if (req.ok) {
            return await req.json();
        }

        return null;
    }
    
    async GetInfo() {
        try {
            // 情報を取得
            const req = await fetch(this.baseURL + '/me', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    "Authorization": localStorage.getItem("token")
                }
            });

            if (!req.ok) {
                return null;
            }

            // 情報を取得
            return await req.json();

        } catch (error) {
            console.error(error);
            return null;
        }
    }
}

if (typeof module !== 'undefined' && module.exports) {
    module.exports = AuthBase;
}

if (typeof window !== 'undefined') {
    window.AuthBase = AuthBase;
}