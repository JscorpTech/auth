package http

import (
	"net/http"

	"github.com/JscorpTech/auth/internal/dto"
	"github.com/JscorpTech/auth/internal/modules/auth"
	"github.com/JscorpTech/auth/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// @Router /api/auth/google [post]
// @Accept json
// @Produce json
// @Tags auth
// @Summary Google authentication
// @Param request body auth.GoogleAuthRequest true "Google auth request"
// @Success 200 {object} dto.BaseResponse{data=auth.AuthLoginResponse}
func (h *AuthHandler) Google(c *gin.Context) {
	var payload auth.GoogleAuthRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		dto.JSON(c, http.StatusBadRequest, utils.FormatValidationErrors(err, &payload), "Invalid request")
		return
	}
	user, err := h.usecase.GoogleAuth(ctx, payload.IDToken)
	if err != nil {
		dto.JSON(c, http.StatusUnauthorized, nil, "Invalid Google ID token")
		return
	}
	dto.JSON(c, http.StatusOK, auth.AuthLoginResponse{
		Token: auth.ToToken(h.usecase.AccessToken(user), h.usecase.RefreshToken(user)),
		User:  auth.ToUser(user),
	}, "")
}

// @Router /api/auth/refresh [post]
// @Accept json
// @Produce json
// @Tags auth
// @Summary Refresh token
// @Success 200 {object} dto.BaseResponse{data=auth.TokenDTO}
// @Param request body auth.AuthRefreshTokenRequest true "Refresh token request"
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var payload auth.AuthRefreshTokenRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		dto.JSON(c, http.StatusBadRequest, utils.FormatValidationErrors(err, &payload), "Invalid request")
		return
	}
	claims, err := h.usecase.ValidateToken(payload.RefreshToken)
	if err != nil {
		dto.JSON(c, http.StatusUnauthorized, nil, "Invalid refresh token")
		return
	}

	user := &auth.User{
		Model: gorm.Model{
			ID: uint((*claims)["user_id"].(float64)),
		},
		Role: (*claims)["role"].(auth.Role),
	}

	dto.JSON(c, http.StatusOK, auth.ToToken(h.usecase.AccessToken(user), ""), "")
}

// Register godoc
// @Summary Login user
// @Router /api/auth/login [post]
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=auth.AuthLoginResponse}
// @Param request body auth.AuthLoginRequest true "Login form"
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var payload auth.AuthLoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		dto.JSON(c, http.StatusBadRequest, utils.FormatValidationErrors(err, &payload), "Invalid request")
		return
	}
	user, err := h.usecase.Login(ctx, payload.Phone, payload.Password)
	if err != nil {
		dto.JSON(c, http.StatusUnauthorized, nil, err.Error())
		return
	}
	dto.JSON(c, http.StatusOK, auth.AuthLoginResponse{
		Token: auth.ToToken(h.usecase.AccessToken(user), h.usecase.RefreshToken(user)),
		User:  auth.ToUser(user),
	}, "")
}

// Register godoc
// @Summary Get user profile
// @Router /api/auth/me [get]
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.BaseResponse{data=auth.AuthMeResponse}
func (h *AuthHandler) Me(c *gin.Context) {
	user := c.MustGet("user").(jwt.MapClaims)
	dto.JSON(c, http.StatusOK, auth.AuthMeResponse{
		User: user,
	}, "")
}

// Register godoc
// @Summary Register user
// @Description Create new user
// @Tags auth
// @Accept json
// @Produce json
// @Router /api/auth/register [post]
// @Success 200 {object} dto.BaseResponse{data=auth.AuthRegisterResponse}
// @Param request body auth.AuthRegisterRequest true "Register request"
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var userPayload auth.AuthRegisterRequest
	if err := c.ShouldBindJSON(&userPayload); err != nil {
		dto.JSON(c, http.StatusBadRequest, utils.FormatValidationErrors(err, &userPayload), "Invalid request")
		return
	}
	password, err := utils.HashPassword(userPayload.Password)
	if err != nil {
		dto.JSON(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}
	userModel := auth.User{
		FirstName: userPayload.FirstName,
		LastName:  userPayload.LastName,
		Phone:     userPayload.Phone,
		Password:  password,
	}
	user, err := h.usecase.Register(ctx, &userModel)
	if err != nil {
		dto.JSON(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	dto.JSON(
		c, http.StatusOK,
		auth.ToRegisterResponse(user, "Tasdiqlash ko'di yuborildi"),
		"",
	)
}

// @Router /api/auth/confirm [post]
// @Summary Confirm phone number
// @Accept json
// @Produce json
// @Param request body auth.AuthConfirmRequest true "Confirm request"
// @Tags auth
// @Success 200 {object} dto.BaseResponse{data=auth.TokenDTO}
func (h *AuthHandler) Confirm(c *gin.Context) {
	var payload auth.AuthConfirmRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		dto.JSON(c, http.StatusBadRequest, utils.FormatValidationErrors(err, &payload), "Validation error")
		return
	}
	if isValid := h.usecase.ValidateOtp(ctx, payload.Phone, payload.Otp); !isValid {
		dto.JSON(c, http.StatusForbidden, nil, "Invalid otp")
		return
	}
	user, err := h.usecase.GetUserByPhone(ctx, payload.Phone)
	if err != nil {
		dto.JSON(c, http.StatusBadRequest, nil, "Invalid phone number")
		return
	}
	h.usecase.Confirm(ctx, user)
	dto.JSON(c, 200, auth.ToToken(h.usecase.AccessToken(user), h.usecase.RefreshToken(user)), "")
}
