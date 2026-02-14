package http

import (
	"net/http"

	"github.com/JscorpTech/auth/internal/modules/auth"
	"github.com/JscorpTech/auth/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var userPayload auth.AuthRegisterRequest
	if err := c.ShouldBindJSON(&userPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": utils.FormatValidationErrors(err, &userPayload)})
		return
	}
	userModel := auth.User{
		FirstName: userPayload.FirstName,
		LastName:  userPayload.LastName,
		Email:     userPayload.Email,
		Phone:     userPayload.Phone,
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
	})
}
