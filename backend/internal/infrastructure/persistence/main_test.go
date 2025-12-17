package persistence_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"backend/internal/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	// テスト全体のセットアップ
	err := setupTestDatabase()
	if err != nil {
		log.Fatalf("テストデータベースのセットアップに失敗しました: %v", err)
	}

	// 全てのテストを実行
	code := m.Run()

	// 後処理
	// DB接続を閉じる
	sqlDB, err := testDB.DB()
	if err == nil {
		sqlDB.Close()
	}

	os.Exit(code)
}

func setupTestDatabase() error {
	// 設定の読み込み
	// ../../config を指定して config.test.yaml を見つける
	// 修正: ../../config -> ../../../config
	cfg, err := config.LoadDBConfig("../../../config", "config.test")
	if err != nil {
		return fmt.Errorf("設定ファイルの読み込みに失敗: %w", err)
	}

	// DSNを構築
	dsn := cfg.DSN()
	migrateDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	// マイグレーションの実行
	// ../../../infra/db-auth/migrations を指定
	// 修正: ../../../infra/db-auth/migrations -> ../../../../infra/db-auth/migrations
	migrationURL := "file://../../../../infra/db-auth/migrations"

	mi, err := migrate.New(migrationURL, migrateDSN)
	if err != nil {
		return fmt.Errorf("migrateインスタンスの作成に失敗: %w", err)
	}

	// DBの状態を一度クリアしてからマイグレーションを実行
	if err := mi.Down(); err != nil && err != migrate.ErrNoChange {
		log.Printf("migrate downに失敗しましたが、テストを続行します: %v", err)
	}
	if err := mi.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate upに失敗: %w", err)
	}

	// GORMでのDB接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("DBへの接続に失敗: %w", err)
	}

	testDB = db
	return nil
}

func cleanupTable(t *testing.T, tableName string) {
	t.Helper()
	err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", tableName)).Error
	if err != nil {
		t.Fatalf("テーブルのクリーンアップに失敗 (%s): %v", tableName, err)
	}
}
