package request

// EquipmentRequests

type CreateEquipmentRequest struct {
	CategoryID      string  `json:"category_id" binding:"required,uuid"`
	Name            string  `json:"name" binding:"required,min=3,max=200"`
	Description     string  `json:"description" binding:"omitempty,max=5000"`
	Brand           string  `json:"brand" binding:"omitempty,max=100"`
	Specifications  string  `json:"specifications" binding:"omitempty"`
	TotalStock      uint    `json:"total_stock" binding:"required,min=1"`
	Condition       string  `json:"condition" binding:"omitempty,oneof=new good fair worn"`
	WeightGrams     uint    `json:"weight_grams" binding:"omitempty"`
	MinRentalDays   uint    `json:"min_rental_days" binding:"omitempty,min=1"`
	MaxRentalDays   uint    `json:"max_rental_days" binding:"omitempty,min=1"`
	RequiresDeposit bool    `json:"requires_deposit"`
	DepositAmount   float64 `json:"deposit_amount" binding:"omitempty,min=0"`
}

type UpdateEquipmentRequest struct {
	Name          *string  `json:"name" binding:"omitempty,min=3,max=200"`
	Description   *string  `json:"description" binding:"omitempty,max=5000"`
	Brand         *string  `json:"brand" binding:"omitempty,max=100"`
	TotalStock    *uint    `json:"total_stock" binding:"omitempty,min=0"`
	Condition     *string  `json:"condition" binding:"omitempty,oneof=new good fair worn"`
	MinRentalDays *uint    `json:"min_rental_days" binding:"omitempty,min=1"`
	MaxRentalDays *uint    `json:"max_rental_days" binding:"omitempty,min=1"`
	DepositAmount *float64 `json:"deposit_amount" binding:"omitempty,min=0"`
	IsActive      *bool    `json:"is_active"`
}

type CheckAvailabilityRequest struct {
	StartDate string `json:"start_date" binding:"required" example:"2026-07-01"`
	EndDate   string `json:"end_date" binding:"required" example:"2026-07-03"`
	Quantity  uint   `json:"quantity" binding:"required,min=1"`
}

type ListEquipmentRequest struct {
	StoreID    string   `form:"store_id" binding:"omitempty,uuid"`
	CategoryID string   `form:"category_id" binding:"omitempty,uuid"`
	City       string   `form:"city" binding:"omitempty,max=100"`
	Search     string   `form:"search" binding:"omitempty,max=200"`
	MinPrice   *float64 `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice   *float64 `form:"max_price" binding:"omitempty,min=0"`
	Page       int      `form:"page" binding:"omitempty,min=1"`
	PerPage    int      `form:"per_page" binding:"omitempty,min=1,max=100"`
}
