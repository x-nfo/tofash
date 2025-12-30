package main

import (
	"log"
	"tofash/internal/config"

	// User Module
	userHandler "tofash/internal/modules/user/handler"
	userRepo "tofash/internal/modules/user/repository"
	userService "tofash/internal/modules/user/service"
	userValidator "tofash/internal/modules/user/utils/validator"

	// Product Module
	productHandler "tofash/internal/modules/product/handlers"
	productRepo "tofash/internal/modules/product/repository"
	productService "tofash/internal/modules/product/service"

	// Order Module
	orderHandler "tofash/internal/modules/order/handlers"
	orderRepo "tofash/internal/modules/order/repository"
	orderService "tofash/internal/modules/order/service"

	// Payment Module
	paymentHandler "tofash/internal/modules/payment/handlers"
	paymentHttpClient "tofash/internal/modules/payment/http_client"
	paymentRepo "tofash/internal/modules/payment/repository"
	paymentService "tofash/internal/modules/payment/service"

	// Notification Module
	notifHandler "tofash/internal/modules/notification/handlers"
	notifMessage "tofash/internal/modules/notification/message"
	notifRepo "tofash/internal/modules/notification/repository"
	notifService "tofash/internal/modules/notification/service"

	"tofash/internal/modules/system/repository"
	"tofash/internal/shared/async"
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

	// Configure CORS with proper settings for credentials
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:4321", "http://localhost:3000"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true, // Important: allows cookies to be sent
		ExposeHeaders:    []string{echo.HeaderContentLength},
		MaxAge:           86400, // 24 hours
	}))

	// Register Custom Validator
	e.Validator = userValidator.NewValidator()

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

	// Redis - REQUIRED for cart functionality
	redisClient := cfg.NewRedisClient()
	if redisClient == nil {
		log.Fatal("[MAIN] Redis connection failed. Redis is REQUIRED for cart functionality. Please ensure Redis is running on localhost:6379")
	}
	log.Println("[MAIN] Redis connected successfully - using Redis cart repository")

	// 5. WIRING: Product Module
	productRepository := productRepo.NewProductRepository(db)
	categoryRepository := productRepo.NewCategoryRepository(db)

	// Use Redis cart repository (Redis is now required)
	cartRepository := productRepo.NewCartRedisRepository(redisClient)
	// productPublisher := productMessage.NewPublishRabbitMQ(cfg) // Removed RabbitMQ

	productSvc := productService.NewProductService(productRepository, nil, categoryRepository)
	categorySvc := productService.NewCategoryService(categoryRepository)
	cartSvc := productService.NewCartService(cartRepository)

	productH := productHandler.NewProductHandler(productSvc)
	categoryH := productHandler.NewCategoryHandler(categorySvc)
	cartH := productHandler.NewCartHandler(cartSvc, productSvc)

	// 6. WIRING: Order Module
	orderRepository := orderRepo.NewOrderRepository(db)
	// RabbitMQ Publisher Removed
	// orderPublisher := orderMessage.NewPublisherRabbitMQ(cfg)

	// Job Queue (System Module)
	jobRepo := repository.NewJobRepository(db)

	orderSvc := orderService.NewOrderService(
		orderRepository,
		cfg,
		jobRepo,
		productSvc,
		userSvc,
	)
	orderH := orderHandler.NewOrderHandler(orderSvc)

	// 7. WIRING: Notification Module
	notifRepository := notifRepo.NewNotificationRepository(db)
	emailSvc := notifMessage.NewMessageEmail(cfg)
	notifSvc := notifService.NewNotificationService(notifRepository, emailSvc)
	notificationH := notifHandler.NewNotificationHandler(notifSvc)
	wsH := notifHandler.NewWebSocketHandler(cfg)

	// consumerRabbit := notifRabbitMQ.NewConsumeRabbitMQ... // Removed

	// 7b. WIRING: Async Worker (Job Queue Consumer)
	jobWorker := async.NewWorker(jobRepo, productSvc, notifSvc)
	go jobWorker.Run()

	// 8. WIRING: Payment Module
	paymentRepository := paymentRepo.NewPaymentRepository(db)
	midtransClient := paymentHttpClient.NewMidtransClient(cfg)
	// paymentPublisher := paymentMessage.NewPublisherRabbitMQ(cfg) // TODO: Refactor payment too if needed
	// For now, pass nil or remove dependency if not critical.
	// Assuming PaymentService still needs refactoring or we ignore for now as per instructions "Refactor Services (Publisher) -> OrderService".
	// But let's check PaymentService signature.

	paymentSvc := paymentService.NewPaymentService(paymentRepository, cfg, midtransClient, orderSvc, userSvc)
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
	auth.GET("/profile", userH.GetProfileUser)

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
