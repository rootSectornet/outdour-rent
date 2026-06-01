package container

import (
	"github.com/rentoutdoor/api/internal/adapter/handler"
	"github.com/rentoutdoor/api/internal/infrastructure/config"
	"github.com/rentoutdoor/api/internal/infrastructure/google"
	"github.com/rentoutdoor/api/internal/infrastructure/persistence/mysql"
	"github.com/rentoutdoor/api/internal/repository"
	"github.com/rentoutdoor/api/internal/usecase"
	"github.com/rentoutdoor/api/internal/usecase/auth"
	"github.com/rentoutdoor/api/internal/usecase/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Container holds all application dependencies.
type Container struct {
	Config *config.Config
	DB     *gorm.DB
	Logger *zap.Logger

	// Repositories
	UserRepo          repository.UserRepository
	SessionRepo       repository.SessionRepository
	PasswordResetRepo repository.PasswordResetRepository
	StoreRepo         repository.StoreRepository
	StorePhotoRepo    repository.StorePhotoRepository
	StoreHoursRepo    repository.StoreOperatingHourRepository
	EquipmentRepo     repository.EquipmentRepository
	CategoryRepo      repository.CategoryRepository
	ReservationRepo   repository.ReservationRepository
	OrderRepo        repository.OrderRepository
	OrderItemRepo    repository.OrderItemRepository
	PaymentRepo      repository.PaymentRepository
	DepositRepo      repository.DepositRepository
	RefundRepo       repository.RefundRepository
	ReviewRepo       repository.ReviewRepository
	NotificationRepo repository.NotificationRepository
	MaintenanceRepo  repository.MaintenanceRepository

	// Transaction Manager
	TxManager repository.TransactionManager

	// Usecases
	AuthUC         usecase.AuthUsecase
	EquipmentUC    usecase.EquipmentUsecase
	RentalUC       usecase.RentalUsecase
	PaymentUC      usecase.PaymentUsecase
	ReviewUC       usecase.ReviewUsecase
	StoreUC        usecase.StoreUsecase
	UserUC         usecase.UserUsecase
	NotificationUC usecase.NotificationUsecase

	// Handlers
	HealthHandler    *handler.HealthHandler
	AuthHandler      *handler.AuthHandler
	StoreHandler     *handler.StoreHandler
	EquipmentHandler *handler.EquipmentHandler
	RentalHandler    *handler.RentalHandler
	PaymentHandler   *handler.PaymentHandler
}

// New creates a new dependency injection container.
func New(cfg *config.Config, db *gorm.DB, logger *zap.Logger) *Container {
	c := &Container{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}

	c.initRepositories()
	c.initUsecases()
	c.initHandlers()

	return c
}

func (c *Container) initRepositories() {
	c.UserRepo = mysql.NewUserRepository(c.DB)
	c.SessionRepo = mysql.NewSessionRepository(c.DB)
	c.PasswordResetRepo = mysql.NewPasswordResetRepository(c.DB)
	c.StoreRepo = mysql.NewStoreRepository(c.DB)
	c.StorePhotoRepo = mysql.NewStorePhotoRepository(c.DB)
	c.StoreHoursRepo = mysql.NewStoreOperatingHourRepository(c.DB)
	c.EquipmentRepo = mysql.NewEquipmentRepository(c.DB)
	c.CategoryRepo = mysql.NewCategoryRepository(c.DB)
	c.ReservationRepo = mysql.NewReservationRepository(c.DB)
	c.TxManager = mysql.NewTxManager(c.DB)

	// TODO: Initialize remaining repositories when implementations are created
	// c.OrderRepo = mysql.NewOrderRepository(c.DB)
	// c.OrderItemRepo = mysql.NewOrderItemRepository(c.DB)
	// c.PaymentRepo = mysql.NewPaymentRepository(c.DB)
	// c.DepositRepo = mysql.NewDepositRepository(c.DB)
	// c.RefundRepo = mysql.NewRefundRepository(c.DB)
	// c.ReviewRepo = mysql.NewReviewRepository(c.DB)
	// c.NotificationRepo = mysql.NewNotificationRepository(c.DB)
	// c.MaintenanceRepo = mysql.NewMaintenanceRepository(c.DB)
}

func (c *Container) initUsecases() {
	googleVerifier := google.NewGoogleVerifier(c.Config.Google.ClientID)

	c.AuthUC = auth.NewAuthUsecase(
		c.UserRepo,
		c.SessionRepo,
		c.PasswordResetRepo,
		c.Config.JWT,
		googleVerifier,
	)

	c.StoreUC = store.NewStoreUsecase(c.StoreRepo, c.StorePhotoRepo, c.StoreHoursRepo, c.UserRepo)

	// TODO: Initialize remaining usecases when implementations are created
	// c.EquipmentUC = equipment.NewEquipmentUsecase(c.EquipmentRepo, c.CategoryRepo, c.ReservationRepo, c.TxManager)
	// c.RentalUC = rental.NewRentalUsecase(c.OrderRepo, c.OrderItemRepo, c.ReservationRepo, c.EquipmentRepo, c.TxManager)
	// c.PaymentUC = payment.NewPaymentUsecase(c.PaymentRepo, c.OrderRepo, c.ReservationRepo, c.Config.Midtrans)
	// c.ReviewUC = review.NewReviewUsecase(c.ReviewRepo, c.OrderRepo, c.EquipmentRepo)
	// c.UserUC = user.NewUserUsecase(c.UserRepo)
	// c.NotificationUC = notification.NewNotificationUsecase(c.NotificationRepo)
}

func (c *Container) initHandlers() {
	c.HealthHandler = handler.NewHealthHandler(c.DB)
	c.AuthHandler = handler.NewAuthHandler(c.AuthUC)
	c.StoreHandler = handler.NewStoreHandler(c.StoreUC)
	c.EquipmentHandler = handler.NewEquipmentHandler(c.EquipmentUC)
	c.RentalHandler = handler.NewRentalHandler(c.RentalUC)
	c.PaymentHandler = handler.NewPaymentHandler(c.PaymentUC)
}
