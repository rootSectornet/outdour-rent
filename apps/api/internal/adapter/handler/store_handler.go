package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rentoutdoor/api/internal/adapter/dto/request"
	"github.com/rentoutdoor/api/internal/adapter/middleware"
	"github.com/rentoutdoor/api/internal/usecase"
	"github.com/rentoutdoor/api/pkg/response"
)

// StoreHandler handles store management endpoints.
type StoreHandler struct {
	storeUC usecase.StoreUsecase
}

// NewStoreHandler creates a new store handler.
func NewStoreHandler(storeUC usecase.StoreUsecase) *StoreHandler {
	return &StoreHandler{storeUC: storeUC}
}

// Create godoc
// @Summary Create a new store
// @Description Creates a new store for the authenticated owner. One store per owner.
// @Tags Store
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body request.CreateStoreRequest true "Create store payload"
// @Success 201 {object} response.Response{data=entity.Store}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /stores [post]
func (h *StoreHandler) Create(c *gin.Context) {
	var req request.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ownerID := middleware.GetUserID(c)
	input := &usecase.CreateStoreInput{
		OwnerID:     ownerID,
		Name:        req.Name,
		Description: req.Description,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		City:        req.City,
		Province:    req.Province,
		PostalCode:  req.PostalCode,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	store, err := h.storeUC.Create(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "store created successfully", store)
}

// GetByID godoc
// @Summary Get store by ID
// @Description Get detailed store information including photos and operating hours.
// @Tags Store
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} response.Response{data=entity.Store}
// @Failure 404 {object} response.ErrorResponse
// @Router /stores/{id} [get]
func (h *StoreHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	store, err := h.storeUC.GetByID(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "store retrieved", store)
}

// GetBySlug godoc
// @Summary Get store by slug
// @Description Get store by its URL-friendly slug.
// @Tags Store
// @Produce json
// @Param slug path string true "Store slug"
// @Success 200 {object} response.Response{data=entity.Store}
// @Failure 404 {object} response.ErrorResponse
// @Router /stores/slug/{slug} [get]
func (h *StoreHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	store, err := h.storeUC.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "store retrieved", store)
}

// GetMyStore godoc
// @Summary Get my store
// @Description Get the authenticated owner's store.
// @Tags Store
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=entity.Store}
// @Failure 404 {object} response.ErrorResponse
// @Router /stores/me [get]
func (h *StoreHandler) GetMyStore(c *gin.Context) {
	ownerID := middleware.GetUserID(c)

	store, err := h.storeUC.GetMyStore(c.Request.Context(), ownerID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "store retrieved", store)
}

// Update godoc
// @Summary Update store
// @Description Update store information. Only the store owner can update.
// @Tags Store
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param body body request.UpdateStoreRequest true "Update store payload"
// @Success 200 {object} response.Response{data=entity.Store}
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /stores/{id} [put]
func (h *StoreHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ownerID := middleware.GetUserID(c)
	input := &usecase.UpdateStoreInput{
		Name:        req.Name,
		Description: req.Description,
		Phone:       req.Phone,
		Email:       req.Email,
		Address:     req.Address,
		City:        req.City,
		Province:    req.Province,
		PostalCode:  req.PostalCode,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		LogoURL:     req.LogoURL,
		BannerURL:   req.BannerURL,
	}

	store, err := h.storeUC.Update(c.Request.Context(), id, input, ownerID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "store updated successfully", store)
}

// List godoc
// @Summary List stores
// @Description List active stores with optional filters.
// @Tags Store
// @Produce json
// @Param city query string false "Filter by city"
// @Param province query string false "Filter by province"
// @Param search query string false "Search by name or city"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} response.Response{data=[]entity.Store}
// @Router /stores [get]
func (h *StoreHandler) List(c *gin.Context) {
	input := &usecase.ListStoreInput{
		City:     c.Query("city"),
		Province: c.Query("province"),
		Search:   c.Query("search"),
		Status:   "active", // Public listing shows only active stores
	}

	page, perPage := extractPagination(c)
	input.Page = page
	input.PerPage = perPage

	stores, meta, err := h.storeUC.List(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "stores retrieved", stores, meta)
}

// --- Photo Endpoints ---

// AddPhoto godoc
// @Summary Add store photo
// @Description Upload a photo for the store. Only the store owner can add photos.
// @Tags Store Photos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param body body request.AddStorePhotoRequest true "Photo payload"
// @Success 201 {object} response.Response{data=entity.StorePhoto}
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /stores/{id}/photos [post]
func (h *StoreHandler) AddPhoto(c *gin.Context) {
	storeID := c.Param("id")
	var req request.AddStorePhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ownerID := middleware.GetUserID(c)
	input := &usecase.AddStorePhotoInput{
		StoreID:   storeID,
		OwnerID:   ownerID,
		PhotoURL:  req.PhotoURL,
		Caption:   req.Caption,
		SortOrder: req.SortOrder,
		IsPrimary: req.IsPrimary,
	}

	photo, err := h.storeUC.AddPhoto(c.Request.Context(), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "photo added", photo)
}

// RemovePhoto godoc
// @Summary Remove store photo
// @Description Remove a photo from the store. Only the store owner can remove.
// @Tags Store Photos
// @Security BearerAuth
// @Produce json
// @Param id path string true "Store ID"
// @Param photoId path string true "Photo ID"
// @Success 200 {object} response.Response
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /stores/{id}/photos/{photoId} [delete]
func (h *StoreHandler) RemovePhoto(c *gin.Context) {
	storeID := c.Param("id")
	photoID := c.Param("photoId")
	ownerID := middleware.GetUserID(c)

	if err := h.storeUC.RemovePhoto(c.Request.Context(), storeID, photoID, ownerID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "photo removed", nil)
}

// SetPrimaryPhoto godoc
// @Summary Set primary store photo
// @Description Set a photo as the primary/cover photo for the store.
// @Tags Store Photos
// @Security BearerAuth
// @Produce json
// @Param id path string true "Store ID"
// @Param photoId path string true "Photo ID"
// @Success 200 {object} response.Response
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /stores/{id}/photos/{photoId}/primary [patch]
func (h *StoreHandler) SetPrimaryPhoto(c *gin.Context) {
	storeID := c.Param("id")
	photoID := c.Param("photoId")
	ownerID := middleware.GetUserID(c)

	if err := h.storeUC.SetPrimaryPhoto(c.Request.Context(), storeID, photoID, ownerID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "primary photo set", nil)
}

// --- Operating Hours Endpoints ---

// SetOperatingHours godoc
// @Summary Set store operating hours
// @Description Set operating hours for each day of the week. Replaces existing hours.
// @Tags Store
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param body body request.SetOperatingHoursRequest true "Operating hours payload"
// @Success 200 {object} response.Response
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /stores/{id}/operating-hours [put]
func (h *StoreHandler) SetOperatingHours(c *gin.Context) {
	storeID := c.Param("id")
	var req request.SetOperatingHoursRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	ownerID := middleware.GetUserID(c)
	hours := make([]usecase.OperatingHourInput, len(req.Hours))
	for i, h := range req.Hours {
		hours[i] = usecase.OperatingHourInput{
			DayOfWeek: h.DayOfWeek,
			OpenTime:  h.OpenTime,
			CloseTime: h.CloseTime,
			IsClosed:  h.IsClosed,
		}
	}

	if err := h.storeUC.SetOperatingHours(c.Request.Context(), storeID, ownerID, hours); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "operating hours updated", nil)
}

// GetOperatingHours godoc
// @Summary Get store operating hours
// @Description Get the operating hours for a store.
// @Tags Store
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} response.Response{data=[]entity.StoreOperatingHour}
// @Router /stores/{id}/operating-hours [get]
func (h *StoreHandler) GetOperatingHours(c *gin.Context) {
	storeID := c.Param("id")

	hours, err := h.storeUC.GetOperatingHours(c.Request.Context(), storeID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "operating hours retrieved", hours)
}

// --- Admin Endpoints ---

// ApproveStore godoc
// @Summary Approve a store (Admin)
// @Description Approve a pending store and set it to active.
// @Tags Store Admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /admin/stores/{id}/approve [patch]
func (h *StoreHandler) ApproveStore(c *gin.Context) {
	storeID := c.Param("id")
	adminID := middleware.GetUserID(c)

	if err := h.storeUC.ApproveStore(c.Request.Context(), storeID, adminID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "store approved", nil)
}

// SuspendStore godoc
// @Summary Suspend a store (Admin)
// @Description Suspend an active store with a reason.
// @Tags Store Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param body body request.SuspendStoreRequest true "Suspend reason"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /admin/stores/{id}/suspend [patch]
func (h *StoreHandler) SuspendStore(c *gin.Context) {
	storeID := c.Param("id")
	adminID := middleware.GetUserID(c)

	var req request.SuspendStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	if err := h.storeUC.SuspendStore(c.Request.Context(), storeID, req.Reason, adminID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "store suspended", nil)
}

// ReactivateStore godoc
// @Summary Reactivate a suspended store (Admin)
// @Description Reactivate a previously suspended store.
// @Tags Store Admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /admin/stores/{id}/reactivate [patch]
func (h *StoreHandler) ReactivateStore(c *gin.Context) {
	storeID := c.Param("id")
	adminID := middleware.GetUserID(c)

	if err := h.storeUC.ReactivateStore(c.Request.Context(), storeID, adminID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "store reactivated", nil)
}

// --- Helpers ---

func extractPagination(c *gin.Context) (int, int) {
	page := 1
	perPage := 20

	if p := c.Query("page"); p != "" {
		if v, err := parseInt(p); err == nil && v > 0 {
			page = v
		}
	}
	if pp := c.Query("per_page"); pp != "" {
		if v, err := parseInt(pp); err == nil && v > 0 {
			perPage = v
		}
	}
	return page, perPage
}

func parseInt(s string) (int, error) {
	var v int
	_, err := fmt.Sscanf(s, "%d", &v)
	return v, err
}
