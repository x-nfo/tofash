package handlers

import (
	"encoding/json"
	"net/http"

	// The following line from the instruction appears to have a typo.
	// Assuming it intended to replace 'adapter "tofash/internal/modules/product"'
	// and add '"tofash/internal/config"'.
	// The instruction provided: adapter "tofash/internal/modules/produc	"net/http"
	// This is syntactically incorrect.
	// I will interpret it as:
	// adapter "tofash/internal/modules/product"
	// "tofash/internal/config"
	// If the intent was different, please provide a corrected instruction.

	"tofash/internal/modules/product/entity"
	"tofash/internal/modules/product/handlers/request"
	"tofash/internal/modules/product/handlers/response"
	"tofash/internal/modules/product/service"
	"tofash/internal/modules/product/utils/conv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type CartHandlerInterface interface {
	AddToCart(c echo.Context) error
	GetCart(c echo.Context) error
	RemoveFromCart(c echo.Context) error
	RemoveAllCart(c echo.Context) error
}

type CartHandler struct {
	CartService    service.CartServiceInterface
	ProductService service.ProductServiceInterface
}

func (ch *CartHandler) RemoveAllCart(c echo.Context) error {
	var (
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[CartHandler-1] RemoveFromCart: %s", "data token not found")
		resp.Message = "data token not found"
		resp.Data = nil
		return c.JSON(http.StatusNotFound, resp)
	}
	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		log.Errorf("[CartHandler-2] RemoveFromCart: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}
	userID := jwtUserData.UserID

	err = ch.CartService.RemoveAllCart(ctx, userID)
	if err != nil {
		log.Errorf("[CartHandler-1] RemoveAllCart: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

// AddToCart implements CartHandlerInterface.
func (ch *CartHandler) AddToCart(c echo.Context) error {
	var (
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		request     = request.CartRequest{}
		jwtUserData = entity.JwtUserData{}
	)

	if err := c.Bind(&request); err != nil {
		log.Errorf("[CartHandler-1] AddToCart: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	if err := c.Validate(request); err != nil {
		log.Errorf("[CartHandler-2] AddToCart: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[CartHandler-3] AddToCart: %s", "data token not found")
		resp.Message = "data token not found"
		resp.Data = nil
		return c.JSON(http.StatusNotFound, resp)
	}

	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		log.Errorf("[CartHandler-4] AddToCart: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	userID := jwtUserData.UserID

	reqEntity := entity.CartItem{
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
		Size:      request.Size,
		Color:     request.Color,
		SKU:       request.SKU,
	}

	err = ch.CartService.AddToCart(ctx, userID, reqEntity)
	if err != nil {
		log.Errorf("[CartHandler-5] AddToCart: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusCreated, resp)
}

// GetCart implements CartHandlerInterface.
func (ch *CartHandler) GetCart(c echo.Context) error {
	var (
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		respList    = []response.CartResponse{}
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[CartHandler-1] GetCart: %s", "data token not found")
		resp.Message = "data token not found"
		resp.Data = nil
		return c.JSON(http.StatusNotFound, resp)
	}
	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		log.Errorf("[CartHandler-2] GetCart: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}
	userID := jwtUserData.UserID
	items, err := ch.CartService.GetCartByUserID(ctx, userID)
	if err != nil {
		log.Errorf("[CartHandler-3] GetCart: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}
	for _, item := range items {
		product, err := ch.ProductService.GetByID(ctx, item.ProductID)
		if err != nil {
			log.Errorf("[CartHandler-4] GetCart: %v", err)
			resp.Message = err.Error()
			resp.Data = nil
			return c.JSON(http.StatusInternalServerError, resp)
		}

		respList = append(respList, response.CartResponse{
			ID:            item.ProductID,
			ProductName:   product.Name,
			ProductImage:  product.Image,
			ProductStatus: product.Status,
			SalePrice:     int64(product.SalePrice),
			Quantity:      item.Quantity,
			Unit:          product.Unit,
			Weight:        int64(product.Weight),
			Size:          item.Size,
			Color:         item.Color,
		})
	}

	resp.Message = "success"
	resp.Data = respList
	return c.JSON(http.StatusOK, resp)
}

// RemoveFromCart implements CartHandlerInterface.
func (ch *CartHandler) RemoveFromCart(c echo.Context) error {
	var (
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)

	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[CartHandler-1] RemoveFromCart: %s", "data token not found")
		resp.Message = "data token not found"
		resp.Data = nil
		return c.JSON(http.StatusNotFound, resp)
	}
	err := json.Unmarshal([]byte(user), &jwtUserData)
	if err != nil {
		log.Errorf("[CartHandler-2] RemoveFromCart: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}
	userID := jwtUserData.UserID
	productID := c.QueryParam("product_id")
	if productID == "" {
		log.Errorf("[CartHandler-3] RemoveFromCart: %s", "product_id is required")
		resp.Message = "product_id is required"
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	prodID, err := conv.StringToInt64(productID)

	err = ch.CartService.RemoveFromCart(ctx, userID, prodID)
	if err != nil {
		log.Errorf("[CartHandler-4] RemoveFromCart: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

func NewCartHandler(cartService service.CartServiceInterface, productService service.ProductServiceInterface) CartHandlerInterface {
	return &CartHandler{
		CartService:    cartService,
		ProductService: productService,
	}
}
