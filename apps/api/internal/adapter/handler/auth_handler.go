package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rentoutdoor/api/internal/adapter/dto/request"
	"github.com/rentoutdoor/api/internal/adapter/middleware"
	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/usecase"
	"github.com/rentoutdoor/api/pkg/response"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authUC usecase.AuthUsecase
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authUC usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new account with email and password. Roles: renter, owner.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.RegisterRequest true "Register payload"
// @Success 201 {object} response.Response{data=usecase.AuthOutput}
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	role := entity.UserRoleRenter
	if req.Role == "owner" {
		role = entity.UserRoleOwner
	}

	input := &usecase.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Phone:    req.Phone,
		Role:     role,
	}

	result, err := h.authUC.Register(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "registration successful", result)
}

// Login godoc
// @Summary Login with email and password
// @Description Authenticate with email/password credentials. Returns JWT access token and refresh token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.LoginRequest true "Login payload"
// @Success 200 {object} response.Response{data=usecase.AuthOutput}
// @Failure 401 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	input := &usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.authUC.Login(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "login successful", result)
}

// GoogleLogin godoc
// @Summary Login with Google
// @Description Authenticate with Google ID token. Creates account if first login.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.GoogleLoginRequest true "Google login payload"
// @Success 200 {object} response.Response{data=usecase.AuthOutput}
// @Failure 401 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /auth/google [post]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req request.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	input := &usecase.GoogleLoginInput{
		IDToken: req.IDToken,
	}

	result, err := h.authUC.GoogleLogin(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "google login successful", result)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Exchange a valid refresh token for a new access token and refresh token pair (rotation).
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.RefreshTokenRequest true "Refresh token payload"
// @Success 200 {object} response.Response{data=usecase.AuthOutput}
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, err := h.authUC.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "token refreshed", result)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send a password reset token for the given email. Does not reveal whether email exists.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.ForgotPasswordRequest true "Forgot password payload"
// @Success 200 {object} response.Response{data=usecase.ForgotPasswordOutput}
// @Failure 422 {object} response.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, err := h.authUC.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "password reset requested", result)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using a valid reset token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body request.ResetPasswordRequest true "Reset password payload"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := h.authUC.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "password reset successful", nil)
}

// Logout godoc
// @Summary Logout and revoke session
// @Description Revoke the current session's refresh token.
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, _ := c.Get(middleware.ContextKeySessionID)

	if err := h.authUC.Logout(c.Request.Context(), sessionID.(string)); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "logged out successfully", nil)
}
