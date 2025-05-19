// ユーザー関連の操作を行うサービス
// 実際の実装ではバックエンドAPIとの通信を行う

export interface User {
  id: string
  name: string
  email: string
  provider: string
  providerId: string
  avatar: string
  labels: string[]
  createdAt: string
  banned: boolean
}

// ユーザー作成用インターフェース（パスワードを含む）
export interface CreateUserRequest {
  id?: string // オプショナル、指定がなければサーバー側で生成
  name: string
  email: string
  password: string
  provider: string
  providerId: string
  avatar: string
  labels: string[]
}

// ユーザー一覧
let users: User[] = [];

// ユーザー一覧を取得
export async function getUsers(): Promise<User[]> {
  // APIからデータ取得
  const req = await fetch("/auth/api/user/all", {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    }
  });

  // ユーザー一覧を取得
  users = await req.json();

  return users;
}

// ユーザーを検索
export async function searchUsers(query: string): Promise<User[]> {
  // 実際の実装ではAPIからデータを取得
  return users.filter(
    (user) =>
      user.name.toLowerCase().includes(query.toLowerCase()) ||
      user.email.toLowerCase().includes(query.toLowerCase()) ||
      user.id.toLowerCase().includes(query.toLowerCase()),
  )
}

// ユーザーを作成
export async function createUser(user: CreateUserRequest): Promise<User> {
  // 実際の実装ではAPIを呼び出してユーザーを作成
  console.log("Create user:", user)

  // パスワードのバリデーション（実際の実装ではサーバーサイドで行う）
  if (!user.password || user.password.length < 8) {
    throw new Error("パスワードは8文字以上で入力してください")
  }

  // 新しいユーザーを作成して返す
  const newUser: User = {
    id: user.id || `usr_${Math.random().toString(36).substring(2, 8)}`,
    name: user.name,
    email: user.email,
    provider: user.provider,
    providerId: user.providerId,
    avatar: user.avatar,
    labels: user.labels,
    createdAt: new Date().toISOString().split("T")[0],
    banned: false,
  }

  users.push(newUser)
  return newUser
}

// ユーザーを更新
export async function updateUser(user: Partial<User> & { id: string }): Promise<void> {
  // 実際の実装ではAPIを呼び出してユーザーを更新
  console.log("Update user:", user)

  // ユーザーを更新する
  const req = await fetch("/auth/api/user", {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(user),
  });

  // 結果を検証
  if (!req.ok) {
    throw new Error("Failed to update user")
  }

  return;
}

// ユーザーを削除
export async function deleteUser(userId: string): Promise<void> {
  // 実際の実装ではAPIを呼び出してユーザーを削除
  console.log("Delete user:", userId)

  // ユーザーを削除する
  const req = await fetch("/auth/api/user", {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      "userId": userId
    },
  });

  // 結果を検証
  if (!req.ok) {
    throw new Error("Failed to delete user")
  }

  return;
}

// ユーザーのBANステータスを切り替え
export async function toggleUserBan(userId: string): Promise<User> {
  // 実際の実装ではAPIを呼び出してユーザーのBANステータスを切り替え
  console.log("Toggle ban status for user:", userId)

  // モックデータを更新して返す
  const index = users.findIndex((u) => u.id === userId)
  if (index !== -1) {
    users[index].banned = !users[index].banned

    // リクエストを送る
    const req = await fetch("/auth/api/user/ban",{
      method: "PUT",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        "IsBanned" : users[index].banned,
        "UserID" : userId,
      })
    });

    // エラー処理
    if (!req.ok) {
      throw new Error("failed to toggle ban")
    }

    return users[index]
  }

  throw new Error("User not found")
}

// ユーザーとしてログイン
export async function loginAsUser(userId: string): Promise<User> {
  // 実際の実装ではAPIを呼び出してユーザーとしてログイン
  console.log("Login as user:", userId)

  // モックデータからユーザーを取得
  const user = users.find((u) => u.id === userId)
  if (!user) {
    throw new Error("User not found")
  }

  // セッションストレージにログイン中のユーザー情報を保存
  sessionStorage.setItem("login_as_user", JSON.stringify(user))

  return user
}

// 現在ログイン中のユーザーを取得
export async function getCurrentLoginAs(): Promise<User | null> {
  // セッションストレージからログイン中のユーザー情報を取得
  const userJson = sessionStorage.getItem("login_as_user")

  if (userJson) {
    return JSON.parse(userJson) as User
  }

  return null
}

// ログイン中のユーザーをクリア
export async function clearLoginAs(): Promise<void> {
  // セッションストレージからログイン中のユーザー情報を削除
  sessionStorage.removeItem("login_as_user")
}
