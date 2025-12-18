package main

import (
	"log"
	"os"

	"backend/internal/infrastructure/persistence"
	"backend/internal/presentation/handler"
	"backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// --- データベース接続の初期化 ---
	// 環境変数からDB接続情報を取得
	dsnAuth := os.Getenv("DSN_AUTH")
	if dsnAuth == "" {
		log.Fatal("environment variable DSN_AUTH is not set")
	}

	db, err := gorm.Open(postgres.Open(dsnAuth), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// --- 依存関係の注入 (Dependency Injection) ---
	deviceRepo := persistence.NewDeviceGormRepository(db)
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo)
	deviceHandler := handler.NewDeviceHandler(deviceUsecase)

	// --- Ginルーターのセットアップ ---
	router := gin.Default()

	// デバイス関連のエンドポイントをグループ化
	deviceRoutes := router.Group("/devices")
	{
		deviceRoutes.POST("", deviceHandler.CreateDevice)
		deviceRoutes.GET("", deviceHandler.ListDevices)
		deviceRoutes.GET("/:id", deviceHandler.GetDevice)
		deviceRoutes.PUT("/:id", deviceHandler.UpdateDevice)
		deviceRoutes.DELETE("/:id", deviceHandler.DeleteDevice)
	}

	// --- サーバーの起動 ---
	log.Println("Server starting on port 8080...")

	if err := router.Run(":8080"); err != nil { //nolint:noinlineerr
		log.Fatalf("failed to start server: %v", err)
	}
}
