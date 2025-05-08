# デバッグ用
FROM golang:1.23.6-bullseye as Develop

# apt 更新
RUN apt update
RUN apt install -y git wget curl

# 作業ディレクトリを設定
WORKDIR /root/src

# リリース用
FROM golang:1.23.6-bullseye as Release

# 作業ディレクトリを設定
WORKDIR /root/app

# プロジェクトをコピー
COPY ./src /root/app

# プロジェクトをビルド
RUN go mod download
RUN CGO_ENABLED=0 go build -x -o meecha_server .

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian12
COPY --from=Release /root/app/meecha_server /
CMD ["/meecha_server"]