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
		setupStoreRoutes(v1, c)
		setupEquipmentRoutes(v1, c)
		setupRentalRoutes(v1, c)
		setupPaymentRoutes(v1, c)
		setupAdminRoutes(v1, c)
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

func setupStoreRoutes(rg *gin.RouterGroup, c *container.Container) {
	stores := rg.Group("/stores")
	{
		// Public routes
		stores.GET("", c.StoreHandler.List)
		stores.GET("/:id", c.StoreHandler.GetByID)
		stores.GET("/slug/:slug", c.StoreHandler.GetBySlug)
		stores.GET("/:id/operating-hours", c.StoreHandler.GetOperatingHours)

		// Authenticated owner routes
		ownerStores := stores.Group("")
		ownerStores.Use(middleware.AuthMiddleware(c.Config.JWT.AccessSecret))
		ownerStores.Use(middleware.RoleMiddleware("owner", "admin"))
		{
			ownerStores.POST("", c.StoreHandler.Create)
			ownerStores.GET("/me", c.StoreHandler.GetMyStore)
			ownerStores.PUT("/:id", c.StoreHandler.Update)

			// Photos
			ownerStores.POST("/:id/photos", c.StoreHandler.AddPhoto)
			ownerStores.DELETE("/:id/photos/:photoId", c.StoreHandler.RemovePhoto)
			ownerStores.PATCH("/:id/photos/:photoId/primary", c.StoreHandler.SetPrimaryPhoto)

			// Operating Hours
			ownerStores.PUT("/:id/operating-hours", c.StoreHandler.SetOperatingHours)
		}
	}
}

func setupAdminRoutes(rg *gin.RouterGroup, c *container.Container) {
	admin := rg.Group("/admin")
	admin.Use(middleware.AuthMiddleware(c.Config.JWT.AccessSecret))
	admin.Use(middleware.RoleMiddleware("admin"))
	{
		// Store management
		admin.PATCH("/stores/:id/approve", c.StoreHandler.ApproveStore)
		admin.PATCH("/stores/:id/suspend", c.StoreHandler.SuspendStore)
		admin.PATCH("/stores/:id/reactivate", c.StoreHandler.ReactivateStore)
	}
}
