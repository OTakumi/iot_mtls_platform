# ----- Builder ----- #
FROM golang:1.25 AS builder

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build
# -ldflags="-w -s": デバッグ情報を削除してバイナリサイズを削減
# CGO_ENABLED=0: 静的リンクされたバイナリを作成
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server/main.go

# ----- runner ----- #
FROM alpine:latest

WORKDIR /app

# ca-certificates: 外部HTTPS通信やTLS検証に必要
# tzdata: ログの時刻をJSTなどで正しく表示するために必要
RUN apk --no-cache add ca-certificates tzdata

# 非rootユーザーの作成
# コンテナが乗っ取られた場合のリスクを最小化するため、rootでの実行を避ける
RUN adduser -D -g '' appuser

COPY --from=builder /app/server .

# ボリュームマウントポイントの権限設定
# docker-composeでマウントするcertsディレクトリへのアクセス権を確保
RUN mkdir -p /app/certs && \
    chown -R appuser:appuser /app/certs && \
    chown -R appuser:appuser /app

# ユーザー切り替え
USER appuser

# ポート公開と起動コマンド
EXPOSE 8080

CMD ["./server"]
