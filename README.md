# Pocketbase 認証キットリポジトリ

## 環境構築
1. 自己署名SSL の鍵を生成するために
   - Mac Linux の人
     - ```./Genkey.sh ```
     を実行して全てエンター
    - Windows の人
      - ``` ./Genkey.bat ```
      を実行して全てエンター
    - しっぱしいた場合
        ファイル内容に書いてあるものをコピーして docker-compose.yaml があるディレクトリで実行してください。
2. コンテナを起動
    - 起動コマンド 
    ```
    docker compose up -d --build
    ```
3. 各種データ設定
   1. (必須) [Google Cloud Console](https://console.cloud.google.com/welcome) にアクセスして　Oauth 用の Client ID と シークレットを作成
   2. (オプション) [Github開発者設定](https://github.com/settings/developers) にアクセスして　Oauth 用の Client ID と シークレットを作成
   3. (オプション) [Discord開発者コンソール](https://discord.com/developers/applications) にアクセスして　Oauth 用の Client ID と シークレットを作成
   - リダイレクト URL には全て 
    ```
    https://dev-meecha.mattuu.com/auth/api/oauth2-redirect
    ```
    を設定してください。
4. Pocketbase の初期設定
    1. [管理画面](https://dev-meecha.mattuu.com/auth/_)にアクセスしてユーザーを作成
    2. [設定画面](https://dev-meecha.mattuu.com/auth/_/#/settings) にアクセス
    3. Application URL を 
    ```
    https://dev-meecha.mattuu.com/auth/
    ```
    にする
    4. [認証プロバイダ管理画面](https://dev-meecha.mattuu.com/auth/_/#/settings/auth-providers)で (使うもののみ) Google, Github,Discord を設定する
5. 認証テスト
   - [ホーム](https://dev-meecha.mattuu.com/statics/)にアクセス
   - ログインしてみる

# 終わり！

# redis の使用状況
- 0
  - フレンドデータのキャッシュに使用
- 1 
  - 位置情報の管理に使用
- 2
  - 過去の応答のキャッシュに使用

# 位置情報の取り扱い
redis の 1番にキー: meecha_geo で位置情報を保存する
キー: meecha_geo_expiry　に 有効期限を保存する
5秒ごとにmeecha_geo_expiry から有効期限切れのキーを取得し存在したら 両方から消す