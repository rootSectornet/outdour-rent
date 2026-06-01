package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rentoutdoor/api/internal/adapter/handler"
	"github.com/rentoutdoor/api/internal/adapter/middleware"
	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/usecase"
	"github.com/rentoutdoor/api/pkg/pagination"
)

// --- Mock StoreUsecase ---

type mockStoreUsecase struct {
	createFn            func(ctx context.Context, input *usecase.CreateStoreInput) (*entity.Store, error)
	getByIDFn           func(ctx context.Context, id string) (*entity.Store, error)
	getBySlugFn         func(ctx context.Context, slug string) (*entity.Store, error)
	getMyStoreFn        func(ctx context.Context, ownerID string) (*entity.Store, error)
	updateFn            func(ctx context.Context, id string, input *usecase.UpdateStoreInput, ownerID string) (*entity.Store, error)
	listFn              func(ctx context.Context, input *usecase.ListStoreInput) ([]entity.Store, *pagination.Meta, error)
	addPhotoFn          func(ctx context.Context, input *usecase.AddStorePhotoInput) (*entity.StorePhoto, error)
	removePhotoFn       func(ctx context.Context, storeID, photoID, ownerID string) error
	setPrimaryPhotoFn   func(ctx context.Context, storeID, photoID, ownerID string) error
	setOperatingHoursFn func(ctx context.Context, storeID, ownerID string, hours []usecase.OperatingHourInput) error
	getOperatingHoursFn func(ctx context.Context, storeID string) ([]entity.StoreOperatingHour, error)
	approveStoreFn      func(ctx context.Context, storeID, adminID string) error
	suspendStoreFn      func(ctx context.Context, storeID, reason, adminID string) error
	reactivateStoreFn   func(ctx context.Context, storeID, adminID string) error
}

func (m *mockStoreUsecase) Create(ctx context.Context, input *usecase.CreateStoreInput) (*entity.Store, error) {
	return m.createFn(ctx, input)
}
func (m *mockStoreUsecase) GetByID(ctx context.Context, id string) (*entity.Store, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockStoreUsecase) GetBySlug(ctx context.Context, slug string) (*entity.Store, error) {
	return m.getBySlugFn(ctx, slug)
}
func (m *mockStoreUsecase) GetMyStore(ctx context.Context, ownerID string) (*entity.Store, error) {
	return m.getMyStoreFn(ctx, ownerID)
}
func (m *mockStoreUsecase) Update(ctx context.Context, id string, input *usecase.UpdateStoreInput, ownerID string) (*entity.Store, error) {
	return m.updateFn(ctx, id, input, ownerID)
}
func (m *mockStoreUsecase) List(ctx context.Context, input *usecase.ListStoreInput) ([]entity.Store, *pagination.Meta, error) {
	return m.listFn(ctx, input)
}
func (m *mockStoreUsecase) AddPhoto(ctx context.Context, input *usecase.AddStorePhotoInput) (*entity.StorePhoto, error) {
	return m.addPhotoFn(ctx, input)
}
func (m *mockStoreUsecase) RemovePhoto(ctx context.Context, storeID, photoID, ownerID string) error {
	return m.removePhotoFn(ctx, storeID, photoID, ownerID)
}
func (m *mockStoreUsecase) SetPrimaryPhoto(ctx context.Context, storeID, photoID, ownerID string) error {
	return m.setPrimaryPhotoFn(ctx, storeID, photoID, ownerID)
}
func (m *mockStoreUsecase) SetOperatingHours(ctx context.Context, storeID, ownerID string, hours []usecase.OperatingHourInput) error {
	return m.setOperatingHoursFn(ctx, storeID, ownerID, hours)
}
func (m *mockStoreUsecase) GetOperatingHours(ctx context.Context, storeID string) ([]entity.StoreOperatingHour, error) {
	return m.getOperatingHoursFn(ctx, storeID)
}
func (m *mockStoreUsecase) ApproveStore(ctx context.Context, storeID, adminID string) error {
	return m.approveStoreFn(ctx, storeID, adminID)
}
func (m *mockStoreUsecase) SuspendStore(ctx context.Context, storeID, reason, adminID string) error {
	return m.suspendStoreFn(ctx, storeID, reason, adminID)
}
func (m *mockStoreUsecase) ReactivateStore(ctx context.Context, storeID, adminID string) error {
	return m.reactivateStoreFn(ctx, storeID, adminID)
}

// --- Helpers ---

func setupStoreRouter(storeUC usecase.StoreUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewStoreHandler(storeUC)

	stores := r.Group("/api/v1/stores")
	{
		stores.GET("", h.List)
		stores.GET("/:id", h.GetByID)
		stores.GET("/slug/:slug", h.GetBySlug)
		stores.GET("/:id/operating-hours", h.GetOperatingHours)

		// Simulate auth middleware by setting user context
		owner := stores.Group("")
		owner.Use(func(c *gin.Context) {
			c.Set(middleware.ContextKeyUserID, "owner-1")
			c.Set(middleware.ContextKeyUserRole, "owner")
			c.Next()
		})
		{
			owner.POST("", h.Create)
			owner.GET("/me", h.GetMyStore)
			owner.PUT("/:id", h.Update)
			owner.POST("/:id/photos", h.AddPhoto)
			owner.DELETE("/:id/photos/:photoId", h.RemovePhoto)
			owner.PATCH("/:id/photos/:photoId/primary", h.SetPrimaryPhoto)
			owner.PUT("/:id/operating-hours", h.SetOperatingHours)
		}
	}

	admin := r.Group("/api/v1/admin/stores")
	admin.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyUserID, "admin-1")
		c.Set(middleware.ContextKeyUserRole, "admin")
		c.Next()
	})
	{
		admin.PATCH("/:id/approve", h.ApproveStore)
		admin.PATCH("/:id/suspend", h.SuspendStore)
		admin.PATCH("/:id/reactivate", h.ReactivateStore)
	}

	return r
}

func testStore() *entity.Store {
	return &entity.Store{
		BaseModel: entity.BaseModel{ID: "store-1"},
		OwnerID:   "owner-1",
		Name:      "Test Store",
		Slug:      "test-store",
		Phone:     "08123456789",
		Email:     "store@example.com",
		Address:   "Jl. Test",
		City:      "Bandung",
		Province:  "Jawa Barat",
		Status:    entity.StoreStatusActive,
	}
}

// --- Tests ---

func TestStoreHandler_Create_Success(t *testing.T) {
	mock := &mockStoreUsecase{
		createFn: func(_ context.Context, input *usecase.CreateStoreInput) (*entity.Store, error) {
			if input.Name != "My Store" {
				t.Fatalf("expected name My Store, got %s", input.Name)
			}
			return testStore(), nil
		},
	}

	r := setupStoreRouter(mock)
	body := map[string]interface{}{
		"name":     "My Store",
		"phone":    "08123456789",
		"email":    "store@example.com",
		"address":  "Jl. Test No. 123 Bandung",
		"city":     "Bandung",
		"province": "Jawa Barat",
	}

	w := performRequest(r, "POST", "/api/v1/stores", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStoreHandler_GetByID_Success(t *testing.T) {
	mock := &mockStoreUsecase{
		getByIDFn: func(_ context.Context, id string) (*entity.Store, error) {
			if id != "store-1" {
				t.Fatalf("expected store-1, got %s", id)
			}
			return testStore(), nil
		},
	}

	r := setupStoreRouter(mock)
	w := performRequest(r, "GET", "/api/v1/stores/store-1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStoreHandler_GetByID_NotFound(t *testing.T) {
	mock := &mockStoreUsecase{
		getByIDFn: func(_ context.Context, _ string) (*entity.Store, error) {
			return nil, usecase.ErrNotFound
		},
	}

	r := setupStoreRouter(mock)
	w := performRequest(r, "GET", "/api/v1/stores/nonexistent", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStoreHandler_List_Success(t *testing.T) {
	mock := &mockStoreUsecase{
		listFn: func(_ context.Context, input *usecase.ListStoreInput) ([]entity.Store, *pagination.Meta, error) {
			if input.Status != "active" {
				t.Fatalf("expected status active, got %s", input.Status)
			}
			return []entity.Store{*testStore()}, pagination.NewMeta(1, 1, 20), nil
		},
	}

	r := setupStoreRouter(mock)
	w := performRequest(r, "GET", "/api/v1/stores?city=Bandung", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["meta"] == nil {
		t.Fatal("expected meta in response")
	}
}

func TestStoreHandler_Update_Success(t *testing.T) {
	mock := &mockStoreUsecase{
		updateFn: func(_ context.Context, id string, input *usecase.UpdateStoreInput, ownerID string) (*entity.Store, error) {
			s := testStore()
			if input.Name != nil {
				s.Name = *input.Name
			}
			return s, nil
		},
	}

	r := setupStoreRouter(mock)
	body := map[string]interface{}{
		"name": "Updated Name",
	}
	w := performRequest(r, "PUT", "/api/v1/stores/store-1", body)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStoreHandler_AddPhoto_Success(t *testing.T) {
	mock := &mockStoreUsecase{
		addPhotoFn: func(_ context.Context, input *usecase.AddStorePhotoInput) (*entity.StorePhoto, error) {
			return &entity.StorePhoto{
				BaseModelCreateOnly: entity.BaseModelCreateOnly{ID: "photo-1"},
				StoreID:             input.StoreID,
				PhotoURL:            input.PhotoURL,
				IsPrimary:           input.IsPrimary,
			}, nil
		},
	}

	r := setupStoreRouter(mock)
	body := map[string]interface{}{
		"photo_url": "https://cdn.example.com/photo.jpg",
		"caption":   "Front view",
	}
	w := performRequest(r, "POST", "/api/v1/stores/store-1/photos", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStoreHandler_SetOperatingHours_Success(t *testing.T) {
	mock := &mockStoreUsecase{
		setOperatingHoursFn: func(_ context.Context, storeID, ownerID string, hours []usecase.OperatingHourInput) error {
			if len(hours) != 2 {
				t.Fatalf("expected 2 hours, got %d", len(hours))
			}
			return nil
		},
	}

	r := setupStoreRouter(mock)
	body := map[string]interface{}{
		"hours": []map[string]interface{}{
			{"day_of_week": 1, "open_time": "08:00", "close_time": "17:00", "is_closed": false},
			{"day_of_week": 7, "open_time": "00:00", "close_time": "00:00", "is_closed": true},
		},
	}
	w := performRequest(r, "PUT", "/api/v1/stores/store-1/operating-hours", body)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStoreHandler_ApproveStore_Success(t *testing.T) {
	mock := &mockStoreUsecase{
		approveStoreFn: func(_ context.Context, storeID, adminID string) error {
			if adminID != "admin-1" {
				t.Fatalf("expected admin-1, got %s", adminID)
			}
			return nil
		},
	}

	r := setupStoreRouter(mock)
	req := httptest.NewRequest("PATCH", "/api/v1/admin/stores/store-1/approve", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStoreHandler_SuspendStore_Success(t *testing.T) {
	mock := &mockStoreUsecase{
		suspendStoreFn: func(_ context.Context, storeID, reason, adminID string) error {
			if reason != "Terms violation" {
				t.Fatalf("expected reason, got %s", reason)
			}
			return nil
		},
	}

	r := setupStoreRouter(mock)
	body := map[string]interface{}{
		"reason": "Terms violation",
	}
	w := performRequest(r, "PATCH", "/api/v1/admin/stores/store-1/suspend", body)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
