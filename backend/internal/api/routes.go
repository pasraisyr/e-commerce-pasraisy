package api

import (
	"encoding/json"
	"net/http"
	"time"

	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/handlers"
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/repository"
	"ecommerce-backend/internal/service"
	"ecommerce-backend/pkg/easyparcel"
	"ecommerce-backend/pkg/toyyibpay"

	"gorm.io/gorm"
)

// wrap executes a chain of middleware functions on a handler
func wrap(h http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
	var handler http.Handler = h
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}

// SetupRouter initializes all repositories, services, handlers, and registers API endpoints
func SetupRouter(cfg *config.Config, db *gorm.DB) http.Handler {
	// 1. Inisialisasi Repositori
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// 2. Inisialisasi Klien Pihak Ketiga
	toyyibClient := toyyibpay.NewClient(cfg.ToyyibPaySecretKey, cfg.ToyyibPayCategoryCode, cfg.ToyyibPaySandbox)
	easyClient := easyparcel.NewClient(cfg.EasyParcelAPIKey, cfg.EasyParcelSandbox)

	// 3. Inisialisasi Servis Bisnes
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, userRepo, toyyibClient, easyClient, cfg.FrontendURL, cfg.BackendURL)

	// 4. Inisialisasi HTTP Handlers
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	webhookHandler := handlers.NewWebhookHandler(orderService)

	// 5. Middlewares
	corsMiddleware := middleware.EnableCORS(cfg.FrontendURL)
	requireAuth := middleware.RequireAuth(cfg.JWTSecret)
	requireAdmin := middleware.RequireAdmin(cfg.JWTSecret)
	optionalAuth := middleware.OptionalAuth(cfg.JWTSecret)

	mux := http.NewServeMux()

	// --- 1. HEALTH CHECK ---
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dbStatus := "connected"
		if db == nil {
			dbStatus = "disconnected"
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "OK",
			"database":  dbStatus,
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "pas-raisy-ecommerce",
		})
	})

	// --- 2. LALUAN AWAM (PUBLIC) ---
	// Auth
	mux.Handle("POST /api/auth/register", wrap(authHandler.Register))
	mux.Handle("POST /api/auth/login", wrap(authHandler.Login))

	// Produk Katalog
	mux.Handle("GET /api/products", wrap(productHandler.GetProducts))
	mux.Handle("GET /api/products/{id}", wrap(productHandler.GetProductByID))

	// EasyParcel Shipping Rate Check
	mux.Handle("POST /api/shipping/rates", wrap(orderHandler.CheckShippingRates))

	// Checkout Pesanan (Menyokong Guest Checkout dan Logged-in Customer)
	mux.Handle("POST /api/orders", wrap(orderHandler.CreateOrder, optionalAuth))
	mux.Handle("GET /api/orders/track", wrap(orderHandler.GetOrderByNumber))

	// Webhook Pembayaran ToyyibPay
	mux.Handle("POST /api/webhooks/toyyibpay", wrap(webhookHandler.ToyyibPayCallback))

	// --- 3. LALUAN PELANGGAN BERDAFTAR (AUTHENTICATED CUSTOMER) ---
	mux.Handle("GET /api/auth/me", wrap(authHandler.GetMe, requireAuth))
	mux.Handle("GET /api/orders/my-orders", wrap(orderHandler.GetMyOrders, requireAuth))

	// --- 4. LALUAN KHAS ADMIN (ADMIN ONLY) ---
	mux.Handle("POST /api/admin/products", wrap(productHandler.CreateProduct, requireAdmin))
	mux.Handle("PUT /api/admin/products/{id}", wrap(productHandler.UpdateProduct, requireAdmin))
	mux.Handle("DELETE /api/admin/products/{id}", wrap(productHandler.DeleteProduct, requireAdmin))
	mux.Handle("GET /api/admin/orders", wrap(orderHandler.GetAllOrders, requireAdmin))

	// Balut semua endpoint dengan middleware CORS
	return corsMiddleware(mux)
}
