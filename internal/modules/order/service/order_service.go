package service

import (
	"context"
	"fmt"
	"tofash/internal/modules/order/config"
	"tofash/internal/modules/order/entity"
	"tofash/internal/modules/order/message"
	"tofash/internal/modules/order/repository"
	"tofash/internal/modules/order/utils"
	"tofash/internal/modules/order/utils/conv"
	productService "tofash/internal/modules/product/service"
	userService "tofash/internal/modules/user/service"

	"github.com/labstack/gommon/log"
)

type OrderServiceInterface interface {
	GetAll(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error)
	GetByID(ctx context.Context, orderID int64) (*entity.OrderEntity, error)
	CreateOrder(ctx context.Context, req entity.OrderEntity) (int64, error)
	UpdateStatus(ctx context.Context, req entity.OrderEntity) error
	GetAllCustomer(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error)
	GetDetailCustomer(ctx context.Context, orderID int64) (*entity.OrderEntity, error)
	DeleteByID(ctx context.Context, orderID int64) error
	GetOrderByOrderCode(ctx context.Context, orderCode string) (*entity.OrderEntity, error)
	GetPublicOrderIDByOrderCode(ctx context.Context, orderCode string) (int64, error)
}

type orderService struct {
	repo              repository.OrderRepositoryInterface
	cfg               *config.Config
	productSvc        productService.ProductServiceInterface
	userSvc           userService.UserServiceInterface
	publisherRabbitMQ message.PublishRabbitMQInterface
	elasticRepo       repository.ElasticRepositoryInterface
}

// GetPublicOrderIDByOrderCode implements OrderServiceInterface.
func (o *orderService) GetPublicOrderIDByOrderCode(ctx context.Context, orderCode string) (int64, error) {
	result, err := o.repo.GetOrderByOrderCode(ctx, orderCode)
	if err != nil {
		log.Errorf("[OrderService-1] GetPublicOrderIDByOrderCode: %v", err)
		return 0, err
	}

	return result.ID, nil
}

// GetOrderByOrderCode implements OrderServiceInterface.
func (o *orderService) GetOrderByOrderCode(ctx context.Context, orderCode string) (*entity.OrderEntity, error) {
	result, err := o.repo.GetOrderByOrderCode(ctx, orderCode)
	if err != nil {
		log.Errorf("[OrderService-1] GetOrderByOrderCode: %v", err)
		return nil, err
	}

	userResponse, err := o.userSvc.GetCustomerByID(ctx, result.BuyerId)
	if err != nil {
		log.Errorf("[OrderService-3] GetOrderByOrderCode: %v", err)
		return nil, err
	}

	result.BuyerName = userResponse.Name
	result.BuyerEmail = userResponse.Email
	result.BuyerPhone = userResponse.Phone
	result.BuyerAddress = userResponse.Address

	for key, val := range result.OrderItems {
		productResponse, err := o.productSvc.GetByID(ctx, val.ProductID)
		if err != nil {
			log.Errorf("[OrderService-4] GetOrderByOrderCode: %v", err)
			return nil, err
		}

		result.OrderItems[key].ProductImage = productResponse.Image
		result.OrderItems[key].ProductName = productResponse.Name
		result.OrderItems[key].Price = int64(productResponse.SalePrice)
	}

	return result, nil
}

// DeleteByID implements OrderServiceInterface.
func (o *orderService) DeleteByID(ctx context.Context, orderID int64) error {
	err := o.repo.DeleteOrder(ctx, orderID)
	if err != nil {
		log.Errorf("[OrderService-1] DeleteByID: %v", err)
		return err
	}

	err = o.publisherRabbitMQ.PublishDeleteOrderFromQueue(orderID)
	if err != nil {
		log.Errorf("[OrderService-2] DeleteByID: %v", err)
		return err
	}

	return nil
}

// GetDetailCustomer implements OrderServiceInterface.
func (o *orderService) GetDetailCustomer(ctx context.Context, orderID int64) (*entity.OrderEntity, error) {
	result, err := o.repo.GetByID(ctx, orderID)
	if err != nil {
		log.Errorf("[OrderService-1] GetByID: %v", err)
		return nil, err
	}

	userResponse, err := o.userSvc.GetCustomerByID(ctx, result.BuyerId)
	if err != nil {
		log.Errorf("[OrderService-3] GetByID: %v", err)
		return nil, err
	}

	result.BuyerName = userResponse.Name
	result.BuyerEmail = userResponse.Email
	result.BuyerPhone = userResponse.Phone
	result.BuyerAddress = userResponse.Address

	for key, val := range result.OrderItems {
		productResponse, err := o.productSvc.GetByID(ctx, val.ProductID)
		if err != nil {
			log.Errorf("[OrderService-3] GetByID: %v", err)
			return nil, err
		}

		result.OrderItems[key].ProductImage = productResponse.Image
		if len(productResponse.Child) > 0 {
			result.OrderItems[key].ProductImage = productResponse.Child[0].Image
		}
		result.OrderItems[key].ProductName = productResponse.Name
		result.OrderItems[key].Price = int64(productResponse.SalePrice)
		result.OrderItems[key].ProductWeight = int64(productResponse.Weight)
		result.OrderItems[key].ProductUnit = productResponse.Unit
	}

	return result, nil
}

// GetAllCustomer implements OrderServiceInterface.
func (o *orderService) GetAllCustomer(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error) {
	results, count, total, err := o.elasticRepo.SearchOrderElasticByBuyerId(ctx, queryString, queryString.BuyerID)
	if err == nil {
		return results, count, total, nil
	} else {
		log.Errorf("[OrderService-1] GetAllCustomer: %v", err)
	}

	results, count, total, err = o.repo.GetAll(ctx, queryString)
	if err != nil {
		log.Errorf("[OrderService-2] GetAllCustomer: %v", err)
		return nil, 0, 0, err
	}

	for key, val := range results {
		userResponse, err := o.userSvc.GetCustomerByID(ctx, val.BuyerId)
		if err != nil {
			log.Errorf("[OrderService-4] GetAllCustomer: %v", err)
			return nil, 0, 0, err
		}

		results[key].BuyerName = userResponse.Name
		results[key].BuyerEmail = userResponse.Email
		results[key].BuyerPhone = userResponse.Phone
		results[key].BuyerAddress = userResponse.Address

		for key2, res := range val.OrderItems {

			productResponse, err := o.productSvc.GetByID(ctx, res.ProductID)
			if err != nil {
				log.Errorf("[OrderService-5] GetAllCustomer: %v", err)
				return nil, 0, 0, err
			}

			val.OrderItems[key2].ProductImage = productResponse.Image
			val.OrderItems[key2].ProductName = productResponse.Name
			val.OrderItems[key2].Price = int64(productResponse.SalePrice)
			val.OrderItems[key2].Quantity = res.Quantity
			val.OrderItems[key2].ProductUnit = productResponse.Unit
			val.OrderItems[key2].ProductWeight = int64(productResponse.Weight)
		}
	}

	return results, count, total, nil
}

// UpdateStatus implements OrderServiceInterface.
func (o *orderService) UpdateStatus(ctx context.Context, req entity.OrderEntity) error {
	buyerID, statusOrder, orderCode, err := o.repo.UpdateStatus(ctx, req)
	if err != nil {
		log.Errorf("[OrderService-1] UpdateStatus: %v", err)
		return err
	}

	userResponse, err := o.userSvc.GetCustomerByID(ctx, buyerID)
	if err != nil {
		log.Errorf("[OrderService-3] UpdateStatus: %v", err)
		return err
	}
	message := fmt.Sprintf("Hello,\n\nYour order with ID %s has been updated to status: %s.\n\nThank you for shopping with us!", orderCode, statusOrder)
	go o.publisherRabbitMQ.PublishSendEmailUpdateStatus(userResponse.Email, message, o.cfg.PublisherName.EmailUpdateStatus, buyerID)
	go o.publisherRabbitMQ.PublishSendPushNotifUpdateStatus(message, utils.PUSH_NOTIF, buyerID)
	go o.publisherRabbitMQ.PublishUpdateStatus(o.cfg.PublisherName.PublisherUpdateStatus, req.ID, req.Status)

	return nil
}

// CreateOrder implements OrderServiceInterface.
func (o *orderService) CreateOrder(ctx context.Context, req entity.OrderEntity) (int64, error) {
	req.OrderCode = conv.GenerateOrderCode()
	shippingFee := 0
	if req.ShippingType == "Delivery" {
		shippingFee = 5000
	}
	req.ShippingFee = int64(shippingFee)
	req.Status = "Pending"
	orderID, err := o.repo.CreateOrder(ctx, req)
	if err != nil {
		log.Errorf("[OrderService-1] CreateOrder: %v", err)
		return 0, err
	}

	resultData, err := o.GetByID(ctx, orderID)
	if err != nil {
		log.Errorf("[OrderService-2] CreateOrder: %v", err)
	}

	if err := o.publisherRabbitMQ.PublishOrderToQueue(*resultData); err != nil {
		log.Errorf("[OrderService-3] CreateOrder: %v", err)
	}

	for _, orderItem := range req.OrderItems {
		o.publisherRabbitMQ.PublishUpdateStock(orderItem.ProductID, orderItem.Quantity)
	}

	return orderID, nil
}

// GetByID implements OrderServiceInterface.
func (o *orderService) GetByID(ctx context.Context, orderID int64) (*entity.OrderEntity, error) {
	result, err := o.repo.GetByID(ctx, orderID)
	if err != nil {
		log.Errorf("[OrderService-1] GetByID: %v", err)
		return nil, err
	}

	userResponse, err := o.userSvc.GetCustomerByID(ctx, result.BuyerId)
	if err != nil {
		log.Errorf("[OrderService-2] GetByID: %v", err)
		return nil, err
	}

	result.BuyerName = userResponse.Name
	result.BuyerEmail = userResponse.Email
	result.BuyerPhone = userResponse.Phone
	result.BuyerAddress = userResponse.Address

	for key, val := range result.OrderItems {
		productResponse, err := o.productSvc.GetByID(ctx, val.ProductID)
		if err != nil {
			log.Errorf("[OrderService-3] GetByID: %v", err)
			return nil, err
		}

		result.OrderItems[key].ProductImage = productResponse.Image
		result.OrderItems[key].ProductName = productResponse.Name
		result.OrderItems[key].Price = int64(productResponse.SalePrice)
	}

	return result, nil
}

// GetAll implements OrderServiceInterface.
func (o *orderService) GetAll(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error) {
	results, count, total, err := o.elasticRepo.SearchOrderElastic(ctx, queryString)
	if err == nil {
		return results, count, total, nil
	} else {
		log.Errorf("[OrderService-1] GetAll: %v", err)
	}

	results, count, total, err = o.repo.GetAll(ctx, queryString)
	if err != nil {
		log.Errorf("[OrderService-2] GetAll: %v", err)
		return nil, 0, 0, err
	}

	for key, val := range results {

		userResponse, err := o.userSvc.GetCustomerByID(ctx, val.BuyerId)
		if err != nil {
			log.Errorf("[OrderService-4] GetAll: %v", err)
			return nil, 0, 0, err
		}
		results[key].BuyerName = userResponse.Name

		for key2, res := range val.OrderItems {

			productResponse, err := o.productSvc.GetByID(ctx, res.ProductID)
			if err != nil {
				log.Errorf("[OrderService-5] GetAll: %v", err)
				return nil, 0, 0, err
			}

			val.OrderItems[key2].ProductImage = productResponse.Image
		}
	}

	return results, count, total, nil
}

func NewOrderService(repo repository.OrderRepositoryInterface, cfg *config.Config, publisherRabbitMQ message.PublishRabbitMQInterface, elasticRepo repository.ElasticRepositoryInterface, productSvc productService.ProductServiceInterface, userSvc userService.UserServiceInterface) OrderServiceInterface {
	return &orderService{
		repo:              repo,
		cfg:               cfg,
		publisherRabbitMQ: publisherRabbitMQ,
		elasticRepo:       elasticRepo,
		productSvc:        productSvc,
		userSvc:           userSvc,
	}
}
