package main

import (
	"tofash/config"
	"tofash/internal/shared/middleware"

	// Import modules
	userH "tofash/internal/modules/user/handler"
	userR "tofash/internal/modules/user/repository"
	userS "tofash/internal/modules/user/service"

	productH "tofash/internal/modules/product/handler"
	productR "tofash/internal/modules/product/repository"
	productS "tofash/internal/modules/product/service"

	// ... import order, payment, notification ...

	"github.com/labstack/echo/v4"
)

func main() {
	// 1. Load Config & Database (Satu koneksi untuk semua)
	cfg := config.LoadConfig()
	db := config.InitDatabase(cfg) // Return *gorm.DB

	// 2. Setup Echo Framework
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 3. WIRING: Setup User Module
	userRepo := userR.NewUserRepository(db)
	userService := userS.NewUserService(userRepo) // Tidak butuh RabbitMQ lagi utk sync data
	userHandler := userH.NewUserHandler(userService)

	// 4. WIRING: Setup Product Module
	// Note: Jika Product butuh data User, kita bisa inject userService ke sini!
	productRepo := productR.NewProductRepository(db)
	productService := productS.NewProductService(productRepo)
	productHandler := productH.NewProductHandler(productService)

	// 5. Setup Routes (Pengganti API Gateway)
	api := e.Group("/api/v1")

	// Routes User
	userGroup := api.Group("/users")
	userGroup.POST("/register", userHandler.Register)
	userGroup.POST("/login", userHandler.Login)

	// Routes Product
	productGroup := api.Group("/products")
	productGroup.GET("", productHandler.GetAll)
	productGroup.POST("", productHandler.Create, middleware.AuthJWT) // Langsung pasang middleware di sini

	// 6. Start Server
	e.Logger.Fatal(e.Start(":8080"))
}
