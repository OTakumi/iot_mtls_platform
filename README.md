# iot-mtls-platform

mTLS（相互TLS）を用いたIoTデバイスの自動プロビジョニング・認証基盤

## 設計

設計ドキュメント：<https://github.com/OTakumi/iot_mtls_platform/wiki>

## 開発環境

### データベースのセットアップと操作

開発環境でデータベースコンテナを起動し、操作する手順は以下の通りです。

#### 1. .envファイルの準備

プロジェクトのルートディレクトリに`.env`ファイルを配置し、データベースの接続情報などの環境変数を設定します。最低限、以下の変数を設定する必要があります。

```planetext
# Auth Database
AUTH_DB_NAME=iot_auth
AUTH_DB_USER=auth_admin
AUTH_DB_PASS=auth_secure_password_123

# Telemetry Database
TELEM_DB_NAME=iot_telemetry
TELEM_DB_USER=telemetry_admin
TELEM_DB_PASS=telemetry_secure_password_123
```

#### 2. サービスの起動

以下のコマンドで、データベースを含む全てのサービスをバックグラウンドで起動します。

```bash
docker compose up -d
```

#### 3. データベースマイグレーションの実行

サービスの起動後、各データベースのテーブル構造をセットアップするためにマイグレーションを実行します。

- **Authデータベースのマイグレーション:**

  ```bash
  docker compose run --rm auth-migrator up
  ```

- **Telemetryデータベースのマイグレーション:**

  ```bash
  docker compose run --rm telemetry-migrator up
  ```

#### 3.1 データベースマイグレーションのロールバック

最新のマイグレーションを元に戻す（ロールバックする）場合は、以下のコマンドを使用します。

- Authデータベースのロールバック (最新1つ)

  ```bash
  docker compose run --rm auth-migrator down 1
  ```

- Telemetryデータベースのロールバック (最新1つ)

  ```bash
  docker compose run --rm telemetry-migrator down 1
  ```

#### 4. pgAdminでのデータベース確認

[http://localhost:5050](http://localhost:5050) からpgAdminにアクセスできます。

- **ログイン情報:**
  - **Email:** `admin@example.com`
  - **Password:** `admin`

ログイン後、以下の情報で手動でサーバーを登録することで、各データベースの内容を確認できます。

- **Authデータベースの接続情報:**
  - **Host name/address:** `db-auth`
  - **Port:** `5432`
  - **Username:** `.env`で設定した`AUTH_DB_USER`の値
  - **Password:** `.env`で設定した`AUTH_DB_PASS`の値

- **Telemetryデータベースの接続情報:**
  - **Host name/address:** `db-telemetry`
  - **Port:** `5432`
  - **Username:** `.env`で設定した`TELEM_DB_USER`の値
  - **Password:** `.env`で設定した`TELEM_DB_PASS`の値
