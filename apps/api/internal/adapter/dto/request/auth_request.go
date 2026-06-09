package request

// AuthRequests

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=8,max=72" example:"securePass123"`
	FullName string `json:"full_name" binding:"required,min=2,max=100" example:"John Doe"`
	Phone    string `json:"phone" binding:"omitempty,min=10,max=20" example:"08123456789"`
	Role     string `json:"role" binding:"omitempty,oneof=renter owner" example:"renter"`
}

type RegisterInvitationRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"securePass123"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required" example:"eyJhbGciOiJSUzI1NiIs..."`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"a1b2c3d4e5f6..."`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" example:"a1b2c3d4e5f6..."`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72" example:"newSecurePass123"`
}
