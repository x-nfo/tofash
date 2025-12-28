package service

import (
	"context"
	"encoding/json"
	"errors"

	"tofash/internal/config"
	orderService "tofash/internal/modules/order/service"
	"tofash/internal/modules/payment/entity"
	httpclient "tofash/internal/modules/payment/http_client"
	"tofash/internal/modules/payment/repository"
	userService "tofash/internal/modules/user/service"

	"github.com/labstack/gommon/log"
)

type PaymentServiceInterface interface {
	ProcessPayment(ctx context.Context, payment entity.PaymentEntity, accessToken string) (*entity.PaymentEntity, error)
	UpdateStatusByOrderCode(ctx context.Context, orderCode, status string) error
	GetAll(ctx context.Context, req entity.PaymentQueryStringRequest, accessToken string) ([]entity.PaymentEntity, int64, int64, error)
	GetDetail(ctx context.Context, paymentID uint, accessToken string) (*entity.PaymentEntity, error)
}

type paymentService struct {
	repo         repository.PaymentRepositoryInterface
	midtrans     httpclient.MidtransClientInterface
	cfg          *config.Config
	orderService orderService.OrderServiceInterface
	userService  userService.UserServiceInterface
}

// GetDetail implements PaymentServiceInterface.
func (p *paymentService) GetDetail(ctx context.Context, paymentID uint, accessToken string) (*entity.PaymentEntity, error) {
	result, err := p.repo.GetDetail(ctx, paymentID)
	if err != nil {
		log.Errorf("[PaymentService] GetDetail-1: %v", err)
		return nil, err
	}

	var token map[string]interface{}
	err = json.Unmarshal([]byte(accessToken), &token)
	if err != nil {
		log.Errorf("[PaymentService] GetDetail-2: %v", err)
		return nil, err
	}

	userID := int64(result.UserID)
	if token["role_name"].(string) == "Super Admin" {
		userID = 0
	}

	orderDetail, err := p.httpClientOrderService(int64(result.OrderID))
	if err != nil {
		log.Errorf("[PaymentService] GetDetail-3: %v", err)
		return nil, err
	}

	userDetail, err := p.httpClientUserService(userID)
	if err != nil {
		log.Errorf("[PaymentService] GetDetail-4: %v", err)
		return nil, err
	}

	result.CustomerName = userDetail.Name
	result.CustomerEmail = userDetail.Email
	result.CustomerAddress = userDetail.Address

	result.OrderCode = orderDetail.OrderCode
	result.OrderShippingType = orderDetail.ShippingType
	result.OrderAt = orderDetail.OrderDatetime
	result.OrderRemarks = orderDetail.Remarks

	return result, nil
}

// GetAll implements PaymentServiceInterface.
func (p *paymentService) GetAll(ctx context.Context, req entity.PaymentQueryStringRequest, accessToken string) ([]entity.PaymentEntity, int64, int64, error) {
	results, count, total, err := p.repo.GetAll(ctx, req)
	if err != nil {
		log.Errorf("[PaymentService] GetAll-1: %v", err)
		return nil, 0, 0, err
	}

	var token map[string]interface{}
	err = json.Unmarshal([]byte(accessToken), &token)
	if err != nil {
		log.Errorf("[PaymentService] GetAll-2: %v", err)
		return nil, 0, 0, err
	}
	for key, val := range results {
		orderDetail, err := p.httpClientOrderService(int64(val.OrderID))
		if err != nil {
			log.Errorf("[PaymentService] GetAll-3: %v", err)
			return nil, 0, 0, err
		}
		results[key].OrderCode = orderDetail.OrderCode
		results[key].OrderShippingType = orderDetail.ShippingType
	}

	return results, count, total, nil
}

// UpdateStatusByOrderCode implements PaymentServiceInterface.
func (p *paymentService) UpdateStatusByOrderCode(ctx context.Context, orderCode string, status string) error {
	orderDetailID, err := p.httpClientPublicOrderIDByCodeService(orderCode)
	if err != nil {
		log.Errorf("[PaymentService] UpdateStatusByOrderCode-1: %v", err)
		return err
	}

	if err = p.repo.UpdateStatusByOrderCode(ctx, uint(orderDetailID), status); err != nil {
		log.Errorf("[PaymentService] UpdateStatusByOrderCode-2: %v", err)
		return err
	}

	return nil
}

// ProcessPayment implements PaymentServiceInterface.
func (p *paymentService) ProcessPayment(ctx context.Context, payment entity.PaymentEntity, accessToken string) (*entity.PaymentEntity, error) {
	err := p.repo.GetByOrderID(ctx, uint(payment.OrderID))
	if err == nil {
		log.Infof("[PaymentService] ProcessPayment-1: Payment already exists")
		return nil, errors.New("Payment already exists")
	}

	if payment.PaymentMethod == "cod" {
		payment.PaymentStatus = "Success"

		if err := p.repo.CreatePayment(ctx, payment); err != nil {
			log.Errorf("[PaymentService] ProcessPayment-2: %v", err)
			return nil, err
		}

		// REMOVED: RabbitMQ Publish PaymentSuccess
		// if err := p.publisherRabbitMQ.PublishPaymentSuccess(payment); err != nil { ... }

		return &payment, nil
	}

	if payment.PaymentMethod == "midtrans" {
		var token map[string]interface{}
		err := json.Unmarshal([]byte(accessToken), &token)
		if err != nil {
			log.Errorf("[PaymentService] ProcessPayment-4: %v", err)
			return nil, err
		}

		userResponse, err := p.httpClientUserService(int64(payment.UserID))
		if err != nil {
			log.Errorf("[PaymentService] ProcessPayment-5: %v", err)
			return nil, err
		}

		orderDetail, err := p.httpClientOrderService(int64(payment.OrderID))
		if err != nil {
			log.Errorf("[PaymentService] ProcessPayment-6: %v", err)
			return nil, err
		}

		transactionID, err := p.midtrans.CreateTransaction(orderDetail.OrderCode, int64(payment.GrossAmount), userResponse.Name, userResponse.Email)
		if err != nil {
			log.Errorf("[PaymentService] ProcessPayment-7: %v", err)
			return nil, err
		}
		payment.PaymentStatus = "Pending"
		payment.PaymentGatewayID = transactionID

		if err := p.repo.CreatePayment(ctx, payment); err != nil {
			log.Errorf("[PaymentService] ProcessPayment-8: %v", err)
			return nil, err
		}

		// REMOVED: RabbitMQ Publish PaymentSuccess
		// if err := p.publisherRabbitMQ.PublishPaymentSuccess(payment); err != nil { ... }

		return &payment, nil
	}

	return nil, errors.New("Invalid payment method")
}

func (p *paymentService) httpClientOrderService(orderId int64) (*entity.OrderDetailHttpResponse, error) {
	order, err := p.orderService.GetDetailCustomer(context.Background(), orderId)
	if err != nil {
		return nil, err
	}
	return &entity.OrderDetailHttpResponse{
		OrderCode:     order.OrderCode,
		ShippingType:  order.ShippingType,
		OrderDatetime: order.OrderDate,
		Remarks:       order.Remarks,
	}, nil
}

func (p *paymentService) httpClientUserService(userID int64) (*entity.ProfileHttpResponse, error) {
	user, err := p.userService.GetCustomerByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return &entity.ProfileHttpResponse{
		Name:    user.Name,
		Email:   user.Email,
		Address: user.Address,
	}, nil
}

func (p *paymentService) httpClientPublicOrderIDByCodeService(orderCode string) (int64, error) {
	return p.orderService.GetPublicOrderIDByOrderCode(context.Background(), orderCode)
}

func NewPaymentService(repo repository.PaymentRepositoryInterface, cfg *config.Config, midtrans httpclient.MidtransClientInterface, orderService orderService.OrderServiceInterface, userService userService.UserServiceInterface) PaymentServiceInterface {
	return &paymentService{
		repo:         repo,
		midtrans:     midtrans,
		cfg:          cfg,
		orderService: orderService,
		userService:  userService,
	}
}
