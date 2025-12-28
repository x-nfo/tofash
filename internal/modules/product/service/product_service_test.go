package service

import (
	"context"
	"errors"
	"testing"

	"tofash/internal/modules/product/entity"

	"github.com/stretchr/testify/assert"
)

// ----- Mock implementations -----

type mockProductRepo struct {
	getAllFn  func(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error)
	getByIDFn func(ctx context.Context, id int64) (*entity.ProductEntity, error)
	createFn  func(ctx context.Context, req entity.ProductEntity) (int64, error)
	updateFn  func(ctx context.Context, req entity.ProductEntity) error
	deleteFn  func(ctx context.Context, id int64) error
	searchFn  func(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error)
}

func (m *mockProductRepo) GetAll(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
	return m.getAllFn(ctx, query)
}
func (m *mockProductRepo) GetByID(ctx context.Context, productID int64) (*entity.ProductEntity, error) {
	return m.getByIDFn(ctx, productID)
}
func (m *mockProductRepo) Create(ctx context.Context, req entity.ProductEntity) (int64, error) {
	return m.createFn(ctx, req)
}
func (m *mockProductRepo) Update(ctx context.Context, req entity.ProductEntity) error {
	return m.updateFn(ctx, req)
}
func (m *mockProductRepo) Delete(ctx context.Context, productID int64) error {
	return m.deleteFn(ctx, productID)
}
func (m *mockProductRepo) SearchProducts(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
	return m.searchFn(ctx, query)
}

type mockCategoryRepo struct {
	getAllFn          func(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.CategoryEntity, int64, int64, error)
	getByIDFn         func(ctx context.Context, categoryID int64) (*entity.CategoryEntity, error)
	getBySlugFn       func(ctx context.Context, slug string) (*entity.CategoryEntity, error)
	createFn          func(ctx context.Context, req entity.CategoryEntity) error
	editFn            func(ctx context.Context, req entity.CategoryEntity) error
	deleteFn          func(ctx context.Context, categoryID int64) error
	getAllPublishedFn func(ctx context.Context) ([]entity.CategoryEntity, error)
}

func (m *mockCategoryRepo) GetAll(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.CategoryEntity, int64, int64, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, queryString)
	}
	return nil, 0, 0, nil
}

func (m *mockCategoryRepo) GetByID(ctx context.Context, categoryID int64) (*entity.CategoryEntity, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, categoryID)
	}
	return nil, nil
}

func (m *mockCategoryRepo) GetBySlug(ctx context.Context, slug string) (*entity.CategoryEntity, error) {
	if m.getBySlugFn != nil {
		return m.getBySlugFn(ctx, slug)
	}
	return nil, nil
}

func (m *mockCategoryRepo) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return nil
}

func (m *mockCategoryRepo) EditCategory(ctx context.Context, req entity.CategoryEntity) error {
	if m.editFn != nil {
		return m.editFn(ctx, req)
	}
	return nil
}

func (m *mockCategoryRepo) DeleteCategory(ctx context.Context, categoryID int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, categoryID)
	}
	return nil
}

func (m *mockCategoryRepo) GetAllPublished(ctx context.Context) ([]entity.CategoryEntity, error) {
	if m.getAllPublishedFn != nil {
		return m.getAllPublishedFn(ctx)
	}
	return nil, nil
}

type mockPublisher struct {
	publishFn func(product entity.ProductEntity) error
	deleteFn  func(productID int64) error
}

func (m *mockPublisher) PublishProductToQueue(product entity.ProductEntity) error {
	return m.publishFn(product)
}
func (m *mockPublisher) DeleteProductFromQueue(productID int64) error {
	return m.deleteFn(productID)
}

// ----- Tests -----

func TestProductService_GetAll_Success(t *testing.T) {
	ctx := context.Background()
	expected := []entity.ProductEntity{{ID: 1, Name: "Test"}}
	mockRepo := &mockProductRepo{getAllFn: func(_ context.Context, _ entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
		return expected, 1, 1, nil
	}}
	svc := NewProductService(mockRepo, nil, nil)
	result, total, page, err := svc.GetAll(ctx, entity.QueryStringProduct{})
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, int64(1), page)
	_ = result
	_ = total
	_ = page
	_ = err
}

func TestProductService_GetByID_Success(t *testing.T) {
	ctx := context.Background()
	prod := &entity.ProductEntity{ID: 1, CategorySlug: "cat-1"}
	cat := &entity.CategoryEntity{Name: "Category 1"}
	mockRepo := &mockProductRepo{getByIDFn: func(_ context.Context, _ int64) (*entity.ProductEntity, error) {
		return prod, nil
	}}
	mockCatRepo := &mockCategoryRepo{getBySlugFn: func(_ context.Context, _ string) (*entity.CategoryEntity, error) {
		return cat, nil
	}}
	svc := NewProductService(mockRepo, nil, mockCatRepo)
	result, err := svc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Category 1", result.CategoryName)
}

func TestProductService_Create_Success(t *testing.T) {
	ctx := context.Background()
	createdID := int64(42)
	prodReq := entity.ProductEntity{ID: createdID, Name: "New", CategorySlug: "cat-1"}
	mockRepo := &mockProductRepo{
		createFn: func(_ context.Context, _ entity.ProductEntity) (int64, error) {
			return createdID, nil
		},
		getByIDFn: func(_ context.Context, id int64) (*entity.ProductEntity, error) {
			if id == createdID {
				return &prodReq, nil
			}
			return nil, errors.New("not found")
		},
	}
	mockCatRepo := &mockCategoryRepo{getBySlugFn: func(_ context.Context, _ string) (*entity.CategoryEntity, error) {
		return &entity.CategoryEntity{Name: "Category 1"}, nil
	}}
	mockPub := &mockPublisher{publishFn: func(p entity.ProductEntity) error { return nil }}
	svc := NewProductService(mockRepo, mockPub, mockCatRepo)
	err := svc.Create(ctx, prodReq)
	assert.NoError(t, err)
}

func TestProductService_Update_Success(t *testing.T) {
	ctx := context.Background()
	prod := entity.ProductEntity{ID: 5, Name: "Old", CategorySlug: "cat-1"}
	mockRepo := &mockProductRepo{
		updateFn: func(_ context.Context, _ entity.ProductEntity) error { return nil },
		getByIDFn: func(_ context.Context, id int64) (*entity.ProductEntity, error) {
			if id == prod.ID {
				return &prod, nil
			}
			return nil, errors.New("not found")
		},
	}
	mockCatRepo := &mockCategoryRepo{getBySlugFn: func(_ context.Context, _ string) (*entity.CategoryEntity, error) {
		return &entity.CategoryEntity{Name: "Category 1"}, nil
	}}
	mockPub := &mockPublisher{publishFn: func(p entity.ProductEntity) error { return nil }}
	svc := NewProductService(mockRepo, mockPub, mockCatRepo)
	err := svc.Update(ctx, prod)
	assert.NoError(t, err)
}

func TestProductService_Delete_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockProductRepo{deleteFn: func(_ context.Context, _ int64) error { return nil }}
	mockPub := &mockPublisher{deleteFn: func(id int64) error { return nil }}
	svc := NewProductService(mockRepo, mockPub, nil)
	err := svc.Delete(ctx, 10)
	assert.NoError(t, err)
}

func TestProductService_GetAll_ErrorPropagation(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockProductRepo{getAllFn: func(_ context.Context, _ entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
		return nil, 0, 0, errors.New("db error")
	}}
	svc := NewProductService(mockRepo, nil, nil)
	_, _, _, err := svc.GetAll(ctx, entity.QueryStringProduct{})
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

// Additional tests for SearchProducts and error cases can be added similarly.
