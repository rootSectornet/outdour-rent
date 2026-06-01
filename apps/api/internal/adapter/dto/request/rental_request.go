package request

// RentalRequests

type CreateOrderRequest struct {
	StoreID         string             `json:"store_id" binding:"required,uuid"`
	Items           []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
	RentalStartDate string             `json:"rental_start_date" binding:"required" example:"2026-07-01"`
	RentalEndDate   string             `json:"rental_end_date" binding:"required" example:"2026-07-03"`
	Notes           string             `json:"notes" binding:"omitempty,max=1000"`
	PickupMethod    string             `json:"pickup_method" binding:"required,oneof=pickup delivery"`
	DeliveryAddress string             `json:"delivery_address" binding:"required_if=PickupMethod delivery,max=500"`
}

type OrderItemRequest struct {
	EquipmentID string `json:"equipment_id" binding:"required,uuid"`
	Quantity    uint   `json:"quantity" binding:"required,min=1"`
}

type ApproveOrderRequest struct{}

type RejectOrderRequest struct {
	Reason string `json:"reason" binding:"required,min=5,max=500"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=500"`
}

// PaymentRequests

type InitiatePaymentRequest struct {
	OrderID string `json:"order_id" binding:"required,uuid"`
}

// ReviewRequests

type CreateReviewRequest struct {
	EquipmentID string `json:"equipment_id" binding:"required,uuid"`
	OrderID     string `json:"order_id" binding:"required,uuid"`
	Rating      uint8  `json:"rating" binding:"required,min=1,max=5"`
	Comment     string `json:"comment" binding:"omitempty,max=2000"`
}

type ReplyReviewRequest struct {
	Reply string `json:"reply" binding:"required,min=1,max=2000"`
}

// StoreRequests

type CreateStoreRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description string  `json:"description" binding:"omitempty,max=5000"`
	Phone       string  `json:"phone" binding:"required,min=10,max=20"`
	Email       string  `json:"email" binding:"required,email,max=255"`
	Address     string  `json:"address" binding:"required,min=10,max=500"`
	City        string  `json:"city" binding:"required,max=100"`
	Province    string  `json:"province" binding:"required,max=100"`
	PostalCode  string  `json:"postal_code" binding:"omitempty,max=10"`
	Latitude    float64 `json:"latitude" binding:"omitempty"`
	Longitude   float64 `json:"longitude" binding:"omitempty"`
}

type UpdateStoreRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty,max=5000"`
	Phone       *string `json:"phone" binding:"omitempty,min=10,max=20"`
	Address     *string `json:"address" binding:"omitempty,min=10,max=500"`
	City        *string `json:"city" binding:"omitempty,max=100"`
	Province    *string `json:"province" binding:"omitempty,max=100"`
}

// UserRequests

type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"omitempty,min=2,max=100"`
	Phone    string `json:"phone" binding:"omitempty,min=10,max=20"`
}
