package persistence_test

import (
	"fmt"
	"log"
	"os"
	"testing"

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
	// 環境変数からテストDBのDSNを取得
	dsnTest := os.Getenv("DSN_TEST")
	if dsnTest == "" {
		return fmt.Errorf("環境変数 DSN_TEST が設定されていません")
	}

	// マイグレーションの実行
	migrationURL := "file://../../../../infra/db-auth/migrations"

	// migrate.New の第二引数に dsnTest を直接使用
	mi, err := migrate.New(migrationURL, dsnTest)
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
	// ロガー設定なしのシンプルな GORM 設定
	db, err := gorm.Open(postgres.Open(dsnTest), &gorm.Config{})
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
