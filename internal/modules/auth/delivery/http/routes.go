package http

import (
	"github.com/JscorpTech/auth/internal/config"
	"github.com/JscorpTech/auth/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(cfg *config.Config, router *gin.RouterGroup, h *AuthHandler) {
	public := router.Group("/auth")
	{
		public.POST("/login", h.Login)
		public.POST("/register", h.Register)
	}
	private := router.Group("/auth")
	private.Use(middlewares.AuthMiddleware(cfg))
	{
		private.GET("/me", h.Me)
	}
}
