package entity

import "time"

type PaymentType string

const (
	PaymentTypeRental  PaymentType = "rental"
	PaymentTypeDeposit PaymentType = "deposit"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusSuccess  PaymentStatus = "success"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusExpired  PaymentStatus = "expired"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Payment struct {
	BaseModelNoSoftDelete
	OrderID              string        `gorm:"type:char(36);not null;index:idx_payments_order_type,priority:1" json:"order_id"`
	PaymentType          PaymentType   `gorm:"type:enum('rental','deposit');not null;index:idx_payments_order_type,priority:2" json:"payment_type"`
	Amount               float64       `gorm:"type:decimal(12,2);not null" json:"amount"`
	Method               *string       `gorm:"type:varchar(50)" json:"method,omitempty"`
	Status               PaymentStatus `gorm:"type:enum('pending','success','failed','expired','refunded');not null;index:idx_payments_status" json:"status"`
	MidtransTransactionID *string      `gorm:"type:varchar(100)" json:"midtrans_transaction_id,omitempty"`
	MidtransOrderID      *string       `gorm:"type:varchar(100);index:idx_payments_midtrans_order" json:"midtrans_order_id,omitempty"`
	MidtransSnapToken    *string       `gorm:"type:varchar(255)" json:"midtrans_snap_token,omitempty"`
	MidtransRedirectURL  *string       `gorm:"type:varchar(500)" json:"midtrans_redirect_url,omitempty"`
	MidtransPaymentType  *string       `gorm:"type:varchar(50)" json:"midtrans_payment_type,omitempty"`
	PaidAt               *time.Time    `json:"paid_at,omitempty"`
	ExpiredAt            *time.Time    `json:"expired_at,omitempty"`
	RawCallback          *string       `gorm:"type:json" json:"raw_callback,omitempty"`

	// Relations
	Order   Order    `gorm:"foreignKey:OrderID;constraint:OnDelete:RESTRICT" json:"order,omitempty"`
	Refunds []Refund `gorm:"foreignKey:PaymentID" json:"refunds,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

type DepositStatus string

const (
	DepositStatusPending           DepositStatus = "pending"
	DepositStatusHeld              DepositStatus = "held"
	DepositStatusReturned          DepositStatus = "returned"
	DepositStatusForfeited         DepositStatus = "forfeited"
	DepositStatusPartiallyReturned DepositStatus = "partially_returned"
)

type Deposit struct {
	BaseModelNoSoftDelete
	OrderID         string        `gorm:"type:char(36);not null;index:idx_deposits_order_id" json:"order_id"`
	PaymentID       *string       `gorm:"type:char(36)" json:"payment_id,omitempty"`
	Amount          float64       `gorm:"type:decimal(12,2);not null" json:"amount"`
	Status          DepositStatus `gorm:"type:enum('pending','held','returned','forfeited','partially_returned');not null;index:idx_deposits_status" json:"status"`
	ReturnedAmount  float64       `gorm:"type:decimal(12,2);not null;default:0.00" json:"returned_amount"`
	ForfeitedAmount float64       `gorm:"type:decimal(12,2);not null;default:0.00" json:"forfeited_amount"`
	ForfeitReason   *string       `gorm:"type:varchar(500)" json:"forfeit_reason,omitempty"`
	ReturnedAt      *time.Time    `json:"returned_at,omitempty"`

	// Relations
	Order   Order    `gorm:"foreignKey:OrderID;constraint:OnDelete:RESTRICT" json:"order,omitempty"`
	Payment *Payment `gorm:"foreignKey:PaymentID;constraint:OnDelete:SET NULL" json:"payment,omitempty"`
}

func (Deposit) TableName() string {
	return "deposits"
}

type RefundStatus string

const (
	RefundStatusPending    RefundStatus = "pending"
	RefundStatusProcessing RefundStatus = "processing"
	RefundStatusSuccess    RefundStatus = "success"
	RefundStatusFailed     RefundStatus = "failed"
)

type Refund struct {
	BaseModelNoSoftDelete
	PaymentID        string       `gorm:"type:char(36);not null;index:idx_refunds_payment_id" json:"payment_id"`
	OrderID          string       `gorm:"type:char(36);not null;index:idx_refunds_order_id" json:"order_id"`
	Amount           float64      `gorm:"type:decimal(12,2);not null" json:"amount"`
	Reason           string       `gorm:"type:varchar(500);not null" json:"reason"`
	Status           RefundStatus `gorm:"type:enum('pending','processing','success','failed');not null" json:"status"`
	MidtransRefundID *string      `gorm:"type:varchar(100)" json:"midtrans_refund_id,omitempty"`
	ProcessedAt      *time.Time   `json:"processed_at,omitempty"`

	// Relations
	Payment Payment `gorm:"foreignKey:PaymentID;constraint:OnDelete:RESTRICT" json:"payment,omitempty"`
	Order   Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:RESTRICT" json:"order,omitempty"`
}

func (Refund) TableName() string {
	return "refunds"
}
