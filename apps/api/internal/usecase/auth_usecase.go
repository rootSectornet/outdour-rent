package usecase

import (
	"context"

	"github.com/rentoutdoor/api/internal/domain/entity"
)

// AuthUsecase defines the interface for authentication business logic.
type AuthUsecase interface {
	Register(ctx context.Context, input *RegisterInput) (*AuthOutput, error)
	Login(ctx context.Context, input *LoginInput) (*AuthOutput, error)
	GoogleLogin(ctx context.Context, input *GoogleLoginInput) (*AuthOutput, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthOutput, error)
	ForgotPassword(ctx context.Context, email string) (*ForgotPasswordOutput, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
	Logout(ctx context.Context, sessionID string) error
}

type RegisterInput struct {
	Email    string
	Password string
	FullName string
	Phone    string
	Role     entity.UserRole
}

type LoginInput struct {
	Email    string
	Password string
}

type GoogleLoginInput struct {
	IDToken string
}

type AuthOutput struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"`
	User         *entity.User `json:"user"`
}

type ForgotPasswordOutput struct {
	Message string `json:"message"`
	// In production, token is sent via email. Exposed here for dev/testing only.
	ResetToken string `json:"reset_token,omitempty"`
}
