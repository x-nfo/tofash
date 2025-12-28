package main

import (
	"log"
	"tofash/internal/config"

	// User Module
	userHandler "tofash/internal/modules/user/handler"
	userRepo "tofash/internal/modules/user/repository"
	userService "tofash/internal/modules/user/service"

	// Product Module
	productHandler "tofash/internal/modules/product/handlers"
	productMessage "tofash/internal/modules/product/message"
	productRepo "tofash/internal/modules/product/repository"
	productService "tofash/internal/modules/product/service"

	// Order Module
	orderHandler "tofash/internal/modules/order/handlers"
	orderMessage "tofash/internal/modules/order/message"
	orderRepo "tofash/internal/modules/order/repository"
	orderService "tofash/internal/modules/order/service"

	// Payment Module
	paymentHandler "tofash/internal/modules/payment/handlers"
	paymentHttpClient "tofash/internal/modules/payment/http_client"
	paymentMessage "tofash/internal/modules/payment/message"
	paymentRepo "tofash/internal/modules/payment/repository"
	paymentService "tofash/internal/modules/payment/service"

	// Notification Module
	notifHandler "tofash/internal/modules/notification/handlers"
	notifMessage "tofash/internal/modules/notification/message"
	notifRabbitMQ "tofash/internal/modules/notification/rabbitmq"
	notifRepo "tofash/internal/modules/notification/repository"
	notifService "tofash/internal/modules/notification/service"

	// Shared
	mid "tofash/internal/shared/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Load Config & Database
	cfg := config.LoadConfig()
	db := config.InitDatabase(cfg)
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	// 2. Setup Echo
	e := echo.New()
	e.Use(middleware.Recover())

	// 3. Setup Dependencies

	// 4. WIRING: User Module
	userRepository := userRepo.NewUserRepository(db)
	roleRepository := userRepo.NewRoleRepository(db)
	verificationTokenRepository := userRepo.NewVerificationTokenRepository(db)

	jwtSvc := userService.NewJwtService(cfg)

	userSvc := userService.NewUserService(
		userRepository,
		cfg,
		jwtSvc,
		verificationTokenRepository,
	)
	roleSvc := userService.NewRoleService(roleRepository)

	userH := userHandler.NewUserHandler(userSvc)
	roleH := userHandler.NewRoleHandler(roleSvc)

	// Middleware
	authMiddleware := mid.NewAuthMiddleware(cfg, jwtSvc)

	// Redis
	redisClient := cfg.NewRedisClient()

	// 5. WIRING: Product Module
	productRepository := productRepo.NewProductRepository(db)
	categoryRepository := productRepo.NewCategoryRepository(db)
	cartRepository := productRepo.NewCartRedisRepository(redisClient)
	productPublisher := productMessage.NewPublishRabbitMQ(cfg)

	productSvc := productService.NewProductService(productRepository, productPublisher, categoryRepository)
	categorySvc := productService.NewCategoryService(categoryRepository)
	cartSvc := productService.NewCartService(cartRepository)

	productH := productHandler.NewProductHandler(productSvc)
	categoryH := productHandler.NewCategoryHandler(categorySvc)
	cartH := productHandler.NewCartHandler(cartSvc, productSvc)

	// 6. WIRING: Order Module
	orderRepository := orderRepo.NewOrderRepository(db)
	// RabbitMQ Publisher
	orderPublisher := orderMessage.NewPublisherRabbitMQ(cfg)

	orderSvc := orderService.NewOrderService(
		orderRepository,
		cfg,
		orderPublisher,
		productSvc,
		userSvc,
	)
	orderH := orderHandler.NewOrderHandler(orderSvc)

	// 7. WIRING: Notification Module
	notifRepository := notifRepo.NewNotificationRepository(db)
	emailSvc := notifMessage.NewMessageEmail(cfg)
	notifSvc := notifService.NewNotificationService(notifRepository)
	notificationH := notifHandler.NewNotificationHandler(notifSvc)
	wsH := notifHandler.NewWebSocketHandler(cfg)

	consumeRabbit := notifRabbitMQ.NewConsumeRabbitMQ(emailSvc, notifRepository, notifSvc)
	go consumeRabbit.ConsumeMessage("notification_queue") // Use config queue name ideally

	// 8. WIRING: Payment Module
	paymentRepository := paymentRepo.NewPaymentRepository(db)
	midtransClient := paymentHttpClient.NewMidtransClient(cfg)
	paymentPublisher := paymentMessage.NewPublisherRabbitMQ(cfg)
	// Direct service injection
	paymentSvc := paymentService.NewPaymentService(paymentRepository, cfg, midtransClient, paymentPublisher, orderSvc, userSvc)
	paymentH := paymentHandler.NewPaymentHandler(paymentSvc)

	// 7. Setup Routes
	api := e.Group("/api/v1")

	// Public Routes
	api.POST("/register", userH.CreateUserAccount)
	api.POST("/login", userH.SignIn)
	api.POST("/forgot-password", userH.ForgotPassword)

	// Product
	api.GET("/products", productH.GetAllShop)
	api.GET("/products/home", productH.GetAllHome)
	api.GET("/products/:id", productH.GetDetailHome)
	api.GET("/categories", categoryH.GetAllShop) // Use Shop or Home variant

	// Secured Routes
	auth := api.Group("", authMiddleware.CheckToken)

	// User
	auth.POST("/roles", roleH.Create)
	auth.GET("/verify-account", userH.VerifyAccount)

	// Cart
	auth.POST("/carts", cartH.AddToCart)
	auth.GET("/carts", cartH.GetCart)
	auth.DELETE("/carts", cartH.RemoveFromCart)
	auth.DELETE("/carts/all", cartH.RemoveAllCart)

	// Order
	auth.POST("/orders", orderH.CreateOrder) // Consider DistanceCheck middleware if needed
	auth.GET("/orders", orderH.GetAllCustomer)
	auth.GET("/orders/:orderID", orderH.GetDetailCustomer)

	// Payment
	auth.POST("/payments", paymentH.Create)
	auth.GET("/payments", paymentH.GetAllCustomer)
	auth.GET("/payments/:id", paymentH.GetDetail)

	// Notification
	auth.GET("/notifications", notificationH.GetAll)
	auth.GET("/notifications/:id", notificationH.GetByID)
	auth.PUT("/notifications/:id", notificationH.MarkAsRead)

	// Admin Routes (reusing auth middleware but logic checks role)
	admin := api.Group("/admin", authMiddleware.CheckToken)
	admin.GET("/payments", paymentH.GetAllAdmin)

	// Webhooks & Public
	api.POST("/midtrans/webhook", paymentH.MidtranswebHookHandler)
	e.GET("/ws", wsH.WebSocketHandler)

	// Start Server
	e.Logger.Fatal(e.Start(":" + cfg.App.AppPort))
}
