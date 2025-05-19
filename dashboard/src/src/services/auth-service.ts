// 認証関連の操作を行うサービス
// 実際の実装ではバックエンドAPIとの通信を行う

export interface AuthUser {
  UserID: string
  Username: string
  IsSystem: number
}

// 管理者が存在するかチェック
export async function checkAdminExists(): Promise<boolean> {
  // 実際の実装ではAPIを呼び出して管理者の存在を確認
  // console.log("Checking if admin exists")

  // API を呼び出す
  const req = await fetch("/auth/admin/status",{
    method: "GET",
  });

  // 結果を検証
  if (!req.ok) {
    throw new Error("Failed to check admin status")
  }

  // 結果を取得
  const res = await req.json();

  console.log(res);

  return res["HasSystemUser"];
}

// ログイン処理
export async function login(email: string, password: string): Promise<void> {
  // API を呼び出す
  const req = await fetch("/auth/admin/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      username: email,
      password,
    }),
  });

  // 結果を検証
  if (!req.ok) {
    throw new Error("Failed to login")
  }

  // レスポンス
  const res = await req.json();

  console.log(res);

  return;
}

// サインアップ処理
export async function signup(email: string, password: string): Promise<void> {
  // 実際の実装ではAPIを呼び出してユーザー登録を行う
  console.log("Signup attempt:", { email, password })

  // API を呼び出してユーザーを作成する
  const req = await fetch("/auth/admin/signup", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      username: email,
      password,
    }),
  });

  // 結果を検証
  if (!req.ok) {
    throw new Error("Failed to create user")
  }

  // レスポンス
  const res = await req.json();

  console.log(res);

  // 成功したらユーザー情報を返す（実際の実装ではAPIからのレスポンスを返す）
  return;
}

// ログアウト処理
export async function logout(): Promise<void> {
  // 実際の実装ではAPIを呼び出してセッションを破棄
  console.log("Logout")

  const req = await fetch("/auth/admin/logout",{
    method: "POST",
  });

  // 結果を検証
  if (!req.ok) {
    throw new Error("Failed to logout")
  }

  return;
}

// 現在のユーザーを取得
export async function getCurrentUser(): Promise<AuthUser | null> {
  // ユーザー情報を取得
  const req = await fetch("/auth/admin/info",{
    method: "GET",
  });

  // 結果を検証
  if (!req.ok) {
    throw new Error("Failed to get current user")
  }

  // 結果を取得
  const res = await req.json();
  console.log(res);

  const reutnData: AuthUser = {
    UserID: res["UserID"],
    Username: res["Username"],
    IsSystem: res["IsSystem"],
  }

  return reutnData;
}

// ユーザーが認証済みかチェック
export async function isAuthenticated(): Promise<boolean> {
  const user = await getCurrentUser()
  return user !== null
}
