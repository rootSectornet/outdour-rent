package usecase

import (
	"context"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/pkg/pagination"
)

// ReviewUsecase defines the interface for review business logic.
type ReviewUsecase interface {
	Create(ctx context.Context, input *CreateReviewInput) (*entity.Review, error)
	Reply(ctx context.Context, reviewID string, reply string, ownerID string) error
	ListByEquipment(ctx context.Context, equipmentID string, params *pagination.Params) ([]entity.Review, *pagination.Meta, error)
	ListByStore(ctx context.Context, storeID string, params *pagination.Params) ([]entity.Review, *pagination.Meta, error)
}

type CreateReviewInput struct {
	UserID      string
	EquipmentID string
	OrderID     string
	StoreID     string
	Rating      uint8
	Comment     string
}

// StoreUsecase defines the interface for store business logic.
// StoreUsecase defines the interface for store business logic.
type StoreUsecase interface {
	Create(ctx context.Context, input *CreateStoreInput) (*entity.Store, error)
	GetByID(ctx context.Context, id string) (*entity.Store, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Store, error)
	GetMyStore(ctx context.Context, ownerID string) (*entity.Store, error)
	Update(ctx context.Context, id string, input *UpdateStoreInput, ownerID string) (*entity.Store, error)
	List(ctx context.Context, input *ListStoreInput) ([]entity.Store, *pagination.Meta, error)

	// Photos
	AddPhoto(ctx context.Context, input *AddStorePhotoInput) (*entity.StorePhoto, error)
	RemovePhoto(ctx context.Context, storeID string, photoID string, ownerID string) error
	SetPrimaryPhoto(ctx context.Context, storeID string, photoID string, ownerID string) error

	// Operating Hours
	SetOperatingHours(ctx context.Context, storeID string, ownerID string, hours []OperatingHourInput) error
	GetOperatingHours(ctx context.Context, storeID string) ([]entity.StoreOperatingHour, error)

	// Admin actions
	ApproveStore(ctx context.Context, storeID string, adminID string) error
	SuspendStore(ctx context.Context, storeID string, reason string, adminID string) error
	ReactivateStore(ctx context.Context, storeID string, adminID string) error
}

type CreateStoreInput struct {
	OwnerID     string
	Name        string
	Description string
	Phone       string
	Email       string
	Address     string
	City        string
	Province    string
	PostalCode  string
	Latitude    float64
	Longitude   float64
}

type UpdateStoreInput struct {
	Name        *string
	Description *string
	Phone       *string
	Email       *string
	Address     *string
	City        *string
	Province    *string
	PostalCode  *string
	Latitude    *float64
	Longitude   *float64
	LogoURL     *string
	BannerURL   *string
}

type ListStoreInput struct {
	City     string
	Province string
	Status   string
	Search   string
	Page     int
	PerPage  int
}

type AddStorePhotoInput struct {
	StoreID   string
	OwnerID   string
	PhotoURL  string
	Caption   string
	SortOrder int
	IsPrimary bool
}

type OperatingHourInput struct {
	DayOfWeek uint8  `json:"day_of_week" binding:"required,min=1,max=7"`
	OpenTime  string `json:"open_time" binding:"required"`
	CloseTime string `json:"close_time" binding:"required"`
	IsClosed  bool   `json:"is_closed"`
}

// UserUsecase defines the interface for user profile business logic.
type UserUsecase interface {
	GetProfile(ctx context.Context, userID string) (*entity.User, error)
	UpdateProfile(ctx context.Context, userID string, input *UpdateProfileInput) (*entity.User, error)
}

type UpdateProfileInput struct {
	FullName  *string
	Phone     *string
	AvatarURL *string
}

// NotificationUsecase defines the interface for notification business logic.
type NotificationUsecase interface {
	List(ctx context.Context, userID string, params *pagination.Params) ([]entity.Notification, *pagination.Meta, error)
	MarkAsRead(ctx context.Context, id string, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	CountUnread(ctx context.Context, userID string) (int64, error)
}
