package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/infrastructure/config"
	"github.com/rentoutdoor/api/internal/usecase"
	"github.com/rentoutdoor/api/internal/usecase/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// --- Mocks ---

type mockUserRepo struct {
	users    map[string]*entity.User
	byEmail  map[string]*entity.User
	byGoogle map[string]*entity.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:    make(map[string]*entity.User),
		byEmail:  make(map[string]*entity.User),
		byGoogle: make(map[string]*entity.User),
	}
}

func (m *mockUserRepo) Create(_ context.Context, user *entity.User) error {
	if _, exists := m.byEmail[user.Email]; exists {
		return errors.New("duplicate email")
	}
	m.users[user.ID] = user
	m.byEmail[user.Email] = user
	if user.GoogleID != nil {
		m.byGoogle[*user.GoogleID] = user
	}
	return nil
}

func (m *mockUserRepo) FindByID(_ context.Context, id string) (*entity.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	if user, ok := m.byEmail[email]; ok {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) FindByGoogleID(_ context.Context, googleID string) (*entity.User, error) {
	if user, ok := m.byGoogle[googleID]; ok {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) Update(_ context.Context, user *entity.User) error {
	m.users[user.ID] = user
	m.byEmail[user.Email] = user
	if user.GoogleID != nil {
		m.byGoogle[*user.GoogleID] = user
	}
	return nil
}

func (m *mockUserRepo) Delete(_ context.Context, id string, _ string) error {
	if user, ok := m.users[id]; ok {
		delete(m.byEmail, user.Email)
		delete(m.users, id)
		return nil
	}
	return gorm.ErrRecordNotFound
}

type mockSessionRepo struct {
	sessions map[string]*entity.UserSession
	byToken  map[string]*entity.UserSession
}

func newMockSessionRepo() *mockSessionRepo {
	return &mockSessionRepo{
		sessions: make(map[string]*entity.UserSession),
		byToken:  make(map[string]*entity.UserSession),
	}
}

func (m *mockSessionRepo) Create(_ context.Context, session *entity.UserSession) error {
	m.sessions[session.ID] = session
	m.byToken[session.RefreshTokenHash] = session
	return nil
}

func (m *mockSessionRepo) FindByID(_ context.Context, id string) (*entity.UserSession, error) {
	if s, ok := m.sessions[id]; ok {
		return s, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockSessionRepo) FindActiveByRefreshTokenHash(_ context.Context, tokenHash string) (*entity.UserSession, error) {
	if s, ok := m.byToken[tokenHash]; ok {
		if s.RevokedAt == nil && time.Now().Before(s.ExpiresAt) {
			return s, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockSessionRepo) Revoke(_ context.Context, id string) error {
	if s, ok := m.sessions[id]; ok {
		now := time.Now()
		s.RevokedAt = &now
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockSessionRepo) RevokeByUserID(_ context.Context, userID string) error {
	now := time.Now()
	for _, s := range m.sessions {
		if s.UserID == userID {
			s.RevokedAt = &now
		}
	}
	return nil
}

func (m *mockSessionRepo) DeleteExpired(_ context.Context) error {
	for id, s := range m.sessions {
		if time.Now().After(s.ExpiresAt) {
			delete(m.sessions, id)
		}
	}
	return nil
}

type mockPasswordResetRepo struct {
	resets  map[string]*entity.PasswordReset
	byToken map[string]*entity.PasswordReset
}

func newMockPasswordResetRepo() *mockPasswordResetRepo {
	return &mockPasswordResetRepo{
		resets:  make(map[string]*entity.PasswordReset),
		byToken: make(map[string]*entity.PasswordReset),
	}
}

func (m *mockPasswordResetRepo) Create(_ context.Context, reset *entity.PasswordReset) error {
	m.resets[reset.ID] = reset
	m.byToken[reset.TokenHash] = reset
	return nil
}

func (m *mockPasswordResetRepo) FindByTokenHash(_ context.Context, tokenHash string) (*entity.PasswordReset, error) {
	if r, ok := m.byToken[tokenHash]; ok {
		if r.UsedAt == nil && time.Now().Before(r.ExpiresAt) {
			return r, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockPasswordResetRepo) MarkUsed(_ context.Context, id string) error {
	if r, ok := m.resets[id]; ok {
		now := time.Now()
		r.UsedAt = &now
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockPasswordResetRepo) DeleteByUserID(_ context.Context, userID string) error {
	for id, r := range m.resets {
		if r.UserID == userID {
			delete(m.byToken, r.TokenHash)
			delete(m.resets, id)
		}
	}
	return nil
}

type mockGoogleVerifier struct {
	claims *auth.GoogleClaims
	err    error
}

func (m *mockGoogleVerifier) Verify(_ context.Context, _ string) (*auth.GoogleClaims, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.claims, nil
}

// --- Helpers ---

func newTestJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		AccessSecret:  "test-access-secret",
		RefreshSecret: "test-refresh-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}
}

func setupAuthUsecase(googleClaims *auth.GoogleClaims, googleErr error) (
	usecase.AuthUsecase,
	*mockUserRepo,
	*mockSessionRepo,
	*mockPasswordResetRepo,
) {
	userRepo := newMockUserRepo()
	sessionRepo := newMockSessionRepo()
	resetRepo := newMockPasswordResetRepo()
	googleVerifier := &mockGoogleVerifier{claims: googleClaims, err: googleErr}

	uc := auth.NewAuthUsecase(userRepo, sessionRepo, resetRepo, newTestJWTConfig(), googleVerifier)
	return uc, userRepo, sessionRepo, resetRepo
}

func createTestUser(userRepo *mockUserRepo, email, password string, role entity.UserRole) *entity.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &entity.User{
		BaseModel: entity.BaseModel{ID: "user-123"},
		Email:     email,
		PasswordHash: string(hash),
		FullName:  "Test User",
		Role:      role,
		IsActive:  true,
		Provider:  "local",
	}
	userRepo.users[user.ID] = user
	userRepo.byEmail[user.Email] = user
	return user
}

// --- Tests ---

func TestRegister_Success(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	input := &usecase.RegisterInput{
		Email:    "new@example.com",
		Password: "password123",
		FullName: "New User",
		Phone:    "08123456789",
		Role:     entity.UserRoleRenter,
	}

	result, err := uc.Register(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if result.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}
	if result.TokenType != "Bearer" {
		t.Fatalf("expected token type Bearer, got %s", result.TokenType)
	}
	if result.User.Email != "new@example.com" {
		t.Fatalf("expected email new@example.com, got %s", result.User.Email)
	}
	if result.User.Role != entity.UserRoleRenter {
		t.Fatalf("expected role renter, got %s", result.User.Role)
	}
	if len(userRepo.users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(userRepo.users))
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	createTestUser(userRepo, "existing@example.com", "password123", entity.UserRoleRenter)

	input := &usecase.RegisterInput{
		Email:    "existing@example.com",
		Password: "password123",
		FullName: "Another User",
		Role:     entity.UserRoleRenter,
	}

	_, err := uc.Register(ctx, input)
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
	if !errors.Is(err, usecase.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestRegister_AdminRoleForbidden(t *testing.T) {
	uc, _, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	input := &usecase.RegisterInput{
		Email:    "admin@example.com",
		Password: "password123",
		FullName: "Admin User",
		Role:     entity.UserRoleAdmin,
	}

	_, err := uc.Register(ctx, input)
	if err == nil {
		t.Fatal("expected error for admin self-registration")
	}
	if !errors.Is(err, usecase.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	createTestUser(userRepo, "user@example.com", "correct-password", entity.UserRoleRenter)

	input := &usecase.LoginInput{
		Email:    "user@example.com",
		Password: "correct-password",
	}

	result, err := uc.Login(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if result.User.Email != "user@example.com" {
		t.Fatalf("expected email user@example.com, got %s", result.User.Email)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	createTestUser(userRepo, "user@example.com", "correct-password", entity.UserRoleRenter)

	input := &usecase.LoginInput{
		Email:    "user@example.com",
		Password: "wrong-password",
	}

	_, err := uc.Login(ctx, input)
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
	if !errors.Is(err, usecase.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	uc, _, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	input := &usecase.LoginInput{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, err := uc.Login(ctx, input)
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !errors.Is(err, usecase.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestLogin_InactiveUser(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	user := createTestUser(userRepo, "inactive@example.com", "password123", entity.UserRoleRenter)
	user.IsActive = false

	input := &usecase.LoginInput{
		Email:    "inactive@example.com",
		Password: "password123",
	}

	_, err := uc.Login(ctx, input)
	if err == nil {
		t.Fatal("expected error for inactive user")
	}
	if !errors.Is(err, usecase.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestGoogleLogin_NewUser(t *testing.T) {
	claims := &auth.GoogleClaims{
		Sub:           "google-123",
		Email:         "google@example.com",
		EmailVerified: true,
		Name:          "Google User",
		Picture:       "https://example.com/photo.jpg",
	}
	uc, userRepo, _, _ := setupAuthUsecase(claims, nil)
	ctx := context.Background()

	input := &usecase.GoogleLoginInput{IDToken: "valid-token"}
	result, err := uc.GoogleLogin(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if result.User.Email != "google@example.com" {
		t.Fatalf("expected email google@example.com, got %s", result.User.Email)
	}
	if result.User.GoogleID == nil || *result.User.GoogleID != "google-123" {
		t.Fatal("expected google ID to be set")
	}
	if len(userRepo.users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(userRepo.users))
	}
}

func TestGoogleLogin_ExistingGoogleUser(t *testing.T) {
	claims := &auth.GoogleClaims{
		Sub:           "google-123",
		Email:         "google@example.com",
		EmailVerified: true,
		Name:          "Google User",
		Picture:       "https://example.com/photo.jpg",
	}
	uc, userRepo, _, _ := setupAuthUsecase(claims, nil)
	ctx := context.Background()

	// Pre-create user with google ID
	googleID := "google-123"
	user := &entity.User{
		BaseModel: entity.BaseModel{ID: "user-456"},
		Email:     "google@example.com",
		FullName:  "Google User",
		GoogleID:  &googleID,
		Provider:  "google",
		Role:      entity.UserRoleRenter,
		IsActive:  true,
	}
	userRepo.users[user.ID] = user
	userRepo.byEmail[user.Email] = user
	userRepo.byGoogle[googleID] = user

	input := &usecase.GoogleLoginInput{IDToken: "valid-token"}
	result, err := uc.GoogleLogin(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.User.ID != "user-456" {
		t.Fatalf("expected user-456, got %s", result.User.ID)
	}
}

func TestGoogleLogin_LinkExistingEmail(t *testing.T) {
	claims := &auth.GoogleClaims{
		Sub:           "google-new",
		Email:         "existing@example.com",
		EmailVerified: true,
		Name:          "Existing User",
		Picture:       "https://example.com/photo.jpg",
	}
	uc, userRepo, _, _ := setupAuthUsecase(claims, nil)
	ctx := context.Background()

	createTestUser(userRepo, "existing@example.com", "password123", entity.UserRoleOwner)

	input := &usecase.GoogleLoginInput{IDToken: "valid-token"}
	result, err := uc.GoogleLogin(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.User.GoogleID == nil || *result.User.GoogleID != "google-new" {
		t.Fatal("expected google ID to be linked")
	}
	if result.User.Role != entity.UserRoleOwner {
		t.Fatalf("expected role to remain owner, got %s", result.User.Role)
	}
}

func TestGoogleLogin_InvalidToken(t *testing.T) {
	uc, _, _, _ := setupAuthUsecase(nil, errors.New("invalid token"))
	ctx := context.Background()

	input := &usecase.GoogleLoginInput{IDToken: "invalid-token"}
	_, err := uc.GoogleLogin(ctx, input)
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	if !errors.Is(err, usecase.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestGoogleLogin_UnverifiedEmail(t *testing.T) {
	claims := &auth.GoogleClaims{
		Sub:           "google-123",
		Email:         "unverified@example.com",
		EmailVerified: false,
		Name:          "Unverified",
	}
	uc, _, _, _ := setupAuthUsecase(claims, nil)
	ctx := context.Background()

	input := &usecase.GoogleLoginInput{IDToken: "valid-token"}
	_, err := uc.GoogleLogin(ctx, input)
	if err == nil {
		t.Fatal("expected error for unverified email")
	}
	if !errors.Is(err, usecase.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestRefreshToken_Success(t *testing.T) {
	uc, userRepo, sessionRepo, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	user := createTestUser(userRepo, "user@example.com", "password123", entity.UserRoleRenter)

	// First login to get a real refresh token
	loginResult, err := uc.Login(ctx, &usecase.LoginInput{
		Email:    "user@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	// Ensure session exists
	if len(sessionRepo.sessions) == 0 {
		t.Fatal("expected session to be created")
	}

	result, err := uc.RefreshToken(ctx, loginResult.RefreshToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.AccessToken == "" {
		t.Fatal("expected new access token")
	}
	if result.RefreshToken == loginResult.RefreshToken {
		t.Fatal("expected new refresh token (rotation)")
	}
	_ = user
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	uc, _, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	_, err := uc.RefreshToken(ctx, "invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	if !errors.Is(err, usecase.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestForgotPassword_ExistingUser(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	createTestUser(userRepo, "user@example.com", "password123", entity.UserRoleRenter)

	result, err := uc.ForgotPassword(ctx, "user@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ResetToken == "" {
		t.Fatal("expected reset token in dev mode")
	}
	if result.Message == "" {
		t.Fatal("expected message")
	}
}

func TestForgotPassword_NonexistentUser(t *testing.T) {
	uc, _, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	result, err := uc.ForgotPassword(ctx, "nobody@example.com")
	if err != nil {
		t.Fatalf("expected no error (info leak prevention), got %v", err)
	}
	if result.Message == "" {
		t.Fatal("expected message")
	}
}

func TestResetPassword_Success(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	createTestUser(userRepo, "user@example.com", "old-password", entity.UserRoleRenter)

	// Request reset
	forgotResult, err := uc.ForgotPassword(ctx, "user@example.com")
	if err != nil {
		t.Fatalf("forgot password failed: %v", err)
	}

	// Reset with new password
	err = uc.ResetPassword(ctx, forgotResult.ResetToken, "new-password-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Login with new password should work
	_, err = uc.Login(ctx, &usecase.LoginInput{
		Email:    "user@example.com",
		Password: "new-password-123",
	})
	if err != nil {
		t.Fatalf("expected login with new password to succeed, got %v", err)
	}

	// Login with old password should fail
	_, err = uc.Login(ctx, &usecase.LoginInput{
		Email:    "user@example.com",
		Password: "old-password",
	})
	if err == nil {
		t.Fatal("expected login with old password to fail")
	}
}

func TestResetPassword_InvalidToken(t *testing.T) {
	uc, _, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	err := uc.ResetPassword(ctx, "invalid-token", "new-password")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	if !errors.Is(err, usecase.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestResetPassword_TokenUsedTwice(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	createTestUser(userRepo, "user@example.com", "password123", entity.UserRoleRenter)

	// Request reset
	forgotResult, err := uc.ForgotPassword(ctx, "user@example.com")
	if err != nil {
		t.Fatalf("forgot password failed: %v", err)
	}

	// First reset should succeed
	err = uc.ResetPassword(ctx, forgotResult.ResetToken, "new-password-1")
	if err != nil {
		t.Fatalf("first reset failed: %v", err)
	}

	// Second reset with same token should fail
	err = uc.ResetPassword(ctx, forgotResult.ResetToken, "new-password-2")
	if err == nil {
		t.Fatal("expected error for reused token")
	}
	if !errors.Is(err, usecase.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestLogout_Success(t *testing.T) {
	uc, userRepo, sessionRepo, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	createTestUser(userRepo, "user@example.com", "password123", entity.UserRoleRenter)

	// Login first
	loginResult, err := uc.Login(ctx, &usecase.LoginInput{
		Email:    "user@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	_ = loginResult

	// Get session ID
	var sessionID string
	for id := range sessionRepo.sessions {
		sessionID = id
		break
	}

	err = uc.Logout(ctx, sessionID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Session should be revoked
	session := sessionRepo.sessions[sessionID]
	if session.RevokedAt == nil {
		t.Fatal("expected session to be revoked")
	}
}

func TestRegister_OwnerRole(t *testing.T) {
	uc, _, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	input := &usecase.RegisterInput{
		Email:    "owner@example.com",
		Password: "password123",
		FullName: "Store Owner",
		Role:     entity.UserRoleOwner,
	}

	result, err := uc.Register(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.User.Role != entity.UserRoleOwner {
		t.Fatalf("expected role owner, got %s", result.User.Role)
	}
}

func TestLogin_GoogleOnlyUser(t *testing.T) {
	uc, userRepo, _, _ := setupAuthUsecase(nil, nil)
	ctx := context.Background()

	// Create user without password (Google only)
	googleID := "google-only-123"
	user := &entity.User{
		BaseModel: entity.BaseModel{ID: "user-google-only"},
		Email:     "googleonly@example.com",
		FullName:  "Google Only",
		GoogleID:  &googleID,
		Provider:  "google",
		Role:      entity.UserRoleRenter,
		IsActive:  true,
	}
	userRepo.users[user.ID] = user
	userRepo.byEmail[user.Email] = user
	userRepo.byGoogle[googleID] = user

	input := &usecase.LoginInput{
		Email:    "googleonly@example.com",
		Password: "any-password",
	}

	_, err := uc.Login(ctx, input)
	if err == nil {
		t.Fatal("expected error for Google-only user trying password login")
	}
	if !errors.Is(err, usecase.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}
