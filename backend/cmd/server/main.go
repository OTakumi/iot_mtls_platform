package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/infrastructure/persistence"
	"backend/internal/presentation/handler"
	"backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	defaultServerPort       = ":8080"
	dbMaxIdleConns          = 10
	dbMaxOpenConns          = 100
	serverShutdownTimeout   = 5 * time.Second
	serverReadHeaderTimeout = 10 * time.Second
)

func main() {
	// --- データベース接続の初期化 ---
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

	// --- データベース接続プールの設定 ---
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(dbMaxIdleConns)
	sqlDB.SetMaxOpenConns(dbMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// --- 依存関係の注入 (Dependency Injection) ---
	deviceRepo := persistence.NewDeviceGormRepository(db)
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo)
	deviceHandler := handler.NewDeviceHandler(deviceUsecase)

	// --- Ginルーターのセットアップ ---
	router := gin.Default()

	// ヘルスチェック用エンドポイント
	router.GET("/health", func(c *gin.Context) {
		err := sqlDB.PingContext(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})

			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// デバイス関連のエンドポイントをグループ化
	deviceRoutes := router.Group("/devices")
	{
		deviceRoutes.POST("", deviceHandler.CreateDevice)
		deviceRoutes.GET("", deviceHandler.ListDevices)
		deviceRoutes.GET("/:id", deviceHandler.GetDevice)
		deviceRoutes.PUT("/:id", deviceHandler.UpdateDevice)
		deviceRoutes.DELETE("/:id", deviceHandler.DeleteDevice)
	}

	// --- サーバーのグレースフルシャットダウン ---
	srv := &http.Server{ //nolint:exhaustruct
		Addr:              defaultServerPort,
		Handler:           router,
		ReadHeaderTimeout: serverReadHeaderTimeout,
	}

	// サーバーをゴルーチンで起動
	go func() {
		log.Println("Server starting on port", defaultServerPort, "...")

		err = srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// OSのシグナルを待機 (Graceful shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// タイムアウト付きコンテキストでシャットダウン処理
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel() // defer ensures cancel is called

	err = srv.Shutdown(ctx) // noinlineerr: avoid inline error handling
	if err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		// Explicitly call cancel before os.Exit to ensure deferred cleanup runs
		cancel()
	}

	log.Println("Server exiting")
}
