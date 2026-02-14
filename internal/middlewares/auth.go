package middlewares

import (
	"strings"
	"time"

	"github.com/JscorpTech/auth/internal/config"
	"github.com/JscorpTech/auth/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenRaw := c.Request.Header.Get("Authorization")
		token := strings.Replace(tokenRaw, "Bearer ", "", 1)
		claims, err := utils.VerifyJWT(token, cfg.PublicKey)
		exp := claims["exp"].(float64)
		if exp < float64(time.Now().Unix()) {
			c.JSON(401, gin.H{
				"error": "Token expired",
			})
			c.Abort()
			return
		}
		if err != nil || claims["type"] != "access" {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Set("user", claims)
		c.Next()
	}
}
