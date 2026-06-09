package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rentoutdoor/api/internal/adapter/handler"
	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/usecase"
)

// --- Mock AuthUsecase ---

type mockAuthUsecase struct {
	registerFn           func(ctx context.Context, input *usecase.RegisterInput) (*usecase.AuthOutput, error)
	registerInvitationFn func(ctx context.Context, input *usecase.RegisterInvitationInput) error
	loginFn              func(ctx context.Context, input *usecase.LoginInput) (*usecase.AuthOutput, error)
	googleLoginFn        func(ctx context.Context, input *usecase.GoogleLoginInput) (*usecase.AuthOutput, error)
	refreshTokenFn       func(ctx context.Context, refreshToken string) (*usecase.AuthOutput, error)
	forgotPasswordFn     func(ctx context.Context, email string) (*usecase.ForgotPasswordOutput, error)
	resetPasswordFn      func(ctx context.Context, token, newPassword string) error
	logoutFn             func(ctx context.Context, sessionID string) error
}

func (m *mockAuthUsecase) Register(ctx context.Context, input *usecase.RegisterInput) (*usecase.AuthOutput, error) {
	return m.registerFn(ctx, input)
}

func (m *mockAuthUsecase) RegisterInvitation(ctx context.Context, input *usecase.RegisterInvitationInput) error {
	return m.registerInvitationFn(ctx, input)
}

func (m *mockAuthUsecase) Login(ctx context.Context, input *usecase.LoginInput) (*usecase.AuthOutput, error) {
	return m.loginFn(ctx, input)
}

func (m *mockAuthUsecase) GoogleLogin(ctx context.Context, input *usecase.GoogleLoginInput) (*usecase.AuthOutput, error) {
	return m.googleLoginFn(ctx, input)
}

func (m *mockAuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*usecase.AuthOutput, error) {
	return m.refreshTokenFn(ctx, refreshToken)
}

func (m *mockAuthUsecase) ForgotPassword(ctx context.Context, email string) (*usecase.ForgotPasswordOutput, error) {
	return m.forgotPasswordFn(ctx, email)
}

func (m *mockAuthUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	return m.resetPasswordFn(ctx, token, newPassword)
}

func (m *mockAuthUsecase) Logout(ctx context.Context, sessionID string) error {
	return m.logoutFn(ctx, sessionID)
}

// --- Helpers ---

func setupRouter(authUC usecase.AuthUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewAuthHandler(authUC)

	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/google", h.GoogleLogin)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
	}
	return r
}

func successAuthOutput() *usecase.AuthOutput {
	return &usecase.AuthOutput{
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-456",
		TokenType:    "Bearer",
		ExpiresIn:    900,
		User: &entity.User{
			BaseModel: entity.BaseModel{ID: "user-1"},
			Email:     "test@example.com",
			FullName:  "Test User",
			Role:      entity.UserRoleRenter,
			IsActive:  true,
		},
	}
}

// --- Tests ---

func TestHandler_Register_Success(t *testing.T) {
	mock := &mockAuthUsecase{
		registerFn: func(_ context.Context, input *usecase.RegisterInput) (*usecase.AuthOutput, error) {
			if input.Email != "test@example.com" {
				t.Fatalf("expected email test@example.com, got %s", input.Email)
			}
			return successAuthOutput(), nil
		},
	}

	r := setupRouter(mock)
	body := map[string]string{
		"email":     "test@example.com",
		"password":  "password123",
		"full_name": "Test User",
	}

	w := performRequest(r, "POST", "/api/v1/auth/register", body)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Register_ValidationError(t *testing.T) {
	mock := &mockAuthUsecase{}
	r := setupRouter(mock)

	// Missing required fields
	body := map[string]string{
		"email": "not-an-email",
	}

	w := performRequest(r, "POST", "/api/v1/auth/register", body)

	if w.Code != http.StatusUnprocessableEntity && w.Code != http.StatusBadRequest {
		t.Fatalf("expected 422 or 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Register_Conflict(t *testing.T) {
	mock := &mockAuthUsecase{
		registerFn: func(_ context.Context, _ *usecase.RegisterInput) (*usecase.AuthOutput, error) {
			return nil, usecase.ErrConflict
		},
	}

	r := setupRouter(mock)
	body := map[string]string{
		"email":     "existing@example.com",
		"password":  "password123",
		"full_name": "Existing User",
	}

	w := performRequest(r, "POST", "/api/v1/auth/register", body)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Login_Success(t *testing.T) {
	mock := &mockAuthUsecase{
		loginFn: func(_ context.Context, input *usecase.LoginInput) (*usecase.AuthOutput, error) {
			if input.Email != "user@example.com" {
				t.Fatalf("expected email user@example.com, got %s", input.Email)
			}
			return successAuthOutput(), nil
		},
	}

	r := setupRouter(mock)
	body := map[string]string{
		"email":    "user@example.com",
		"password": "password123",
	}

	w := performRequest(r, "POST", "/api/v1/auth/login", body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != true {
		t.Fatal("expected success to be true")
	}
}

func TestHandler_Login_Unauthorized(t *testing.T) {
	mock := &mockAuthUsecase{
		loginFn: func(_ context.Context, _ *usecase.LoginInput) (*usecase.AuthOutput, error) {
			return nil, usecase.ErrUnauthorized
		},
	}

	r := setupRouter(mock)
	body := map[string]string{
		"email":    "user@example.com",
		"password": "wrong",
	}

	w := performRequest(r, "POST", "/api/v1/auth/login", body)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GoogleLogin_Success(t *testing.T) {
	mock := &mockAuthUsecase{
		googleLoginFn: func(_ context.Context, input *usecase.GoogleLoginInput) (*usecase.AuthOutput, error) {
			if input.IDToken != "google-id-token" {
				t.Fatalf("expected id_token google-id-token, got %s", input.IDToken)
			}
			return successAuthOutput(), nil
		},
	}

	r := setupRouter(mock)
	body := map[string]string{
		"id_token": "google-id-token",
	}

	w := performRequest(r, "POST", "/api/v1/auth/google", body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_RefreshToken_Success(t *testing.T) {
	mock := &mockAuthUsecase{
		refreshTokenFn: func(_ context.Context, token string) (*usecase.AuthOutput, error) {
			if token != "my-refresh-token" {
				t.Fatalf("expected my-refresh-token, got %s", token)
			}
			return successAuthOutput(), nil
		},
	}

	r := setupRouter(mock)
	body := map[string]string{
		"refresh_token": "my-refresh-token",
	}

	w := performRequest(r, "POST", "/api/v1/auth/refresh", body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ForgotPassword_Success(t *testing.T) {
	mock := &mockAuthUsecase{
		forgotPasswordFn: func(_ context.Context, email string) (*usecase.ForgotPasswordOutput, error) {
			return &usecase.ForgotPasswordOutput{
				Message:    "if the email exists, a reset link has been sent",
				ResetToken: "reset-token-123",
			}, nil
		},
	}

	r := setupRouter(mock)
	body := map[string]string{
		"email": "user@example.com",
	}

	w := performRequest(r, "POST", "/api/v1/auth/forgot-password", body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ResetPassword_Success(t *testing.T) {
	mock := &mockAuthUsecase{
		resetPasswordFn: func(_ context.Context, token, newPassword string) error {
			if token != "valid-token" {
				t.Fatalf("expected valid-token, got %s", token)
			}
			return nil
		},
	}

	r := setupRouter(mock)
	body := map[string]string{
		"token":        "valid-token",
		"new_password": "newpass123",
	}

	w := performRequest(r, "POST", "/api/v1/auth/reset-password", body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ResetPassword_InvalidToken(t *testing.T) {
	mock := &mockAuthUsecase{
		resetPasswordFn: func(_ context.Context, _, _ string) error {
			return usecase.ErrUnauthorized
		},
	}

	r := setupRouter(mock)
	body := map[string]string{
		"token":        "invalid-token",
		"new_password": "newpass123",
	}

	w := performRequest(r, "POST", "/api/v1/auth/reset-password", body)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}
