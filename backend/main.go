package main

import (
	"fmt"
	"log"
	"net/http"

	"ecommerce-backend/internal/api"
	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/database"
)

func main() {
	// 1. Muat konfigurasi .env
	cfg := config.LoadConfig()

	// 2. Sambung ke PostgreSQL & AutoMigrate
	db, err := database.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Printf("⚠️ Amaran Database: %v\n", err)
		log.Println("👉 Sila pastikan DATABASE_URL dalam backend/.env adalah tepat.")
	}

	// 3. Inisialisasi semua laluan API dari package internal/api
	router := api.SetupRouter(cfg, db)

	serverAddr := ":" + cfg.Port
	fmt.Println("==============================================================")
	fmt.Printf("🚀 E-Commerce Backend (Pas Raisy) berjalan di http://localhost%s\n", serverAddr)
	fmt.Println("📦 Modul sedia ada: Auth (Admin/Customer), Product CRUD, Orders (Guest/Member)")
	fmt.Println("💳 Gateway Pembayaran: ToyyibPay")
	fmt.Println("🚚 Pengiraan Pos & Kurier: EasyParcel")
	fmt.Println("📍 Health Check: http://localhost:8080/api/health")
	fmt.Println("📍 Senarai Produk: http://localhost:8080/api/products")
	fmt.Println("==============================================================")

	// 4. Jalankan HTTP server
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Server gagal dimulakan: %v", err)
	}
}
