// Package main предоставляет API-шлюз для Kinos.
// Обрабатывает HTTP-запросы, управляет аутентификацией и маршрутизирует запросы к gRPC-сервисам.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"kinos/api-service/internal/api/cart"
	"kinos/api-service/internal/api/catalog"
	"kinos/api-service/internal/api/inventory"
	"kinos/api-service/internal/api/middleware"
	"kinos/api-service/internal/api/order"
	"kinos/api-service/internal/api/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	httpPort := getEnv("API_HTTP_PORT", "8080")
	userAddr := getEnv("USER_GRPC_ADDR", "localhost:8081")
	catalogAddr := getEnv("CATALOG_GRPC_ADDR", "localhost:8082")
	inventoryAddr := getEnv("INVENTORY_GRPC_ADDR", "localhost:8083")
	cartAddr := getEnv("CART_GRPC_ADDR", "localhost:8084")
	orderAddr := getEnv("ORDER_GRPC_ADDR", "localhost:8085")
	corsOrigins := splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost,http://localhost:5173,http://localhost:3000"))

	router := gin.Default()

	// Security headers
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Next()
	})

	// Настройка CORS для frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userClient, err := users.NewUserClient(userAddr)
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	catalogClient, err := catalog.NewCatalogClient(catalogAddr)
	if err != nil {
		log.Fatalf("failed to connect to catalog service: %v", err)
	}
	inventoryClient, err := inventory.NewInventoryClient(inventoryAddr)
	if err != nil {
		log.Fatalf("failed to connect to inventory service: %v", err)
	}
	cartClient, err := cart.NewCartClient(cartAddr)
	if err != nil {
		log.Fatalf("failed to connect to cart service: %v", err)
	}
	orderClient, err := order.NewOrderClient(orderAddr)
	if err != nil {
		log.Fatalf("failed to connect to order service: %v", err)
	}
	defer userClient.Close()
	defer catalogClient.Close()
	defer inventoryClient.Close()
	defer cartClient.Close()
	defer orderClient.Close()

	userHandler := users.NewHandler(userClient)
	catalogHandler := catalog.NewHandler(catalogClient)
	inventoryHandler := inventory.NewHandler(inventoryClient)
	cartHandler := cart.NewHandler(cartClient)
	orderHandler := order.NewHandler(orderClient, cartClient)

	authMiddleware := users.NewAuthMiddleware(userClient)
	api := router.Group("/api")

	{
		apiInventory := api.Group("/inventory")
		{
			apiInventory.GET("", inventoryHandler.GetInventory)
			apiInventory.GET("/list", inventoryHandler.GetListInventory)
			apiInventory.POST("/reserve", inventoryHandler.ReserveStock)
			apiInventory.POST("/release", inventoryHandler.ReleaseReservation)
			apiInventory.GET("/warehouses/list", inventoryHandler.GetListWarehouse)
		}

		apiCatalog := api.Group("/catalog")
		{
			apiCatalog.GET("/category", catalogHandler.GetCategory)
			apiCatalog.GET("/categories", catalogHandler.GetListCategory)
			apiCatalog.GET("/product", catalogHandler.GetProduct)
			apiCatalog.GET("/products", catalogHandler.GetProductList)
			apiCatalog.GET("/manufacturer", catalogHandler.GetManufacturers)
			apiCatalog.GET("/manufacturers", catalogHandler.GetManufacturersList)
		}

		// Backward-compatible aliases for older frontend paths.
		api.GET("/category", catalogHandler.GetCategory)
		api.GET("/categories", catalogHandler.GetListCategory)
		api.GET("/product", catalogHandler.GetProduct)
		api.GET("/products", catalogHandler.GetProductList)
		api.GET("/manufacturer", catalogHandler.GetManufacturers)
		api.GET("/manufacturers", catalogHandler.GetManufacturersList)

		// Rate limiter для публичных endpoints (10 запросов в минуту на IP)
		publicRateLimit := middleware.RateLimit(10, time.Minute)

		apiPublic := api.Group("/users")
		apiPublic.Use(publicRateLimit)
		{
			apiPublic.POST("/register", userHandler.Register)
			apiPublic.POST("/login", userHandler.Login)
			apiPublic.POST("/refresh", userHandler.Refresh)
			apiPublic.POST("/revoke", userHandler.Revoke)
		}

		apiLegacyAuth := api.Group("/")
		apiLegacyAuth.Use(publicRateLimit)
		{
			apiLegacyAuth.POST("/register", userHandler.Register)
			apiLegacyAuth.POST("/login", userHandler.Login)
			apiLegacyAuth.POST("/refresh", userHandler.Refresh)
			apiLegacyAuth.POST("/revoke", userHandler.Revoke)
		}

		apiAdmin := api.Group("/admin")
		apiAdmin.Use(authMiddleware.AuthMiddleware(), authMiddleware.AdminOnly())
		{
			apiAdminCatalog := apiAdmin.Group("/catalog")
			{
				apiAdminCatalog.POST("/categories", catalogHandler.CreateCategory)
				apiAdminCatalog.PUT("/categories/:id", catalogHandler.UpdateCategory)
				apiAdminCatalog.DELETE("/categories/:id", catalogHandler.DeleteCategory)

				apiAdminCatalog.POST("/manufacturers", catalogHandler.CreateManufacturer)
				apiAdminCatalog.PUT("/manufacturers/:id", catalogHandler.UpdateManufacturer)
				apiAdminCatalog.DELETE("/manufacturers/:id", catalogHandler.DeleteManufacturer)

				apiAdminCatalog.POST("/products", catalogHandler.CreateProduct)
				apiAdminCatalog.PUT("/products/:id", catalogHandler.UpdateProduct)
				apiAdminCatalog.DELETE("/products/:id", catalogHandler.DeleteProduct)
			}
			apiAdminUsers := apiAdmin.Group("/users")
			{
				apiAdminUsers.PUT("/role", userHandler.UpdateRole)
				apiAdminUsers.DELETE("/:id", userHandler.DeleteUser)
				apiAdminUsers.GET("", userHandler.GetUsers)
			}

			apiAdminInventory := apiAdmin.Group("/inventory")
			{
				apiAdminInventory.POST("", inventoryHandler.CreateInventory)
				apiAdminInventory.PUT("/:id", inventoryHandler.UpdateInventory)
				apiAdminInventory.DELETE("/:id", inventoryHandler.DeleteInventory)
				apiAdminInventory.POST("/warehouses", inventoryHandler.CreateWarehouse)
				apiAdminInventory.PUT("/warehouses/:id", inventoryHandler.UpdateWarehouse)
				apiAdminInventory.DELETE("/warehouses/:id", inventoryHandler.DeleteWarehouse)
			}
		}

		apiAuth := api.Group("/")
		apiAuth.Use(authMiddleware.AuthMiddleware())
		{
			apiAuth.GET("/profile", userHandler.GetProfile)
			apiAuth.PUT("/profile", userHandler.UpdateProfile)
		}

		// Cart routes (защищённые)
		apiCart := api.Group("/cart")
		apiCart.Use(authMiddleware.AuthMiddleware())
		{
			apiCart.GET("", cartHandler.GetCart)
			apiCart.POST("/items", cartHandler.AddItem)
			apiCart.DELETE("/items/:product_id", cartHandler.RemoveItem)
			apiCart.PUT("/items/:product_id", cartHandler.UpdateItem)
			apiCart.POST("/clear", cartHandler.ClearCart)
			apiCart.GET("/count", cartHandler.GetItemsCount)
		}

		// Order routes (защищённые)
		apiOrders := api.Group("/orders")
		apiOrders.Use(authMiddleware.AuthMiddleware())
		{
			apiOrders.POST("", orderHandler.CreateOrder)
			apiOrders.GET("/my", orderHandler.GetUserOrders)
			apiOrders.GET("/:id", orderHandler.GetOrder)
			apiOrders.POST("/:id/cancel", orderHandler.CancelOrder)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Printf("Server started on :%s", httpPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-quit
	log.Println("Shutdown Server ...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return []string{"http://localhost"}
	}
	return result
}
