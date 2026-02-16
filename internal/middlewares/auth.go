package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/JscorpTech/auth/internal/config"
	"github.com/JscorpTech/auth/internal/dto"
	"github.com/JscorpTech/auth/pkg/utils"
	"github.com/gin-gonic/gin"
)

func RateLimiterPerIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := utils.GetVisitor(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status": false,
				"error":  "Too many requests from your IP",
			})
			return
		}

		c.Next()
	}
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenRaw := c.Request.Header.Get("Authorization")
		token := strings.Replace(tokenRaw, "Bearer ", "", 1)
		claims, err := utils.VerifyJWT(token, cfg.PublicKey)

		exp, ok := claims["exp"]
		if !ok {
			dto.JSON(c, http.StatusUnauthorized, nil, "Invalid token")
			c.Abort()
			return
		}
		if exp.(float64) < float64(time.Now().Unix()) {
			dto.JSON(c, http.StatusUnauthorized, nil, "Token expired")
			c.Abort()
			return
		}
		if err != nil || claims["type"] != "access" {
			dto.JSON(c, http.StatusUnauthorized, nil, "Invalid token")
			c.Abort()
			return
		}
		c.Set("user", claims)
		c.Next()
	}
}
