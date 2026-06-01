package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/repository"
	"github.com/rentoutdoor/api/internal/usecase"
	storeuc "github.com/rentoutdoor/api/internal/usecase/store"
	"github.com/rentoutdoor/api/pkg/pagination"
	"gorm.io/gorm"
)

// --- Mocks ---

type mockStoreRepo struct {
	stores  map[string]*entity.Store
	byOwner map[string]*entity.Store
	bySlug  map[string]*entity.Store
}

func newMockStoreRepo() *mockStoreRepo {
	return &mockStoreRepo{
		stores:  make(map[string]*entity.Store),
		byOwner: make(map[string]*entity.Store),
		bySlug:  make(map[string]*entity.Store),
	}
}

func (m *mockStoreRepo) Create(_ context.Context, s *entity.Store) error {
	m.stores[s.ID] = s
	m.byOwner[s.OwnerID] = s
	m.bySlug[s.Slug] = s
	return nil
}

func (m *mockStoreRepo) FindByID(_ context.Context, id string) (*entity.Store, error) {
	if s, ok := m.stores[id]; ok {
		return s, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockStoreRepo) FindByIDWithRelations(_ context.Context, id string) (*entity.Store, error) {
	if s, ok := m.stores[id]; ok {
		return s, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockStoreRepo) FindByOwnerID(_ context.Context, ownerID string) (*entity.Store, error) {
	if s, ok := m.byOwner[ownerID]; ok {
		return s, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockStoreRepo) FindBySlug(_ context.Context, slug string) (*entity.Store, error) {
	if s, ok := m.bySlug[slug]; ok {
		return s, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockStoreRepo) SlugExists(_ context.Context, slug string) (bool, error) {
	_, exists := m.bySlug[slug]
	return exists, nil
}

func (m *mockStoreRepo) Update(_ context.Context, s *entity.Store) error {
	m.stores[s.ID] = s
	m.byOwner[s.OwnerID] = s
	m.bySlug[s.Slug] = s
	return nil
}

func (m *mockStoreRepo) Delete(_ context.Context, id string, _ string) error {
	if s, ok := m.stores[id]; ok {
		delete(m.byOwner, s.OwnerID)
		delete(m.bySlug, s.Slug)
		delete(m.stores, id)
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockStoreRepo) List(_ context.Context, params *repository.StoreListParams) ([]entity.Store, *pagination.Meta, error) {
	var result []entity.Store
	for _, s := range m.stores {
		if params.Status != "" && string(s.Status) != params.Status {
			continue
		}
		result = append(result, *s)
	}
	meta := pagination.NewMeta(int64(len(result)), params.Page, params.PerPage)
	return result, meta, nil
}

type mockPhotoRepo struct {
	photos map[string]*entity.StorePhoto
}

func newMockPhotoRepo() *mockPhotoRepo {
	return &mockPhotoRepo{photos: make(map[string]*entity.StorePhoto)}
}

func (m *mockPhotoRepo) Create(_ context.Context, p *entity.StorePhoto) error {
	m.photos[p.ID] = p
	return nil
}

func (m *mockPhotoRepo) FindByStoreID(_ context.Context, storeID string) ([]entity.StorePhoto, error) {
	var result []entity.StorePhoto
	for _, p := range m.photos {
		if p.StoreID == storeID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockPhotoRepo) Delete(_ context.Context, id string) error {
	delete(m.photos, id)
	return nil
}

func (m *mockPhotoRepo) SetPrimary(_ context.Context, storeID string, photoID string) error {
	for _, p := range m.photos {
		if p.StoreID == storeID {
			p.IsPrimary = (p.ID == photoID)
		}
	}
	return nil
}

func (m *mockPhotoRepo) UpdateSortOrder(_ context.Context, id string, order int) error {
	if p, ok := m.photos[id]; ok {
		p.SortOrder = order
	}
	return nil
}

type mockHoursRepo struct {
	hours map[string][]entity.StoreOperatingHour
}

func newMockHoursRepo() *mockHoursRepo {
	return &mockHoursRepo{hours: make(map[string][]entity.StoreOperatingHour)}
}

func (m *mockHoursRepo) Upsert(_ context.Context, hours []entity.StoreOperatingHour) error {
	if len(hours) == 0 {
		return nil
	}
	storeID := hours[0].StoreID
	m.hours[storeID] = hours
	return nil
}

func (m *mockHoursRepo) FindByStoreID(_ context.Context, storeID string) ([]entity.StoreOperatingHour, error) {
	return m.hours[storeID], nil
}

func (m *mockHoursRepo) DeleteByStoreID(_ context.Context, storeID string) error {
	delete(m.hours, storeID)
	return nil
}

type mockUserRepo struct {
	users map[string]*entity.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*entity.User)}
}

func (m *mockUserRepo) Create(_ context.Context, u *entity.User) error {
	m.users[u.ID] = u
	return nil
}
func (m *mockUserRepo) FindByID(_ context.Context, id string) (*entity.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *mockUserRepo) FindByEmail(_ context.Context, _ string) (*entity.User, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockUserRepo) FindByGoogleID(_ context.Context, _ string) (*entity.User, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockUserRepo) Update(_ context.Context, u *entity.User) error {
	m.users[u.ID] = u
	return nil
}
func (m *mockUserRepo) Delete(_ context.Context, id string, _ string) error {
	delete(m.users, id)
	return nil
}

// --- Helpers ---

func setupStoreUsecase() (usecase.StoreUsecase, *mockStoreRepo, *mockPhotoRepo, *mockHoursRepo, *mockUserRepo) {
	storeRepo := newMockStoreRepo()
	photoRepo := newMockPhotoRepo()
	hoursRepo := newMockHoursRepo()
	userRepo := newMockUserRepo()

	uc := storeuc.NewStoreUsecase(storeRepo, photoRepo, hoursRepo, userRepo)
	return uc, storeRepo, photoRepo, hoursRepo, userRepo
}

func createOwnerUser(userRepo *mockUserRepo) *entity.User {
	user := &entity.User{
		BaseModel: entity.BaseModel{ID: "owner-1"},
		Email:     "owner@example.com",
		FullName:  "Store Owner",
		Role:      entity.UserRoleOwner,
		IsActive:  true,
	}
	userRepo.users[user.ID] = user
	return user
}

func createActiveStore(storeRepo *mockStoreRepo, ownerID string) *entity.Store {
	s := &entity.Store{
		BaseModel: entity.BaseModel{ID: "store-1"},
		OwnerID:   ownerID,
		Name:      "Test Store",
		Slug:      "test-store",
		Phone:     "08123456789",
		Email:     "store@example.com",
		Address:   "Jl. Test No. 1",
		City:      "Bandung",
		Province:  "Jawa Barat",
		Status:    entity.StoreStatusActive,
	}
	storeRepo.stores[s.ID] = s
	storeRepo.byOwner[s.OwnerID] = s
	storeRepo.bySlug[s.Slug] = s
	return s
}

// --- Tests ---

func TestCreateStore_Success(t *testing.T) {
	uc, storeRepo, _, _, userRepo := setupStoreUsecase()
	ctx := context.Background()
	createOwnerUser(userRepo)

	input := &usecase.CreateStoreInput{
		OwnerID:     "owner-1",
		Name:        "My Outdoor Store",
		Description: "Best gear in town",
		Phone:       "08123456789",
		Email:       "mystore@example.com",
		Address:     "Jl. Outdoor No. 1",
		City:        "Bandung",
		Province:    "Jawa Barat",
	}

	store, err := uc.Create(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.Name != "My Outdoor Store" {
		t.Fatalf("expected name My Outdoor Store, got %s", store.Name)
	}
	if store.Status != entity.StoreStatusPendingApproval {
		t.Fatalf("expected status pending_approval, got %s", store.Status)
	}
	if store.Slug == "" {
		t.Fatal("expected slug to be generated")
	}
	if len(storeRepo.stores) != 1 {
		t.Fatalf("expected 1 store, got %d", len(storeRepo.stores))
	}
}

func TestCreateStore_DuplicateOwner(t *testing.T) {
	uc, storeRepo, _, _, userRepo := setupStoreUsecase()
	ctx := context.Background()
	createOwnerUser(userRepo)
	createActiveStore(storeRepo, "owner-1")

	input := &usecase.CreateStoreInput{
		OwnerID: "owner-1",
		Name:    "Second Store",
		Phone:   "08123456789",
		Email:   "store2@example.com",
		Address: "Jl. Second No. 2",
		City:    "Jakarta",
		Province: "DKI Jakarta",
	}

	_, err := uc.Create(ctx, input)
	if err == nil {
		t.Fatal("expected error for duplicate owner")
	}
	if !errors.Is(err, usecase.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestCreateStore_RenterForbidden(t *testing.T) {
	uc, _, _, _, userRepo := setupStoreUsecase()
	ctx := context.Background()

	// Create a renter user
	userRepo.users["renter-1"] = &entity.User{
		BaseModel: entity.BaseModel{ID: "renter-1"},
		Role:      entity.UserRoleRenter,
		IsActive:  true,
	}

	input := &usecase.CreateStoreInput{
		OwnerID: "renter-1",
		Name:    "Forbidden Store",
		Phone:   "08123456789",
		Email:   "nope@example.com",
		Address: "Jl. Forbidden No. 1",
		City:    "Bandung",
		Province: "Jawa Barat",
	}

	_, err := uc.Create(ctx, input)
	if err == nil {
		t.Fatal("expected error for renter creating store")
	}
	if !errors.Is(err, usecase.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestUpdateStore_Success(t *testing.T) {
	uc, storeRepo, _, _, userRepo := setupStoreUsecase()
	ctx := context.Background()
	createOwnerUser(userRepo)
	createActiveStore(storeRepo, "owner-1")

	newName := "Updated Store Name"
	input := &usecase.UpdateStoreInput{
		Name: &newName,
	}

	store, err := uc.Update(ctx, "store-1", input, "owner-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.Name != "Updated Store Name" {
		t.Fatalf("expected Updated Store Name, got %s", store.Name)
	}
}

func TestUpdateStore_NotOwner(t *testing.T) {
	uc, storeRepo, _, _, userRepo := setupStoreUsecase()
	ctx := context.Background()
	createOwnerUser(userRepo)
	createActiveStore(storeRepo, "owner-1")

	newName := "Hacked"
	input := &usecase.UpdateStoreInput{Name: &newName}

	_, err := uc.Update(ctx, "store-1", input, "other-user")
	if err == nil {
		t.Fatal("expected error for non-owner update")
	}
	if !errors.Is(err, usecase.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestGetByID_Success(t *testing.T) {
	uc, storeRepo, _, _, _ := setupStoreUsecase()
	ctx := context.Background()
	createActiveStore(storeRepo, "owner-1")

	store, err := uc.GetByID(ctx, "store-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if store.ID != "store-1" {
		t.Fatalf("expected store-1, got %s", store.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	uc, _, _, _, _ := setupStoreUsecase()
	ctx := context.Background()

	_, err := uc.GetByID(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing store")
	}
	if !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestAddPhoto_Success(t *testing.T) {
	uc, storeRepo, photoRepo, _, _ := setupStoreUsecase()
	ctx := context.Background()
	createActiveStore(storeRepo, "owner-1")

	input := &usecase.AddStorePhotoInput{
		StoreID:   "store-1",
		OwnerID:   "owner-1",
		PhotoURL:  "https://cdn.example.com/photo.jpg",
		Caption:   "Front view",
		SortOrder: 1,
		IsPrimary: true,
	}

	photo, err := uc.AddPhoto(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if photo.PhotoURL != "https://cdn.example.com/photo.jpg" {
		t.Fatalf("expected photo URL, got %s", photo.PhotoURL)
	}
	if len(photoRepo.photos) != 1 {
		t.Fatalf("expected 1 photo, got %d", len(photoRepo.photos))
	}
}

func TestAddPhoto_NotOwner(t *testing.T) {
	uc, storeRepo, _, _, _ := setupStoreUsecase()
	ctx := context.Background()
	createActiveStore(storeRepo, "owner-1")

	input := &usecase.AddStorePhotoInput{
		StoreID:  "store-1",
		OwnerID:  "hacker",
		PhotoURL: "https://evil.com/photo.jpg",
	}

	_, err := uc.AddPhoto(ctx, input)
	if err == nil {
		t.Fatal("expected error for non-owner adding photo")
	}
	if !errors.Is(err, usecase.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestSetOperatingHours_Success(t *testing.T) {
	uc, storeRepo, _, hoursRepo, _ := setupStoreUsecase()
	ctx := context.Background()
	createActiveStore(storeRepo, "owner-1")

	hours := []usecase.OperatingHourInput{
		{DayOfWeek: 1, OpenTime: "08:00", CloseTime: "17:00", IsClosed: false},
		{DayOfWeek: 2, OpenTime: "08:00", CloseTime: "17:00", IsClosed: false},
		{DayOfWeek: 7, OpenTime: "00:00", CloseTime: "00:00", IsClosed: true},
	}

	err := uc.SetOperatingHours(ctx, "store-1", "owner-1", hours)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	stored := hoursRepo.hours["store-1"]
	if len(stored) != 3 {
		t.Fatalf("expected 3 hours entries, got %d", len(stored))
	}
}

func TestSetOperatingHours_DuplicateDay(t *testing.T) {
	uc, storeRepo, _, _, _ := setupStoreUsecase()
	ctx := context.Background()
	createActiveStore(storeRepo, "owner-1")

	hours := []usecase.OperatingHourInput{
		{DayOfWeek: 1, OpenTime: "08:00", CloseTime: "17:00"},
		{DayOfWeek: 1, OpenTime: "09:00", CloseTime: "18:00"}, // duplicate!
	}

	err := uc.SetOperatingHours(ctx, "store-1", "owner-1", hours)
	if err == nil {
		t.Fatal("expected error for duplicate day")
	}
	if !errors.Is(err, usecase.ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestApproveStore_Success(t *testing.T) {
	uc, storeRepo, _, _, _ := setupStoreUsecase()
	ctx := context.Background()

	// Create pending store
	s := &entity.Store{
		BaseModel: entity.BaseModel{ID: "store-pending"},
		OwnerID:   "owner-1",
		Name:      "Pending Store",
		Slug:      "pending-store",
		Status:    entity.StoreStatusPendingApproval,
	}
	storeRepo.stores[s.ID] = s

	err := uc.ApproveStore(ctx, "store-pending", "admin-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if storeRepo.stores["store-pending"].Status != entity.StoreStatusActive {
		t.Fatalf("expected active, got %s", storeRepo.stores["store-pending"].Status)
	}
	if storeRepo.stores["store-pending"].VerifiedAt == nil {
		t.Fatal("expected verified_at to be set")
	}
}

func TestApproveStore_InvalidTransition(t *testing.T) {
	uc, storeRepo, _, _, _ := setupStoreUsecase()
	ctx := context.Background()

	// Already active store
	s := &entity.Store{
		BaseModel: entity.BaseModel{ID: "store-active"},
		OwnerID:   "owner-1",
		Status:    entity.StoreStatusActive,
	}
	storeRepo.stores[s.ID] = s

	err := uc.ApproveStore(ctx, "store-active", "admin-1")
	if err == nil {
		t.Fatal("expected error for invalid transition")
	}
	if !errors.Is(err, usecase.ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestSuspendStore_Success(t *testing.T) {
	uc, storeRepo, _, _, _ := setupStoreUsecase()
	ctx := context.Background()
	createActiveStore(storeRepo, "owner-1")

	err := uc.SuspendStore(ctx, "store-1", "Terms violation", "admin-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if storeRepo.stores["store-1"].Status != entity.StoreStatusSuspended {
		t.Fatalf("expected suspended, got %s", storeRepo.stores["store-1"].Status)
	}
	if storeRepo.stores["store-1"].SuspendedAt == nil {
		t.Fatal("expected suspended_at to be set")
	}
}

func TestReactivateStore_Success(t *testing.T) {
	uc, storeRepo, _, _, _ := setupStoreUsecase()
	ctx := context.Background()

	// Create suspended store
	s := &entity.Store{
		BaseModel: entity.BaseModel{ID: "store-suspended"},
		OwnerID:   "owner-1",
		Status:    entity.StoreStatusSuspended,
	}
	storeRepo.stores[s.ID] = s

	err := uc.ReactivateStore(ctx, "store-suspended", "admin-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if storeRepo.stores["store-suspended"].Status != entity.StoreStatusActive {
		t.Fatalf("expected active, got %s", storeRepo.stores["store-suspended"].Status)
	}
}

func TestReactivateStore_InvalidTransition(t *testing.T) {
	uc, storeRepo, _, _, _ := setupStoreUsecase()
	ctx := context.Background()

	// Pending store can't be reactivated
	s := &entity.Store{
		BaseModel: entity.BaseModel{ID: "store-pending"},
		OwnerID:   "owner-1",
		Status:    entity.StoreStatusPendingApproval,
	}
	storeRepo.stores[s.ID] = s

	err := uc.ReactivateStore(ctx, "store-pending", "admin-1")
	if err == nil {
		t.Fatal("expected error for invalid transition")
	}
	if !errors.Is(err, usecase.ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestList_FiltersActive(t *testing.T) {
	uc, storeRepo, _, _, _ := setupStoreUsecase()
	ctx := context.Background()

	// Active store
	storeRepo.stores["s1"] = &entity.Store{
		BaseModel: entity.BaseModel{ID: "s1"},
		Status:    entity.StoreStatusActive,
		City:      "Bandung",
	}
	// Pending store (shouldn't appear in public listing)
	storeRepo.stores["s2"] = &entity.Store{
		BaseModel: entity.BaseModel{ID: "s2"},
		Status:    entity.StoreStatusPendingApproval,
		City:      "Bandung",
	}

	input := &usecase.ListStoreInput{
		Status:  "active",
		Page:    1,
		PerPage: 20,
	}

	stores, _, err := uc.List(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(stores) != 1 {
		t.Fatalf("expected 1 active store, got %d", len(stores))
	}
	if stores[0].ID != "s1" {
		t.Fatalf("expected s1, got %s", stores[0].ID)
	}
}

func TestGetMyStore_NotFound(t *testing.T) {
	uc, _, _, _, _ := setupStoreUsecase()
	ctx := context.Background()

	_, err := uc.GetMyStore(ctx, "nobody")
	if err == nil {
		t.Fatal("expected error for no store")
	}
	if !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
