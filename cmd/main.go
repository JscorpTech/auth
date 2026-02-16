package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/JscorpTech/auth/docs"
	"github.com/JscorpTech/auth/internal/config"
	"github.com/JscorpTech/auth/internal/middlewares"
	"github.com/JscorpTech/auth/internal/modules/auth"
	authHttp "github.com/JscorpTech/auth/internal/modules/auth/delivery/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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
	db.AutoMigrate(&auth.Otp{})

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	api.Use(middlewares.RateLimiterPerIP())

	// Auth routes
	authRepository := auth.NewAuthRepository(db)
	authUsecase := auth.NewAuthUsecase(authRepository, cfg, logger)
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
