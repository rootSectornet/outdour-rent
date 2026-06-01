package response

import "time"

// AuthResponse

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int64        `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	FullName        string     `json:"full_name"`
	Phone           *string    `json:"phone,omitempty"`
	AvatarURL       *string    `json:"avatar_url,omitempty"`
	Role            string     `json:"role"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// EquipmentResponse

type EquipmentResponse struct {
	ID              string                   `json:"id"`
	StoreID         string                   `json:"store_id"`
	CategoryID      string                   `json:"category_id"`
	Name            string                   `json:"name"`
	Slug            string                   `json:"slug"`
	Description     *string                  `json:"description,omitempty"`
	Brand           *string                  `json:"brand,omitempty"`
	Specifications  interface{}              `json:"specifications,omitempty"`
	TotalStock      uint                     `json:"total_stock"`
	Condition       string                   `json:"condition"`
	WeightGrams     *uint                    `json:"weight_grams,omitempty"`
	MinRentalDays   uint                     `json:"min_rental_days"`
	MaxRentalDays   uint                     `json:"max_rental_days"`
	RequiresDeposit bool                     `json:"requires_deposit"`
	DepositAmount   float64                  `json:"deposit_amount"`
	IsActive        bool                     `json:"is_active"`
	RatingAvg       float64                  `json:"rating_avg"`
	RatingCount     uint                     `json:"rating_count"`
	RentalCount     uint                     `json:"rental_count"`
	Photos          []EquipmentPhotoResponse `json:"photos,omitempty"`
	Pricing         []PricingResponse        `json:"pricing,omitempty"`
	Store           *StoreMinimalResponse    `json:"store,omitempty"`
	Category        *CategoryResponse        `json:"category,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
}

type EquipmentPhotoResponse struct {
	ID           string  `json:"id"`
	PhotoURL     string  `json:"photo_url"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
	SortOrder    int     `json:"sort_order"`
	IsPrimary    bool    `json:"is_primary"`
}

type PricingResponse struct {
	ID          string  `json:"id"`
	PricingType string  `json:"pricing_type"`
	MinDays     uint    `json:"min_days"`
	MaxDays     *uint   `json:"max_days,omitempty"`
	PricePerDay float64 `json:"price_per_day"`
}

type AvailabilityResponse struct {
	Available    bool `json:"available"`
	TotalStock   uint `json:"total_stock"`
	PeakUsage    uint `json:"peak_usage"`
	AvailableQty uint `json:"available_qty"`
}

// StoreResponse

type StoreResponse struct {
	ID          string                    `json:"id"`
	OwnerID     string                    `json:"owner_id"`
	Name        string                    `json:"name"`
	Slug        string                    `json:"slug"`
	Description *string                   `json:"description,omitempty"`
	Phone       string                    `json:"phone"`
	Email       string                    `json:"email"`
	Address     string                    `json:"address"`
	City        string                    `json:"city"`
	Province    string                    `json:"province"`
	PostalCode  *string                   `json:"postal_code,omitempty"`
	Latitude    *float64                  `json:"latitude,omitempty"`
	Longitude   *float64                  `json:"longitude,omitempty"`
	LogoURL     *string                   `json:"logo_url,omitempty"`
	BannerURL   *string                   `json:"banner_url,omitempty"`
	Status      string                    `json:"status"`
	VerifiedAt  *time.Time                `json:"verified_at,omitempty"`
	RatingAvg   float64                   `json:"rating_avg"`
	RatingCount uint                      `json:"rating_count"`
	Photos      []StorePhotoResponse      `json:"photos,omitempty"`
	Hours       []OperatingHourResponse   `json:"operating_hours,omitempty"`
	CreatedAt   time.Time                 `json:"created_at"`
}

type StorePhotoResponse struct {
	ID           string  `json:"id"`
	PhotoURL     string  `json:"photo_url"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
	Caption      *string `json:"caption,omitempty"`
	SortOrder    int     `json:"sort_order"`
	IsPrimary    bool    `json:"is_primary"`
}

type OperatingHourResponse struct {
	DayOfWeek uint8  `json:"day_of_week"`
	DayName   string `json:"day_name"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
	IsClosed  bool   `json:"is_closed"`
}

type StoreMinimalResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	City     string  `json:"city"`
	LogoURL  *string `json:"logo_url,omitempty"`
	Status   string  `json:"status"`
}

// CategoryResponse

type CategoryResponse struct {
	ID       string             `json:"id"`
	ParentID *string            `json:"parent_id,omitempty"`
	Name     string             `json:"name"`
	Slug     string             `json:"slug"`
	IconURL  *string            `json:"icon_url,omitempty"`
	Children []CategoryResponse `json:"children,omitempty"`
}

// OrderResponse

type OrderResponse struct {
	ID              string              `json:"id"`
	OrderNumber     string              `json:"order_number"`
	RenterID        string              `json:"renter_id"`
	StoreID         string              `json:"store_id"`
	Status          string              `json:"status"`
	RentalStartDate string              `json:"rental_start_date"`
	RentalEndDate   string              `json:"rental_end_date"`
	RentalDays      uint                `json:"rental_days"`
	Subtotal        float64             `json:"subtotal"`
	ServiceFee      float64             `json:"service_fee"`
	DepositTotal    float64             `json:"deposit_total"`
	TotalAmount     float64             `json:"total_amount"`
	Notes           *string             `json:"notes,omitempty"`
	PickupMethod    string              `json:"pickup_method"`
	PaymentDeadline time.Time           `json:"payment_deadline"`
	Items           []OrderItemResponse `json:"items,omitempty"`
	Store           *StoreMinimalResponse `json:"store,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
}

type OrderItemResponse struct {
	ID              string  `json:"id"`
	EquipmentID     string  `json:"equipment_id"`
	EquipmentName   string  `json:"equipment_name"`
	Quantity        uint    `json:"quantity"`
	PricePerDay     float64 `json:"price_per_day"`
	RentalDays      uint    `json:"rental_days"`
	Subtotal        float64 `json:"subtotal"`
	DepositPerUnit  float64 `json:"deposit_per_unit"`
	DepositSubtotal float64 `json:"deposit_subtotal"`
}

// PaymentResponse

type PaymentResponse struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	SnapToken   string  `json:"snap_token,omitempty"`
	RedirectURL string  `json:"redirect_url,omitempty"`
}

// ReviewResponse

type ReviewResponse struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	UserName       string     `json:"user_name"`
	EquipmentID    string     `json:"equipment_id"`
	Rating         uint8      `json:"rating"`
	Comment        *string    `json:"comment,omitempty"`
	OwnerReply     *string    `json:"owner_reply,omitempty"`
	OwnerRepliedAt *time.Time `json:"owner_replied_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
