package service

import (
	"context"
	"errors"
	"tofash/internal/modules/product/entity"
	"tofash/internal/modules/product/message"
	"tofash/internal/modules/product/repository"

	"github.com/labstack/gommon/log"
)

type ProductServiceInterface interface {
	GetAll(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error)
	GetByID(ctx context.Context, productID int64) (*entity.ProductEntity, error)
	Create(ctx context.Context, req entity.ProductEntity) error
	Update(ctx context.Context, req entity.ProductEntity) error
	Delete(ctx context.Context, productID int64) error
	SearchProducts(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error)
	UpdateStock(ctx context.Context, productID int64, quantity int) error
}

type productService struct {
	repo              repository.ProductRepositoryInterface
	publisherRabbitMQ message.PublishRabbitMQInterface
	repoCat           repository.CategoryRepositoryInterface
}

// SearchProducts implements ProductServiceInterface.
func (p *productService) SearchProducts(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
	return p.repo.SearchProducts(ctx, query)
}

// Create implements ProductServiceInterface.
func (p *productService) Create(ctx context.Context, req entity.ProductEntity) error {
	// productID, err := p.repo.Create(ctx, req)
	_, err := p.repo.Create(ctx, req)
	if err != nil {
		log.Errorf("[ProductService-1] Create: %v", err)
		return err
	}

	// getProductByID, err := p.GetByID(ctx, productID)
	// if err != nil {
	// 	log.Errorf("[ProductService-2] Create: %v", err)
	// }

	// RabbitMQ usage to NOT be removed here?
	// The user requested to remove RabbitMQ from Order Service and to create a Worker.
	// But did they ask to remove RabbitMQ from Product Service for syncing?
	// Step 5 says: Remove all RabbitMQ-related code.
	// So I should remove RabbitMQ calls here too.

	// if err := p.publisherRabbitMQ.PublishProductToQueue(*getProductByID); err != nil {
	// 	log.Errorf("[ProductService-3] Create: %v", err)
	// }

	return nil
}

// Delete implements ProductServiceInterface.
func (p *productService) Delete(ctx context.Context, productID int64) error {
	err := p.repo.Delete(ctx, productID)
	if err != nil {
		log.Errorf("[ProductService-1] Delete: %v", err)
		return err
	}

	// RabbitMQ removal
	// if err := p.publisherRabbitMQ.DeleteProductFromQueue(productID); err != nil { ... }

	return nil
}

// GetAll implements ProductServiceInterface.
func (p *productService) GetAll(ctx context.Context, query entity.QueryStringProduct) ([]entity.ProductEntity, int64, int64, error) {
	return p.repo.GetAll(ctx, query)
}

// GetByID implements ProductServiceInterface.
func (p *productService) GetByID(ctx context.Context, productID int64) (*entity.ProductEntity, error) {
	result, err := p.repo.GetByID(ctx, productID)
	if err != nil {
		log.Errorf("[ProductService-1] GetByID: %v", err)
		return nil, err
	}

	resultCat, err := p.repoCat.GetBySlug(ctx, result.CategorySlug)
	if err != nil {
		log.Errorf("[ProductService-2] GetByID: %v", err)
		return nil, err
	}
	if resultCat == nil {
		return nil, errors.New("category not found")
	}
	result.CategoryName = resultCat.Name
	return result, nil
}

// Update implements ProductServiceInterface.
func (p *productService) Update(ctx context.Context, req entity.ProductEntity) error {
	err := p.repo.Update(ctx, req)
	if err != nil {
		log.Errorf("[ProductService-1] Update: %v", err)
		return err
	}

	// RabbitMQ removal
	// if err := p.publisherRabbitMQ.PublishProductToQueue(*getProductByID); err != nil { ... }

	return nil
}

// UpdateStock implements ProductServiceInterface.
func (p *productService) UpdateStock(ctx context.Context, productID int64, quantity int) error {
	product, err := p.repo.GetByID(ctx, productID)
	if err != nil {
		log.Errorf("[ProductService] UpdateStock GetByID: %v", err)
		return err
	}

	if product == nil {
		return errors.New("product not found")
	}

	// Assuming Product.Stock is int or int64?
	// Error was: mismatched types int and int64
	// If product.Stock is int, then cast quantity to int.
	// Let's check entity.

	// Assuming ProductEntity struct has Stock as int based on error.
	if product.Stock < int(quantity) {
		return errors.New("insufficient stock")
	}

	product.Stock -= int(quantity)

	return p.repo.Update(ctx, *product)
}

func NewProductService(repo repository.ProductRepositoryInterface, publisherRabbitMQ message.PublishRabbitMQInterface, repoCat repository.CategoryRepositoryInterface) ProductServiceInterface {
	return &productService{repo: repo, publisherRabbitMQ: publisherRabbitMQ, repoCat: repoCat}
}
