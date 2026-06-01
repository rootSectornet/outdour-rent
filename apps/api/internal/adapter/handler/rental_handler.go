package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rentoutdoor/api/internal/adapter/dto/request"
	"github.com/rentoutdoor/api/internal/adapter/middleware"
	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/usecase"
	"github.com/rentoutdoor/api/pkg/pagination"
	"github.com/rentoutdoor/api/pkg/response"
)

// RentalHandler handles rental/order endpoints.
type RentalHandler struct {
	rentalUC usecase.RentalUsecase
}

// NewRentalHandler creates a new rental handler.
func NewRentalHandler(rentalUC usecase.RentalUsecase) *RentalHandler {
	return &RentalHandler{rentalUC: rentalUC}
}

// CreateOrder godoc
// @Summary Create a new rental order
// @Tags Rentals
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.CreateOrderRequest true "Order payload"
// @Success 201 {object} response.Response
// @Router /rentals [post]
func (h *RentalHandler) CreateOrder(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	userID := middleware.GetUserID(c)

	items := make([]usecase.OrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = usecase.OrderItemInput{
			EquipmentID: item.EquipmentID,
			Quantity:    item.Quantity,
		}
	}

	input := &usecase.CreateOrderInput{
		RenterID:        userID,
		StoreID:         req.StoreID,
		Items:           items,
		RentalStartDate: req.RentalStartDate,
		RentalEndDate:   req.RentalEndDate,
		Notes:           req.Notes,
		PickupMethod:    entity.PickupMethod(req.PickupMethod),
		DeliveryAddress: req.DeliveryAddress,
	}

	order, err := h.rentalUC.CreateOrder(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "order created", order)
}

// GetOrder godoc
// @Summary Get order details
// @Tags Rentals
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response
// @Router /rentals/{id} [get]
func (h *RentalHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := middleware.GetUserID(c)

	order, err := h.rentalUC.GetOrder(c.Request.Context(), orderID, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "order retrieved", order)
}

// ListMyOrders godoc
// @Summary List renter's orders
// @Tags Rentals
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} response.Response
// @Router /rentals [get]
func (h *RentalHandler) ListMyOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, perPage := pagination.FromQuery(c)

	params := &usecase.OrderListInput{
		Status: c.Query("status"),
		Params: pagination.NewParams(page, perPage),
	}

	orders, meta, err := h.rentalUC.ListRenterOrders(c.Request.Context(), userID, params)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "orders retrieved", orders, meta)
}

// ListStoreOrders godoc
// @Summary List store's incoming orders
// @Tags Rentals
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Success 200 {object} response.Response
// @Router /rentals/incoming [get]
func (h *RentalHandler) ListStoreOrders(c *gin.Context) {
	storeID := c.Param("store_id")
	page, perPage := pagination.FromQuery(c)

	params := &usecase.OrderListInput{
		Status: c.Query("status"),
		Params: pagination.NewParams(page, perPage),
	}

	orders, meta, err := h.rentalUC.ListStoreOrders(c.Request.Context(), storeID, params)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "orders retrieved", orders, meta)
}

// ApproveOrder godoc
// @Summary Approve an order (store owner)
// @Tags Rentals
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response
// @Router /rentals/{id}/approve [patch]
func (h *RentalHandler) ApproveOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.rentalUC.ApproveOrder(c.Request.Context(), orderID, userID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "order approved", nil)
}

// RejectOrder godoc
// @Summary Reject an order (store owner)
// @Tags Rentals
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param body body request.RejectOrderRequest true "Rejection reason"
// @Success 200 {object} response.Response
// @Router /rentals/{id}/reject [patch]
func (h *RentalHandler) RejectOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := middleware.GetUserID(c)

	var req request.RejectOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := h.rentalUC.RejectOrder(c.Request.Context(), orderID, req.Reason, userID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "order rejected", nil)
}

// CancelOrder godoc
// @Summary Cancel an order
// @Tags Rentals
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response
// @Router /rentals/{id}/cancel [patch]
func (h *RentalHandler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.rentalUC.CancelOrder(c.Request.Context(), orderID, userID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "order cancelled", nil)
}

// CompleteOrder godoc
// @Summary Mark order as completed (equipment returned)
// @Tags Rentals
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response
// @Router /rentals/{id}/complete [patch]
func (h *RentalHandler) CompleteOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.rentalUC.CompleteOrder(c.Request.Context(), orderID, userID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "order completed", nil)
}
