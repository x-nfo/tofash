package service

import (
	"context"
	"errors"
	"testing"

	"tofash/internal/config"
	"tofash/internal/modules/order/entity"
	productEntity "tofash/internal/modules/product/entity"
	userEntity "tofash/internal/modules/user/entity"

	"github.com/stretchr/testify/assert"
)

// ----- Mock implementations -----

type mockOrderRepo struct {
	getAllFn              func(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error)
	getByIDFn             func(ctx context.Context, orderID int64) (*entity.OrderEntity, error)
	createOrderFn         func(ctx context.Context, req entity.OrderEntity) (int64, error)
	updateStatusFn        func(ctx context.Context, req entity.OrderEntity) (int64, string, string, error)
	deleteOrderFn         func(ctx context.Context, orderID int64) error
	getOrderByOrderCodeFn func(ctx context.Context, orderCode string) (*entity.OrderEntity, error)
}

func (m *mockOrderRepo) GetAll(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, queryString)
	}
	return nil, 0, 0, nil
}

func (m *mockOrderRepo) GetByID(ctx context.Context, orderID int64) (*entity.OrderEntity, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, orderID)
	}
	return nil, nil
}

func (m *mockOrderRepo) CreateOrder(ctx context.Context, req entity.OrderEntity) (int64, error) {
	if m.createOrderFn != nil {
		return m.createOrderFn(ctx, req)
	}
	return 0, nil
}

func (m *mockOrderRepo) UpdateStatus(ctx context.Context, req entity.OrderEntity) (int64, string, string, error) {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, req)
	}
	return 0, "", "", nil
}

func (m *mockOrderRepo) DeleteOrder(ctx context.Context, orderID int64) error {
	if m.deleteOrderFn != nil {
		return m.deleteOrderFn(ctx, orderID)
	}
	return nil
}

func (m *mockOrderRepo) GetOrderByOrderCode(ctx context.Context, orderCode string) (*entity.OrderEntity, error) {
	if m.getOrderByOrderCodeFn != nil {
		return m.getOrderByOrderCodeFn(ctx, orderCode)
	}
	return nil, nil
}

type mockElasticRepo struct {
	searchOrderElasticFn          func(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error)
	searchOrderElasticByBuyerIdFn func(ctx context.Context, queryString entity.QueryStringEntity, buyerId int64) ([]entity.OrderEntity, int64, int64, error)
}

func (m *mockElasticRepo) SearchOrderElastic(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error) {
	if m.searchOrderElasticFn != nil {
		return m.searchOrderElasticFn(ctx, queryString)
	}
	return nil, 0, 0, nil
}

func (m *mockElasticRepo) SearchOrderElasticByBuyerId(ctx context.Context, queryString entity.QueryStringEntity, buyerId int64) ([]entity.OrderEntity, int64, int64, error) {
	if m.searchOrderElasticByBuyerIdFn != nil {
		return m.searchOrderElasticByBuyerIdFn(ctx, queryString, buyerId)
	}
	return nil, 0, 0, nil
}

type mockOrderPublisher struct {
	publishUpdateStockFn               func(productID int64, quantity int64)
	publishOrderToQueueFn              func(order entity.OrderEntity) error
	publishSendEmailUpdateStatusFn     func(email, message, queuename string, userID int64) error
	publishDeleteOrderFromQueueFn      func(orderID int64) error
	publishSendPushNotifUpdateStatusFn func(message, queuename string, userID int64) error
	publishUpdateStatusFn              func(queuename string, orderID int64, status string) error
}

func (m *mockOrderPublisher) PublishUpdateStock(productID int64, quantity int64) {
	if m.publishUpdateStockFn != nil {
		m.publishUpdateStockFn(productID, quantity)
	}
}

func (m *mockOrderPublisher) PublishOrderToQueue(order entity.OrderEntity) error {
	if m.publishOrderToQueueFn != nil {
		return m.publishOrderToQueueFn(order)
	}
	return nil
}

func (m *mockOrderPublisher) PublishSendEmailUpdateStatus(email, message, queuename string, userID int64) error {
	if m.publishSendEmailUpdateStatusFn != nil {
		return m.publishSendEmailUpdateStatusFn(email, message, queuename, userID)
	}
	return nil
}

func (m *mockOrderPublisher) PublishDeleteOrderFromQueue(orderID int64) error {
	if m.publishDeleteOrderFromQueueFn != nil {
		return m.publishDeleteOrderFromQueueFn(orderID)
	}
	return nil
}

func (m *mockOrderPublisher) PublishSendPushNotifUpdateStatus(message, queuename string, userID int64) error {
	if m.publishSendPushNotifUpdateStatusFn != nil {
		return m.publishSendPushNotifUpdateStatusFn(message, queuename, userID)
	}
	return nil
}

func (m *mockOrderPublisher) PublishUpdateStatus(queuename string, orderID int64, status string) error {
	if m.publishUpdateStatusFn != nil {
		return m.publishUpdateStatusFn(queuename, orderID, status)
	}
	return nil
}

type mockProductService struct {
	getByIDFn func(ctx context.Context, productID int64) (*productEntity.ProductEntity, error)
}

func (m *mockProductService) GetAll(ctx context.Context, query productEntity.QueryStringProduct) ([]productEntity.ProductEntity, int64, int64, error) {
	return nil, 0, 0, nil
}

func (m *mockProductService) GetByID(ctx context.Context, productID int64) (*productEntity.ProductEntity, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, productID)
	}
	return nil, nil
}

func (m *mockProductService) Create(ctx context.Context, req productEntity.ProductEntity) error {
	return nil
}

func (m *mockProductService) Update(ctx context.Context, req productEntity.ProductEntity) error {
	return nil
}

func (m *mockProductService) Delete(ctx context.Context, productID int64) error {
	return nil
}

func (m *mockProductService) SearchProducts(ctx context.Context, query productEntity.QueryStringProduct) ([]productEntity.ProductEntity, int64, int64, error) {
	return nil, 0, 0, nil
}

type mockUserService struct {
	getByIDFn func(ctx context.Context, userID int64) (*userEntity.UserEntity, error)
}

func (m *mockUserService) SignIn(ctx context.Context, req userEntity.UserEntity) (*userEntity.UserEntity, string, error) {
	return nil, "", nil
}

func (m *mockUserService) CreateUserAccount(ctx context.Context, req userEntity.UserEntity) error {
	return nil
}

func (m *mockUserService) ForgotPassword(ctx context.Context, req userEntity.UserEntity) error {
	return nil
}

func (m *mockUserService) VerifyToken(ctx context.Context, token string) (*userEntity.UserEntity, error) {
	return nil, nil
}

func (m *mockUserService) UpdatePassword(ctx context.Context, req userEntity.UserEntity) error {
	return nil
}

func (m *mockUserService) GetProfileUser(ctx context.Context, userID int64) (*userEntity.UserEntity, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserService) UpdateDataUser(ctx context.Context, req userEntity.UserEntity) error {
	return nil
}

func (m *mockUserService) GetCustomerAll(ctx context.Context, query userEntity.QueryStringCustomer) ([]userEntity.UserEntity, int64, int64, error) {
	return nil, 0, 0, nil
}

func (m *mockUserService) GetCustomerByID(ctx context.Context, customerID int64) (*userEntity.UserEntity, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, customerID)
	}
	return nil, nil
}

func (m *mockUserService) CreateCustomer(ctx context.Context, req userEntity.UserEntity) error {
	return nil
}

func (m *mockUserService) UpdateCustomer(ctx context.Context, req userEntity.UserEntity) error {
	return nil
}

func (m *mockUserService) DeleteCustomer(ctx context.Context, customerID int64) error {
	return nil
}

// ----- Tests -----

func TestOrderService_GetAll_Success(t *testing.T) {
	ctx := context.Background()
	expected := []entity.OrderEntity{{ID: 1, OrderCode: "ORD-001"}}
	mockRepo := &mockOrderRepo{getAllFn: func(_ context.Context, _ entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error) {
		return expected, 1, 1, nil
	}}
	mockElastic := &mockElasticRepo{searchOrderElasticFn: func(_ context.Context, _ entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error) {
		return nil, 0, 0, errors.New("elastic not available")
	}}
	mockUserSvc := &mockUserService{getByIDFn: func(_ context.Context, _ int64) (*userEntity.UserEntity, error) {
		return &userEntity.UserEntity{ID: 10, Name: "User"}, nil
	}}
	mockProductSvc := &mockProductService{getByIDFn: func(_ context.Context, _ int64) (*productEntity.ProductEntity, error) {
		return &productEntity.ProductEntity{ID: 1, Name: "Product"}, nil
	}}
	cfg := &config.Config{}
	svc := NewOrderService(mockRepo, cfg, nil, mockElastic, mockProductSvc, mockUserSvc)
	result, total, page, err := svc.GetAll(ctx, entity.QueryStringEntity{})
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, int64(1), page)
}

func TestOrderService_GetByID_Success(t *testing.T) {
	ctx := context.Background()
	order := &entity.OrderEntity{ID: 1, OrderCode: "ORD-001", BuyerId: 10}
	user := &userEntity.UserEntity{ID: 10, Name: "John Doe"}
	mockRepo := &mockOrderRepo{getByIDFn: func(_ context.Context, _ int64) (*entity.OrderEntity, error) {
		return order, nil
	}}
	mockUserSvc := &mockUserService{getByIDFn: func(_ context.Context, _ int64) (*userEntity.UserEntity, error) {
		return user, nil
	}}
	mockProductSvc := &mockProductService{getByIDFn: func(_ context.Context, _ int64) (*productEntity.ProductEntity, error) {
		return &productEntity.ProductEntity{ID: 1, Name: "Product"}, nil
	}}
	cfg := &config.Config{}
	svc := NewOrderService(mockRepo, cfg, nil, nil, mockProductSvc, mockUserSvc)
	result, err := svc.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", result.BuyerName)
}

func TestOrderService_CreateOrder_Success(t *testing.T) {
	ctx := context.Background()
	createdID := int64(42)
	orderReq := entity.OrderEntity{
		OrderCode:   "ORD-001",
		BuyerId:     10,
		OrderDate:   "2025-12-28",
		OrderTime:   "13:00:00",
		Status:      "Pending",
		TotalAmount: 100000,
		OrderItems: []entity.OrderItemEntity{
			{ProductID: 1, Quantity: 2},
		},
	}
	mockRepo := &mockOrderRepo{
		createOrderFn: func(_ context.Context, _ entity.OrderEntity) (int64, error) {
			return createdID, nil
		},
		getByIDFn: func(_ context.Context, id int64) (*entity.OrderEntity, error) {
			if id == createdID {
				return &orderReq, nil
			}
			return nil, errors.New("not found")
		},
	}
	mockPub := &mockOrderPublisher{publishOrderToQueueFn: func(o entity.OrderEntity) error { return nil }}
	mockElastic := &mockElasticRepo{}
	mockUserSvc := &mockUserService{getByIDFn: func(_ context.Context, _ int64) (*userEntity.UserEntity, error) {
		return &userEntity.UserEntity{ID: 10, Name: "John"}, nil
	}}
	mockProductSvc := &mockProductService{getByIDFn: func(_ context.Context, _ int64) (*productEntity.ProductEntity, error) {
		return &productEntity.ProductEntity{ID: 1, Name: "Product"}, nil
	}}
	cfg := &config.Config{}
	svc := NewOrderService(mockRepo, cfg, mockPub, mockElastic, mockProductSvc, mockUserSvc)
	result, err := svc.CreateOrder(ctx, orderReq)
	assert.NoError(t, err)
	assert.Equal(t, createdID, result)
}

func TestOrderService_UpdateStatus_Success(t *testing.T) {
	ctx := context.Background()
	orderReq := entity.OrderEntity{ID: 5, Status: "Confirmed"}
	mockRepo := &mockOrderRepo{
		updateStatusFn: func(_ context.Context, _ entity.OrderEntity) (int64, string, string, error) {
			return 10, "Confirmed", "ORD-001", nil
		},
	}
	mockPub := &mockOrderPublisher{publishUpdateStatusFn: func(queuename string, orderID int64, status string) error { return nil }}
	mockUserSvc := &mockUserService{getByIDFn: func(_ context.Context, _ int64) (*userEntity.UserEntity, error) {
		return &userEntity.UserEntity{ID: 10, Name: "User"}, nil
	}}
	cfg := &config.Config{}
	svc := NewOrderService(mockRepo, cfg, mockPub, nil, nil, mockUserSvc)
	err := svc.UpdateStatus(ctx, orderReq)
	assert.NoError(t, err)
}

func TestOrderService_DeleteByID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockOrderRepo{deleteOrderFn: func(_ context.Context, _ int64) error { return nil }}
	mockElastic := &mockElasticRepo{}
	mockPub := &mockOrderPublisher{publishDeleteOrderFromQueueFn: func(orderID int64) error { return nil }}
	cfg := &config.Config{}
	svc := NewOrderService(mockRepo, cfg, mockPub, mockElastic, nil, nil)
	err := svc.DeleteByID(ctx, 10)
	assert.NoError(t, err)
}

func TestOrderService_GetOrderByOrderCode_Success(t *testing.T) {
	ctx := context.Background()
	order := &entity.OrderEntity{ID: 1, OrderCode: "ORD-001", BuyerId: 10}
	user := &userEntity.UserEntity{ID: 10, Name: "Jane Doe"}
	mockRepo := &mockOrderRepo{getOrderByOrderCodeFn: func(_ context.Context, _ string) (*entity.OrderEntity, error) {
		return order, nil
	}}
	mockUserSvc := &mockUserService{getByIDFn: func(_ context.Context, _ int64) (*userEntity.UserEntity, error) {
		return user, nil
	}}
	mockProductSvc := &mockProductService{getByIDFn: func(_ context.Context, _ int64) (*productEntity.ProductEntity, error) {
		return &productEntity.ProductEntity{ID: 1, Name: "Product"}, nil
	}}
	cfg := &config.Config{}
	svc := NewOrderService(mockRepo, cfg, nil, nil, mockProductSvc, mockUserSvc)
	result, err := svc.GetOrderByOrderCode(ctx, "ORD-001")
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", result.BuyerName)
}

func TestOrderService_GetAll_ErrorPropagation(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockOrderRepo{getAllFn: func(_ context.Context, _ entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error) {
		return nil, 0, 0, errors.New("db error")
	}}
	cfg := &config.Config{}
	mockElastic := &mockElasticRepo{searchOrderElasticFn: func(_ context.Context, _ entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error) {
		return nil, 0, 0, errors.New("elastic not available")
	}}
	svc := NewOrderService(mockRepo, cfg, nil, mockElastic, nil, nil)
	_, _, _, err := svc.GetAll(ctx, entity.QueryStringEntity{})
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
