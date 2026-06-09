package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/infrastructure/config"
	"github.com/rentoutdoor/api/internal/repository"
	"github.com/rentoutdoor/api/internal/usecase"
	"github.com/rentoutdoor/api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authUsecase struct {
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	resetRepo    repository.PasswordResetRepository
	cfg          config.Config
	jwtCfg       config.JWTConfig
	googleVerify GoogleTokenVerifier
}

// GoogleTokenVerifier verifies Google ID tokens.
type GoogleTokenVerifier interface {
	Verify(ctx context.Context, idToken string) (*GoogleClaims, error)
}

// GoogleClaims holds claims from a verified Google ID token.
type GoogleClaims struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// NewAuthUsecase creates a new auth usecase instance.
func NewAuthUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	resetRepo repository.PasswordResetRepository,
	cfg config.Config,
	jwtCfg config.JWTConfig,
	googleVerify GoogleTokenVerifier,
) usecase.AuthUsecase {
	return &authUsecase{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		resetRepo:    resetRepo,
		cfg:          cfg,
		jwtCfg:       jwtCfg,
		googleVerify: googleVerify,
	}
}

func (uc *authUsecase) Register(ctx context.Context, input *usecase.RegisterInput) (*usecase.AuthOutput, error) {
	// Check if user already exists
	existing, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: email already registered", usecase.ErrConflict)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	passwordHash := string(hashedPassword)

	// Validate role
	role := input.Role
	if role == "" {
		role = entity.UserRoleRenter
	}
	if role == entity.UserRoleAdmin {
		return nil, fmt.Errorf("%w: cannot self-register as admin", usecase.ErrForbidden)
	}

	user := &entity.User{
		Email:        input.Email,
		PasswordHash: &passwordHash,
		FullName:     input.FullName,
		Role:         role,
		IsActive:     true,
		Provider:     "local",
	}
	if input.Phone != "" {
		user.Phone = &input.Phone
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return uc.generateTokenPair(ctx, user, "")
}

func (uc *authUsecase) RegisterInvitation(ctx context.Context, input *usecase.RegisterInvitationInput) error {
	existing, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if existing != nil {
		return fmt.Errorf("%w: email already registered", usecase.ErrConflict)
	}

	user := &entity.User{
		Email:    input.Email,
		FullName: input.FullName,
		Role:     entity.UserRoleRenter,
		IsActive: true,
		Provider: "local",
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	rawToken, tokenHash, err := utils.GenerateResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	reset := &entity.PasswordReset{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := uc.resetRepo.Create(ctx, reset); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", uc.cfg.App.FrontEndURL, rawToken)

	body := fmt.Sprintf(`
Welcome!

Your account has been created.

Please set your password:

%s

This link expires in 24 hours.
`, resetLink)

	if err := utils.Send(ctx, user.Email, "Set Your Password", body, uc.cfg.Email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (uc *authUsecase) Login(ctx context.Context, input *usecase.LoginInput) (*usecase.AuthOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: invalid credentials", usecase.ErrUnauthorized)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("%w: account is deactivated", usecase.ErrForbidden)
	}

	if user.PasswordHash == nil {
		return nil, fmt.Errorf("%w: please login with Google", usecase.ErrUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("%w: invalid credentials", usecase.ErrUnauthorized)
	}

	return uc.generateTokenPair(ctx, user, "")
}

func (uc *authUsecase) GoogleLogin(ctx context.Context, input *usecase.GoogleLoginInput) (*usecase.AuthOutput, error) {
	claims, err := uc.googleVerify.Verify(ctx, input.IDToken)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid Google token", usecase.ErrUnauthorized)
	}

	if !claims.EmailVerified {
		return nil, fmt.Errorf("%w: Google email not verified", usecase.ErrUnauthorized)
	}

	// Try to find existing user by Google ID
	user, err := uc.userRepo.FindByGoogleID(ctx, claims.Sub)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to find user by google id: %w", err)
	}

	if user != nil {
		if !user.IsActive {
			return nil, fmt.Errorf("%w: account is deactivated", usecase.ErrForbidden)
		}
		return uc.generateTokenPair(ctx, user, "")
	}

	// Try to find user by email and link Google account
	user, err = uc.userRepo.FindByEmail(ctx, claims.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if user != nil {
		// Link Google account to existing user
		user.GoogleID = &claims.Sub
		user.Provider = "google"
		if user.AvatarURL == nil && claims.Picture != "" {
			user.AvatarURL = &claims.Picture
		}
		now := time.Now()
		if user.EmailVerifiedAt == nil {
			user.EmailVerifiedAt = &now
		}
		if err := uc.userRepo.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to link google account: %w", err)
		}
		return uc.generateTokenPair(ctx, user, "")
	}

	// Create new user from Google claims
	now := time.Now()
	user = &entity.User{
		Email:           claims.Email,
		FullName:        claims.Name,
		GoogleID:        &claims.Sub,
		Provider:        "google",
		Role:            entity.UserRoleRenter,
		IsActive:        true,
		EmailVerifiedAt: &now,
	}
	if claims.Picture != "" {
		user.AvatarURL = &claims.Picture
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create google user: %w", err)
	}

	return uc.generateTokenPair(ctx, user, "")
}

func (uc *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*usecase.AuthOutput, error) {
	tokenHash := hashToken(refreshToken)

	session, err := uc.sessionRepo.FindActiveByRefreshTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: invalid refresh token", usecase.ErrUnauthorized)
		}
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	if session.RevokedAt != nil {
		// Token reuse detected - revoke all sessions for this user (security measure)
		_ = uc.sessionRepo.RevokeByUserID(ctx, session.UserID)
		return nil, fmt.Errorf("%w: refresh token already used", usecase.ErrUnauthorized)
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("%w: refresh token expired", usecase.ErrUnauthorized)
	}

	// Revoke the old session (rotation)
	if err := uc.sessionRepo.Revoke(ctx, session.ID); err != nil {
		return nil, fmt.Errorf("failed to revoke old session: %w", err)
	}

	// Get the user
	user, err := uc.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("%w: account is deactivated", usecase.ErrForbidden)
	}

	return uc.generateTokenPair(ctx, user, "")
}

func (uc *authUsecase) ForgotPassword(ctx context.Context, email string) (*usecase.ForgotPasswordOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal whether email exists
			return &usecase.ForgotPasswordOutput{
				Message: "if the email exists, a reset link has been sent",
			}, nil
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Delete any existing reset tokens
	_ = uc.resetRepo.DeleteByUserID(ctx, user.ID)

	// Generate a secure random token
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate reset token: %w", err)
	}

	reset := &entity.PasswordReset{
		UserID:    user.ID,
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	if err := uc.resetRepo.Create(ctx, reset); err != nil {
		return nil, fmt.Errorf("failed to create password reset: %w", err)
	}

	// TODO: In production, send token via email service.
	// For development, include the token in the response.
	return &usecase.ForgotPasswordOutput{
		Message:    "if the email exists, a reset link has been sent",
		ResetToken: token,
	}, nil
}

func (uc *authUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	tokenHash := hashToken(token)

	reset, err := uc.resetRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: invalid or expired reset token", usecase.ErrUnauthorized)
		}
		return fmt.Errorf("failed to find reset token: %w", err)
	}

	if reset.UsedAt != nil {
		return fmt.Errorf("%w: reset token already used", usecase.ErrUnauthorized)
	}

	if time.Now().After(reset.ExpiresAt) {
		return fmt.Errorf("%w: reset token expired", usecase.ErrUnauthorized)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	user, err := uc.userRepo.FindByID(ctx, reset.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	hash := string(hashedPassword)
	user.PasswordHash = &hash
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark reset token as used
	if err := uc.resetRepo.MarkUsed(ctx, reset.ID); err != nil {
		return fmt.Errorf("failed to mark reset token: %w", err)
	}

	// Revoke all existing sessions (force re-login)
	_ = uc.sessionRepo.RevokeByUserID(ctx, user.ID)

	return nil
}

func (uc *authUsecase) Logout(ctx context.Context, sessionID string) error {
	if err := uc.sessionRepo.Revoke(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}
	return nil
}

// --- Private helpers ---

func (uc *authUsecase) generateTokenPair(ctx context.Context, user *entity.User, userAgent string) (*usecase.AuthOutput, error) {
	now := time.Now()

	// Generate refresh token (opaque random string)
	refreshTokenStr, err := generateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store session with hashed refresh token
	session := &entity.UserSession{
		UserID:           user.ID,
		RefreshTokenHash: hashToken(refreshTokenStr),
		ExpiresAt:        now.Add(uc.jwtCfg.RefreshExpiry),
	}
	if userAgent != "" {
		session.UserAgent = &userAgent
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Generate access token with session ID
	accessClaims := jwt.MapClaims{
		"sub":  user.ID,
		"role": string(user.Role),
		"sid":  session.ID,
		"iat":  now.Unix(),
		"exp":  now.Add(uc.jwtCfg.AccessExpiry).Unix(),
		"type": "access",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(uc.jwtCfg.AccessSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	return &usecase.AuthOutput{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		TokenType:    "Bearer",
		ExpiresIn:    int64(uc.jwtCfg.AccessExpiry.Seconds()),
		User:         user,
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
