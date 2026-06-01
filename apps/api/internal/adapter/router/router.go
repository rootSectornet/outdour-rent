package router

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/rentoutdoor/api/docs"
	"github.com/rentoutdoor/api/internal/adapter/middleware"
	"github.com/rentoutdoor/api/internal/infrastructure/container"
)

// Setup configures all routes for the application.
func Setup(r *gin.Engine, c *container.Container) {
	// Global middleware
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.CORSMiddleware([]string{"*"})) // TODO: Restrict in production
	r.Use(middleware.LoggerMiddleware(c.Logger))
	r.Use(middleware.RecoveryMiddleware(c.Logger))

	// Rate limiter: 100 requests per minute per IP
	rl := middleware.NewRateLimiter(100, time.Minute)
	r.Use(middleware.RateLimitMiddleware(rl))

	// Health check (no auth)
	r.GET("/health", c.HealthHandler.Health)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	v1 := r.Group("/api/v1")
	{
		setupAuthRoutes(v1, c)
		setupEquipmentRoutes(v1, c)
		setupRentalRoutes(v1, c)
		setupPaymentRoutes(v1, c)
	}
}

func setupAuthRoutes(rg *gin.RouterGroup, c *container.Container) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", c.AuthHandler.Register)
		auth.POST("/login", c.AuthHandler.Login)
		auth.POST("/google", c.AuthHandler.GoogleLogin)
		auth.POST("/refresh", c.AuthHandler.RefreshToken)
		auth.POST("/forgot-password", c.AuthHandler.ForgotPassword)
		auth.POST("/reset-password", c.AuthHandler.ResetPassword)
		auth.POST("/logout", middleware.AuthMiddleware(c.Config.JWT.AccessSecret), c.AuthHandler.Logout)
	}
}

func setupEquipmentRoutes(rg *gin.RouterGroup, c *container.Container) {
	equipment := rg.Group("/equipment")
	{
		// Public routes
		equipment.GET("", c.EquipmentHandler.List)
		equipment.GET("/:id", c.EquipmentHandler.GetByID)
		equipment.POST("/:id/availability", c.EquipmentHandler.CheckAvailability)

		// Authenticated routes (owner only)
		ownerEquipment := equipment.Group("")
		ownerEquipment.Use(middleware.AuthMiddleware(c.Config.JWT.AccessSecret))
		ownerEquipment.Use(middleware.RoleMiddleware("owner", "admin"))
		{
			ownerEquipment.POST("", c.EquipmentHandler.Create)
			ownerEquipment.PUT("/:id", c.EquipmentHandler.Update)
			ownerEquipment.DELETE("/:id", c.EquipmentHandler.Delete)
		}
	}
}

func setupRentalRoutes(rg *gin.RouterGroup, c *container.Container) {
	rentals := rg.Group("/rentals")
	rentals.Use(middleware.AuthMiddleware(c.Config.JWT.AccessSecret))
	{
		// Renter routes
		rentals.POST("", c.RentalHandler.CreateOrder)
		rentals.GET("", c.RentalHandler.ListMyOrders)
		rentals.GET("/:id", c.RentalHandler.GetOrder)
		rentals.PATCH("/:id/cancel", c.RentalHandler.CancelOrder)

		// Owner routes
		owner := rentals.Group("")
		owner.Use(middleware.RoleMiddleware("owner", "admin"))
		{
			owner.GET("/incoming", c.RentalHandler.ListStoreOrders)
			owner.PATCH("/:id/approve", c.RentalHandler.ApproveOrder)
			owner.PATCH("/:id/reject", c.RentalHandler.RejectOrder)
			owner.PATCH("/:id/complete", c.RentalHandler.CompleteOrder)
		}
	}
}

func setupPaymentRoutes(rg *gin.RouterGroup, c *container.Container) {
	payments := rg.Group("/payments")
	{
		// Webhook (no auth - verified by Midtrans signature)
		payments.POST("/callback", c.PaymentHandler.HandleCallback)

		// Authenticated routes
		authenticated := payments.Group("")
		authenticated.Use(middleware.AuthMiddleware(c.Config.JWT.AccessSecret))
		{
			authenticated.POST("", c.PaymentHandler.InitiatePayment)
			authenticated.GET("/:id/status", c.PaymentHandler.CheckStatus)
		}
	}
}
