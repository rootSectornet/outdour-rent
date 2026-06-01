package equipment_test

import (
	"context"
	"testing"
	"time"

	"github.com/rentoutdoor/api/internal/domain/entity"
	"github.com/rentoutdoor/api/internal/repository"
	"github.com/rentoutdoor/api/internal/usecase"
	equipmentuc "github.com/rentoutdoor/api/internal/usecase/equipment"
	"github.com/rentoutdoor/api/pkg/pagination"
	"gorm.io/gorm"
)

// --- Mocks ---

type mockEquipmentRepo struct {
	items map[string]*entity.Equipment
}

func newMockEquipmentRepo() *mockEquipmentRepo {
	return &mockEquipmentRepo{items: make(map[string]*entity.Equipment)}
}

func (m *mockEquipmentRepo) Create(_ context.Context, e *entity.Equipment) error {
	if e.ID == "" {
		e.ID = "equip-" + e.Slug
	}
	m.items[e.ID] = e
	return nil
}

func (m *mockEquipmentRepo) FindByID(_ context.Context, id string) (*entity.Equipment, error) {
	if e, ok := m.items[id]; ok {
		return e, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockEquipmentRepo) FindByIDForUpdate(_ context.Context, _ *gorm.DB, id string) (*entity.Equipment, error) {
	if e, ok := m.items[id]; ok {
		return e, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockEquipmentRepo) Update(_ context.Context, e *entity.Equipment) error {
	m.items[e.ID] = e
	return nil
}

func (m *mockEquipmentRepo) Delete(_ context.Context, id string, _ string) error {
	if _, ok := m.items[id]; ok {
		delete(m.items, id)
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockEquipmentRepo) List(_ context.Context, params *repository.EquipmentListParams) ([]entity.Equipment, *pagination.Meta, error) {
	var result []entity.Equipment
	for _, e := range m.items {
		if params.Status != "" && string(e.Status) != params.Status {
			continue
		}
		if params.StoreID != "" && e.StoreID != params.StoreID {
			continue
		}
		if params.IsActive != nil && e.IsActive != *params.IsActive {
			continue
		}
		result = append(result, *e)
	}
	meta := pagination.NewMeta(int64(len(result)), params.Page, params.PerPage)
	return result, meta, nil
}

type mockCategoryRepo struct {
	categories map[string]*entity.EquipmentCategory
}

func newMockCategoryRepo() *mockCategoryRepo {
	return &mockCategoryRepo{categories: make(map[string]*entity.EquipmentCategory)}
}

func (m *mockCategoryRepo) Create(_ context.Context, c *entity.EquipmentCategory) error {
	m.categories[c.ID] = c
	return nil
}
func (m *mockCategoryRepo) FindByID(_ context.Context, id string) (*entity.EquipmentCategory, error) {
	if c, ok := m.categories[id]; ok {
		return c, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *mockCategoryRepo) FindBySlug(_ context.Context, _ string) (*entity.EquipmentCategory, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockCategoryRepo) Update(_ context.Context, c *entity.EquipmentCategory) error {
	m.categories[c.ID] = c
	return nil
}
func (m *mockCategoryRepo) Delete(_ context.Context, id string, _ string) error {
	delete(m.categories, id)
	return nil
}
func (m *mockCategoryRepo) ListRoots(_ context.Context) ([]entity.EquipmentCategory, error) {
	return nil, nil
}
func (m *mockCategoryRepo) ListByParent(_ context.Context, _ string) ([]entity.EquipmentCategory, error) {
	return nil, nil
}

type mockReservationRepo struct{}

func (m *mockReservationRepo) Create(_ context.Context, _ *gorm.DB, _ *entity.InventoryReservation) error {
	return nil
}
func (m *mockReservationRepo) FindByID(_ context.Context, _ string) (*entity.InventoryReservation, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockReservationRepo) Update(_ context.Context, _ *entity.InventoryReservation) error {
	return nil
}
func (m *mockReservationRepo) UpdateStatus(_ context.Context, _ *gorm.DB, _ string, _ entity.ReservationStatus, _ string) error {
	return nil
}
func (m *mockReservationRepo) GetOverlappingReservations(_ context.Context, _ *gorm.DB, _ string, _, _ string) ([]entity.InventoryReservation, error) {
	return nil, nil
}
func (m *mockReservationRepo) GetPeakUsage(_ context.Context, _ *gorm.DB, _ string, _, _ string) (uint, error) {
	return 0, nil
}
func (m *mockReservationRepo) CreateDateLocks(_ context.Context, _ *gorm.DB, _ []entity.ReservationDateLock) error {
	return nil
}
func (m *mockReservationRepo) DeleteDateLocksByReservation(_ context.Context, _ *gorm.DB, _ string) error {
	return nil
}
func (m *mockReservationRepo) FindExpired(_ context.Context, _ time.Time, _ int) ([]entity.InventoryReservation, error) {
	return nil, nil
}
func (m *mockReservationRepo) ExpireReservation(_ context.Context, _ *gorm.DB, _ string) error {
	return nil
}

type mockStoreRepoForEquip struct {
	byOwner map[string]*entity.Store
}

func newMockStoreRepoForEquip() *mockStoreRepoForEquip {
	return &mockStoreRepoForEquip{byOwner: make(map[string]*entity.Store)}
}

func (m *mockStoreRepoForEquip) Create(_ context.Context, _ *entity.Store) error { return nil }
func (m *mockStoreRepoForEquip) FindByID(_ context.Context, _ string) (*entity.Store, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockStoreRepoForEquip) FindByIDWithRelations(_ context.Context, _ string) (*entity.Store, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockStoreRepoForEquip) FindByOwnerID(_ context.Context, ownerID string) (*entity.Store, error) {
	if s, ok := m.byOwner[ownerID]; ok {
		return s, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *mockStoreRepoForEquip) FindBySlug(_ context.Context, _ string) (*entity.Store, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockStoreRepoForEquip) SlugExists(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (m *mockStoreRepoForEquip) Update(_ context.Context, _ *entity.Store) error { return nil }
func (m *mockStoreRepoForEquip) Delete(_ context.Context, _ string, _ string) error {
	return nil
}
func (m *mockStoreRepoForEquip) List(_ context.Context, _ *repository.StoreListParams) ([]entity.Store, *pagination.Meta, error) {
	return nil, nil, nil
}

type mockTxManager struct{}

func (m *mockTxManager) WithTransaction(_ context.Context, fn func(tx *gorm.DB) error) error {
	return fn(nil)
}
func (m *mockTxManager) WithTransactionIsolation(_ context.Context, _ string, fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

// --- Helpers ---

func setupEquipmentUsecase() (usecase.EquipmentUsecase, *mockEquipmentRepo, *mockCategoryRepo, *mockStoreRepoForEquip) {
	equipRepo := newMockEquipmentRepo()
	catRepo := newMockCategoryRepo()
	storeRepo := newMockStoreRepoForEquip()
	resRepo := &mockReservationRepo{}
	txMgr := &mockTxManager{}

	uc := equipmentuc.NewEquipmentUsecase(equipRepo, catRepo, resRepo, storeRepo, txMgr)
	return uc, equipRepo, catRepo, storeRepo
}

func seedTestData(equipRepo *mockEquipmentRepo, catRepo *mockCategoryRepo, storeRepo *mockStoreRepoForEquip) {
	storeRepo.byOwner["owner-1"] = &entity.Store{
		BaseModel: entity.BaseModel{ID: "store-1"},
		OwnerID:   "owner-1",
		Name:      "Test Store",
		Status:    entity.StoreStatusActive,
	}

	catRepo.categories["cat-1"] = &entity.EquipmentCategory{
		BaseModel: entity.BaseModel{ID: "cat-1"},
		Name:      "Tenda",
		Slug:      "tenda",
		IsActive:  true,
	}

	equipRepo.items["equip-1"] = &entity.Equipment{
		BaseModel:   entity.BaseModel{ID: "equip-1"},
		StoreID:     "store-1",
		CategoryID:  "cat-1",
		Name:        "Tenda Dome 4P",
		Slug:        "tenda-dome-4p",
		TotalStock:  5,
		Condition:   entity.EquipmentConditionNew,
		Status:      entity.EquipmentStatusAvailable,
		IsActive:    true,
		DepositAmount: 100000,
	}
}

// --- Tests ---

func TestCreateEquipment_Success(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()

	storeRepo.byOwner["owner-1"] = &entity.Store{
		BaseModel: entity.BaseModel{ID: "store-1"},
		OwnerID:   "owner-1",
	}
	catRepo.categories["cat-1"] = &entity.EquipmentCategory{
		BaseModel: entity.BaseModel{ID: "cat-1"},
		Name:      "Tenda",
	}

	input := &usecase.CreateEquipmentInput{
		CategoryID:      "cat-1",
		Name:            "Tenda Dome 4 Person",
		Description:     "High quality dome tent",
		TotalStock:      5,
		Condition:       entity.EquipmentConditionNew,
		DepositAmount:   100000,
		RequiresDeposit: true,
		CreatedBy:       "owner-1",
	}

	equipment, err := uc.Create(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if equipment.Name != "Tenda Dome 4 Person" {
		t.Fatalf("expected name Tenda Dome 4 Person, got %s", equipment.Name)
	}
	if equipment.Status != entity.EquipmentStatusAvailable {
		t.Fatalf("expected status available, got %s", equipment.Status)
	}
	if equipment.StoreID != "store-1" {
		t.Fatalf("expected store-1, got %s", equipment.StoreID)
	}
	if len(equipRepo.items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(equipRepo.items))
	}
}

func TestCreateEquipment_NoStore(t *testing.T) {
	uc, _, catRepo, _ := setupEquipmentUsecase()
	ctx := context.Background()

	catRepo.categories["cat-1"] = &entity.EquipmentCategory{
		BaseModel: entity.BaseModel{ID: "cat-1"},
	}

	input := &usecase.CreateEquipmentInput{
		CategoryID: "cat-1",
		Name:       "Test",
		TotalStock: 1,
		CreatedBy:  "user-without-store",
	}

	_, err := uc.Create(ctx, input)
	if err == nil {
		t.Fatal("expected error for user without store")
	}
}

func TestCreateEquipment_InvalidCategory(t *testing.T) {
	uc, _, _, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()

	storeRepo.byOwner["owner-1"] = &entity.Store{
		BaseModel: entity.BaseModel{ID: "store-1"},
		OwnerID:   "owner-1",
	}

	input := &usecase.CreateEquipmentInput{
		CategoryID: "nonexistent-cat",
		Name:       "Test",
		TotalStock: 1,
		CreatedBy:  "owner-1",
	}

	_, err := uc.Create(ctx, input)
	if err == nil {
		t.Fatal("expected error for invalid category")
	}
}

func TestUpdateEquipment_Success(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	newName := "Updated Tent Name"
	newDeposit := float64(150000)
	input := &usecase.UpdateEquipmentInput{
		Name:          &newName,
		DepositAmount: &newDeposit,
	}

	equipment, err := uc.Update(ctx, "equip-1", input, "owner-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if equipment.Name != "Updated Tent Name" {
		t.Fatalf("expected Updated Tent Name, got %s", equipment.Name)
	}
	if equipment.DepositAmount != 150000 {
		t.Fatalf("expected 150000, got %f", equipment.DepositAmount)
	}
}

func TestUpdateEquipment_NotOwner(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	newName := "Hacked"
	input := &usecase.UpdateEquipmentInput{Name: &newName}

	_, err := uc.Update(ctx, "equip-1", input, "other-user")
	if err == nil {
		t.Fatal("expected error for non-owner")
	}
}

func TestUpdateEquipment_NotFound(t *testing.T) {
	uc, _, _, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	storeRepo.byOwner["owner-1"] = &entity.Store{BaseModel: entity.BaseModel{ID: "store-1"}}

	newName := "test"
	input := &usecase.UpdateEquipmentInput{Name: &newName}

	_, err := uc.Update(ctx, "nonexistent", input, "owner-1")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestDeleteEquipment_Success(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	err := uc.Delete(ctx, "equip-1", "owner-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(equipRepo.items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(equipRepo.items))
	}
}

func TestDeleteEquipment_NotOwner(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	err := uc.Delete(ctx, "equip-1", "other-user")
	if err == nil {
		t.Fatal("expected error for non-owner")
	}
}

func TestGetByID_Success(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	equipment, err := uc.GetByID(ctx, "equip-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if equipment.ID != "equip-1" {
		t.Fatalf("expected equip-1, got %s", equipment.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	uc, _, _, _ := setupEquipmentUsecase()
	ctx := context.Background()

	_, err := uc.GetByID(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestChangeStatus_AvailableToMaintenance(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	equipment, err := uc.ChangeStatus(ctx, "equip-1", entity.EquipmentStatusMaintenance, "owner-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if equipment.Status != entity.EquipmentStatusMaintenance {
		t.Fatalf("expected maintenance, got %s", equipment.Status)
	}
}

func TestChangeStatus_AvailableToRetired(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	equipment, err := uc.ChangeStatus(ctx, "equip-1", entity.EquipmentStatusRetired, "owner-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if equipment.Status != entity.EquipmentStatusRetired {
		t.Fatalf("expected retired, got %s", equipment.Status)
	}
	if equipment.IsActive {
		t.Fatal("expected IsActive to be false after retirement")
	}
}

func TestChangeStatus_InvalidTransition(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	// Retire first
	equipRepo.items["equip-1"].Status = entity.EquipmentStatusRetired

	// Try to go back to available from retired (not allowed)
	_, err := uc.ChangeStatus(ctx, "equip-1", entity.EquipmentStatusAvailable, "owner-1")
	if err == nil {
		t.Fatal("expected error for invalid transition from retired")
	}
}

func TestChangeStatus_NotOwner(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	_, err := uc.ChangeStatus(ctx, "equip-1", entity.EquipmentStatusMaintenance, "other-user")
	if err == nil {
		t.Fatal("expected error for non-owner")
	}
}

func TestList_FilterByStatus(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	// Add a maintenance item
	equipRepo.items["equip-2"] = &entity.Equipment{
		BaseModel:  entity.BaseModel{ID: "equip-2"},
		StoreID:    "store-1",
		CategoryID: "cat-1",
		Name:       "Sleeping Bag",
		Status:     entity.EquipmentStatusMaintenance,
		IsActive:   true,
	}

	input := &usecase.ListEquipmentInput{
		Status: "available",
	}
	input.Params = pagination.NewParams(1, 20)

	items, _, err := uc.List(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 available item, got %d", len(items))
	}
}

func TestList_FilterByStore(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	// Add item for different store
	equipRepo.items["equip-other"] = &entity.Equipment{
		BaseModel: entity.BaseModel{ID: "equip-other"},
		StoreID:   "store-other",
		Name:      "Other Store Item",
		Status:    entity.EquipmentStatusAvailable,
		IsActive:  true,
	}

	input := &usecase.ListEquipmentInput{
		StoreID: "store-1",
	}
	input.Params = pagination.NewParams(1, 20)

	items, _, err := uc.List(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item for store-1, got %d", len(items))
	}
}

func TestCheckAvailability_Available(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	input := &usecase.AvailabilityInput{
		EquipmentID: "equip-1",
		StartDate:   "2026-07-01",
		EndDate:     "2026-07-03",
		Quantity:    2,
	}

	result, err := uc.CheckAvailability(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Available {
		t.Fatal("expected available to be true")
	}
	if result.TotalStock != 5 {
		t.Fatalf("expected total stock 5, got %d", result.TotalStock)
	}
	if result.AvailableQty != 5 {
		t.Fatalf("expected available qty 5, got %d", result.AvailableQty)
	}
}

func TestCheckAvailability_RetiredEquipment(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)
	equipRepo.items["equip-1"].Status = entity.EquipmentStatusRetired

	input := &usecase.AvailabilityInput{
		EquipmentID: "equip-1",
		StartDate:   "2026-07-01",
		EndDate:     "2026-07-03",
		Quantity:    1,
	}

	result, err := uc.CheckAvailability(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Available {
		t.Fatal("expected not available for retired equipment")
	}
}

func TestCheckAvailability_InvalidDates(t *testing.T) {
	uc, equipRepo, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()
	seedTestData(equipRepo, catRepo, storeRepo)

	input := &usecase.AvailabilityInput{
		EquipmentID: "equip-1",
		StartDate:   "2026-07-05",
		EndDate:     "2026-07-01", // end before start
		Quantity:    1,
	}

	_, err := uc.CheckAvailability(ctx, input)
	if err == nil {
		t.Fatal("expected error for invalid dates")
	}
}

func TestCreateEquipment_WithPurchaseDate(t *testing.T) {
	uc, _, catRepo, storeRepo := setupEquipmentUsecase()
	ctx := context.Background()

	storeRepo.byOwner["owner-1"] = &entity.Store{
		BaseModel: entity.BaseModel{ID: "store-1"},
		OwnerID:   "owner-1",
	}
	catRepo.categories["cat-1"] = &entity.EquipmentCategory{
		BaseModel: entity.BaseModel{ID: "cat-1"},
	}

	purchaseDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	input := &usecase.CreateEquipmentInput{
		CategoryID:   "cat-1",
		Name:         "Tent with Date",
		TotalStock:   3,
		PurchaseDate: &purchaseDate,
		CreatedBy:    "owner-1",
	}

	equipment, err := uc.Create(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if equipment.PurchaseDate == nil {
		t.Fatal("expected purchase_date to be set")
	}
	if !equipment.PurchaseDate.Equal(purchaseDate) {
		t.Fatalf("expected %v, got %v", purchaseDate, *equipment.PurchaseDate)
	}
}
