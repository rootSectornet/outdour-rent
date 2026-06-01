package usecase

import (
	"context"

	"github.com/rentoutdoor/api/internal/domain/entity"
)

// PaymentUsecase defines the interface for payment business logic.
type PaymentUsecase interface {
	InitiatePayment(ctx context.Context, input *InitiatePaymentInput) (*PaymentOutput, error)
	HandleCallback(ctx context.Context, input *PaymentCallbackInput) error
	CheckStatus(ctx context.Context, paymentID string) (*entity.Payment, error)
}

type InitiatePaymentInput struct {
	OrderID   string
	UserID    string
}

type PaymentOutput struct {
	PaymentID   string
	SnapToken   string
	RedirectURL string
}

type PaymentCallbackInput struct {
	OrderID           string
	TransactionID     string
	TransactionStatus string
	PaymentType       string
	FraudStatus       string
	RawPayload        string
}
