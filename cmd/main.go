package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/JscorpTech/auth/internal/config"
	"github.com/JscorpTech/auth/internal/modules/auth"
	authHttp "github.com/JscorpTech/auth/internal/modules/auth/delivery/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(".env load error")
	}

	logger, _ := zap.NewDevelopment()
	cfg := config.NewConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("failed to connect databse")
	}

	// migrations
	db.AutoMigrate(&auth.User{})

	router := gin.Default()
	api := router.Group("/api")

	// Auth routes
	authRepository := auth.NewAuthRepository(db)
	authUsecase := auth.NewAuthUsecase(authRepository, cfg)
	authHandler := authHttp.NewAuthHandler(authUsecase, logger)
	authHttp.RegisterAuthRoutes(cfg, api, authHandler)

	srv := http.Server{
		Handler: router,
		Addr:    cfg.Addr,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		logger.Info("Server ishga tushdi 🚀 " + cfg.Addr)
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	<-stop
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Serverni o'chirishda xatolik yuz berdi", zap.Error(err))
	}
	logger.Info("Server muvaffaqiyatli to'xtatildi")
}
