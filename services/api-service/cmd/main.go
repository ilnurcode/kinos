// Package main предоставляет API-шлюз для Kinos.
// Обрабатывает HTTP-запросы, управляет аутентификацией и маршрутизирует запросы к gRPC-сервисам.
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kinos/api-service/internal/api/catalog"
	"kinos/api-service/internal/api/inventory"
	"kinos/api-service/internal/api/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Настройка CORS для frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userClient := users.NewUserClient("user-service:8081")
	catalogClient := catalog.NewCatalogClient("catalog-service:8082")
	inventoryClient := inventory.NewInventoryClient("inventory-service:8083")

	userHandler := users.NewHandler(userClient)
	catalogHandler := catalog.NewHandler(catalogClient)
	inventoryHandler := inventory.NewHandler(inventoryClient)

	authMiddleware := users.NewAuthMiddleware(userClient)
	api := router.Group("/api")

	{
		apiInventory := api.Group("/inventory")
		{
			apiInventory.GET("", inventoryHandler.GetInventory)
			apiInventory.GET("/list", inventoryHandler.GetListInventory)
			apiInventory.POST("", inventoryHandler.CreateInventory)
			apiInventory.PUT("/:id", inventoryHandler.UpdateInventory)
			apiInventory.DELETE("/:id", inventoryHandler.DeleteInventory)
			apiInventory.POST("/reserve", inventoryHandler.ReserveStock)
			apiInventory.POST("/release", inventoryHandler.ReleaseReservation)

			// Склады
			apiInventory.GET("/warehouses/list", inventoryHandler.GetListWarehouse)
			apiInventory.POST("/warehouses", inventoryHandler.CreateWarehouse)
			apiInventory.PUT("/warehouses/:id", inventoryHandler.UpdateWarehouse)
			apiInventory.DELETE("/warehouses/:id", inventoryHandler.DeleteWarehouse)
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

		apiPublic := api.Group("/users")
		{
			apiPublic.POST("/register", userHandler.Register)
			apiPublic.POST("/login", userHandler.Login)
			apiPublic.POST("/refresh", userHandler.Refresh)
			apiPublic.POST("/revoke", userHandler.Revoke)
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
				apiAdminUsers.GET("", userHandler.GetUsers)
			}
		}

		apiAuth := api.Group("/")
		apiAuth.Use(authMiddleware.AuthMiddleware())
		{
			apiAuth.GET("/profile", userHandler.GetProfile)
			apiAuth.PUT("/profile", userHandler.UpdateProfile)
		}

	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-quit
	log.Println("Shutdown Server ...")
	if err := srv.Close(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
