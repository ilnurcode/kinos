// Package main предоставляет API-шлюз для Kinos.
// Обрабатывает HTTP-запросы, управляет аутентификацией и маршрутизирует запросы к gRPC-сервисам.
package main

import (
	"context"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	assets "kinos/api-service"
	"kinos/api-service/internal/api/catalog"
	"kinos/api-service/internal/api/inventory"
	"kinos/api-service/internal/api/users"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Явно указываем все шаблоны для правильного сохранения имён
	tmpl := template.Must(template.ParseFS(assets.FS,
		"templates/index.html",
		"templates/login.html",
		"templates/registration.html",
		"templates/profile.html",
		"templates/edit_profile.html",
		"templates/catalog.html",
		"templates/admin/users.html",
		"templates/admin/categories.html",
		"templates/admin/manufacturers.html",
		"templates/admin/products.html",
		"templates/admin/inventory.html",
		"templates/admin/warehouses.html",
	))
	router.SetHTMLTemplate(tmpl)

	staticFS, err := fs.Sub(assets.FS, "static")
	if err != nil {
		log.Fatal(err)
	}
	router.StaticFS("/static", http.FS(staticFS))

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

	router.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", gin.H{"title": "Welcome to Kinos!"}) })
	router.GET("/catalog", func(c *gin.Context) {
		c.HTML(http.StatusOK, "catalog.html", gin.H{"title": "Каталог товаров"})
	})
	router.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "registration.html", gin.H{"title": "Registration"}) })
	router.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", gin.H{"title": "Login"}) })

	// HTML страницы админ-панели — без middleware, проверка через JS
	router.GET("/admin/users", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users.html", gin.H{"title": "Управление пользователями"})
	})
	router.GET("/admin/categories", func(c *gin.Context) {
		c.HTML(http.StatusOK, "categories.html", gin.H{"title": "Управление категориями"})
	})
	router.GET("/admin/manufacturers", func(c *gin.Context) {
		c.HTML(http.StatusOK, "manufacturers.html", gin.H{"title": "Управление производителями"})
	})
	router.GET("/admin/products", func(c *gin.Context) {
		c.HTML(http.StatusOK, "products.html", gin.H{"title": "Управление товарами"})
	})
	router.GET("/admin/inventory", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inventory.html", gin.H{"title": "Управление запасами"})
	})
	router.GET("/admin/warehouses", func(c *gin.Context) {
		c.HTML(http.StatusOK, "warehouses.html", gin.H{"title": "Управление складами"})
	})

	// HTML страницы профиля — без middleware, аутентификация через JS
	router.GET("/profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", gin.H{"title": "Профиль"})
	})
	router.GET("/profile/edit", func(c *gin.Context) {
		c.HTML(http.StatusOK, "edit_profile.html", gin.H{"title": "Редактирование профиля"})
	})

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
