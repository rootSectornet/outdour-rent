package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rentoutdoor/api/internal/adapter/dto/request"
	"github.com/rentoutdoor/api/internal/adapter/middleware"
	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/usecase"
	"github.com/rentoutdoor/api/pkg/pagination"
	"github.com/rentoutdoor/api/pkg/response"
)

// EquipmentHandler handles equipment endpoints.
type EquipmentHandler struct {
	equipmentUC usecase.EquipmentUsecase
}

// NewEquipmentHandler creates a new equipment handler.
func NewEquipmentHandler(equipmentUC usecase.EquipmentUsecase) *EquipmentHandler {
	return &EquipmentHandler{equipmentUC: equipmentUC}
}

// List godoc
// @Summary List equipment
// @Description List equipment with optional filters by store, category, status, and search.
// @Tags Equipment
// @Produce json
// @Param store_id query string false "Store ID"
// @Param category_id query string false "Category ID"
// @Param status query string false "Filter by status" Enums(available, reserved, rented, maintenance, damaged, retired)
// @Param search query string false "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} response.Response
// @Router /equipment [get]
func (h *EquipmentHandler) List(c *gin.Context) {
	var req request.ListEquipmentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	input := &usecase.ListEquipmentInput{
		StoreID:    req.StoreID,
		CategoryID: req.CategoryID,
		City:       req.City,
		Search:     req.Search,
		Status:     req.Status,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Params:     pagination.NewParams(req.Page, req.PerPage),
	}

	items, meta, err := h.equipmentUC.List(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "equipment retrieved", items, meta)
}

// GetByID godoc
// @Summary Get equipment by ID
// @Description Get detailed equipment information including photos, pricing, and category.
// @Tags Equipment
// @Produce json
// @Param id path string true "Equipment ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.ErrorResponse
// @Router /equipment/{id} [get]
func (h *EquipmentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	equipment, err := h.equipmentUC.GetByID(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "equipment retrieved", equipment)
}

// Create godoc
// @Summary Create new equipment
// @Description Create a new equipment item for a store. Only store owners can create.
// @Tags Equipment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.CreateEquipmentRequest true "Equipment payload"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Router /equipment [post]
func (h *EquipmentHandler) Create(c *gin.Context) {
	var req request.CreateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	userID := middleware.GetUserID(c)

	input := &usecase.CreateEquipmentInput{
		CategoryID:      req.CategoryID,
		Name:            req.Name,
		Description:     req.Description,
		Brand:           req.Brand,
		Specifications:  req.Specifications,
		TotalStock:      req.TotalStock,
		Condition:       entity.EquipmentCondition(req.Condition),
		WeightGrams:     req.WeightGrams,
		MinRentalDays:   req.MinRentalDays,
		MaxRentalDays:   req.MaxRentalDays,
		RequiresDeposit: req.RequiresDeposit,
		DepositAmount:   req.DepositAmount,
		CreatedBy:       userID,
	}

	if req.PurchaseDate != "" {
		t, err := time.Parse("2006-01-02", req.PurchaseDate)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid purchase_date format, use YYYY-MM-DD", err.Error())
			return
		}
		input.PurchaseDate = &t
	}

	equipment, err := h.equipmentUC.Create(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "equipment created", equipment)
}

// Update godoc
// @Summary Update equipment
// @Description Update equipment details. Only the store owner can update.
// @Tags Equipment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Equipment ID"
// @Param body body request.UpdateEquipmentRequest true "Update payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /equipment/{id} [put]
func (h *EquipmentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	userID := middleware.GetUserID(c)

	input := &usecase.UpdateEquipmentInput{
		Name:          req.Name,
		Description:   req.Description,
		Brand:         req.Brand,
		TotalStock:    req.TotalStock,
		MinRentalDays: req.MinRentalDays,
		MaxRentalDays: req.MaxRentalDays,
		DepositAmount: req.DepositAmount,
		IsActive:      req.IsActive,
	}
	if req.Condition != nil {
		cond := entity.EquipmentCondition(*req.Condition)
		input.Condition = &cond
	}
	if req.PurchaseDate != nil {
		t, err := time.Parse("2006-01-02", *req.PurchaseDate)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid purchase_date format, use YYYY-MM-DD", err.Error())
			return
		}
		input.PurchaseDate = &t
	}

	equipment, err := h.equipmentUC.Update(c.Request.Context(), id, input, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "equipment updated", equipment)
}

// Delete godoc
// @Summary Delete equipment (soft)
// @Description Soft delete equipment. Only the store owner can delete.
// @Tags Equipment
// @Security BearerAuth
// @Produce json
// @Param id path string true "Equipment ID"
// @Success 200 {object} response.Response
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /equipment/{id} [delete]
func (h *EquipmentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.equipmentUC.Delete(c.Request.Context(), id, userID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "equipment deleted", nil)
}

// ChangeStatus godoc
// @Summary Change equipment status
// @Description Change equipment status (available, reserved, rented, maintenance, damaged, retired). Validates allowed transitions.
// @Tags Equipment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Equipment ID"
// @Param body body request.ChangeEquipmentStatusRequest true "New status"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /equipment/{id}/status [patch]
func (h *EquipmentHandler) ChangeStatus(c *gin.Context) {
	id := c.Param("id")
	var req request.ChangeEquipmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	userID := middleware.GetUserID(c)
	status := entity.EquipmentStatus(req.Status)

	equipment, err := h.equipmentUC.ChangeStatus(c.Request.Context(), id, status, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "equipment status updated", equipment)
}

// CheckAvailability godoc
// @Summary Check equipment availability for date range
// @Description Check if equipment is available for rental in the given date range and quantity.
// @Tags Equipment
// @Accept json
// @Produce json
// @Param id path string true "Equipment ID"
// @Param body body request.CheckAvailabilityRequest true "Availability check"
// @Success 200 {object} response.Response{data=usecase.AvailabilityOutput}
// @Failure 404 {object} response.ErrorResponse
// @Router /equipment/{id}/availability [post]
func (h *EquipmentHandler) CheckAvailability(c *gin.Context) {
	id := c.Param("id")
	var req request.CheckAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	input := &usecase.AvailabilityInput{
		EquipmentID: id,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Quantity:    req.Quantity,
	}

	result, err := h.equipmentUC.CheckAvailability(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "availability checked", result)
}
