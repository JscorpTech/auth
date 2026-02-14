package http

import (
	"net/http"

	"github.com/JscorpTech/auth/internal/modules/auth"
	"github.com/JscorpTech/auth/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthHandler struct {
	usecase auth.AuthUsecase
	logger  *zap.Logger
}

func NewAuthHandler(usecase auth.AuthUsecase, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var payload auth.AuthRefreshTokenRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": utils.FormatValidationErrors(err, &payload),
		})
		return
	}
	claims, err := h.usecase.ValidateToken(payload.RefreshToken)
	if err != nil || (*claims)["type"] != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	user := &auth.User{
		Model: gorm.Model{
			ID: uint((*claims)["user_id"].(float64)),
		},
		Role: (*claims)["role"].(string),
	}

	c.JSON(http.StatusOK, gin.H{
		"access":  h.usecase.AccessToken(user),
		"refresh": h.usecase.RefreshToken(user),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var payload auth.AuthLoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": utils.FormatValidationErrors(err, &payload)})
		return
	}
	user, err := h.usecase.Login(ctx, payload.Phone, payload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access":  h.usecase.AccessToken(user),
		"refresh": h.usecase.RefreshToken(user),
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user": c.MustGet("user"),
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var userPayload auth.AuthRegisterRequest
	if err := c.ShouldBindJSON(&userPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": utils.FormatValidationErrors(err, &userPayload)})
		return
	}
	password, err := utils.HashPassword(userPayload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	userModel := auth.User{
		FirstName: userPayload.FirstName,
		LastName:  userPayload.LastName,
		Email:     userPayload.Email,
		Phone:     userPayload.Phone,
		Password:  password,
	}
	user, err := h.usecase.Register(ctx, &userModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user":    auth.ToRegisterResponse(user),
		"token": map[string]string{
			"access":  h.usecase.AccessToken(user),
			"refresh": h.usecase.RefreshToken(user),
		},
	})
}
