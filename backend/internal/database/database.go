package database

import (
	"fmt"
	"log"
	"os"

	"ecommerce-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB initializes the PostgreSQL connection and runs migrations
func ConnectDB(databaseURL string) (*gorm.DB, error) {
	dsn := databaseURL
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "ecommerce"
		}
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host, user, password, dbname, port, sslmode)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	log.Println("✅ Berjaya berhubung dengan PostgreSQL!")

	// Auto-migrate all models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Shipment{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate schemas: %w", err)
	}
	log.Println("✅ Auto-migration selesai (users, products, orders, order_items, payments, shipments)!")

	// Seed Initial Admin & Products
	SeedInitialData(DB)

	return DB, nil
}

// SeedInitialData sets up default admin and sample inventory
func SeedInitialData(db *gorm.DB) {
	// 1. Seed Admin User jika belum wujud
	var adminCount int64
	db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err == nil {
			admin := models.User{
				Name:         "Super Administrator",
				Email:        "admin@ecommerce.com",
				PasswordHash: string(hashedPassword),
				Phone:        "0123456789",
				Role:         models.RoleAdmin,
				Address:      "HQ Ecommerce Store, Level 10",
				City:         "Kuala Lumpur",
				State:        "Wilayah Persekutuan Kuala Lumpur",
				Postcode:     "50450",
			}
			if err := db.Create(&admin).Error; err == nil {
				log.Println("👤 Akaun Admin awal dijana: admin@ecommerce.com / admin123")
			}
		}
	}

	// 2. Seed Sample Products jika table kosong
	var productCount int64
	db.Model(&models.Product{}).Count(&productCount)
	if productCount == 0 {
		log.Println("🌱 Menambah produk sampel dengan data EasyParcel weight...")
		products := []models.Product{
			{
				SKU:         "TSH-001",
				Name:        "Classic Cotton T-Shirt",
				Description: "Premium breathable 100% cotton casual t-shirt.",
				Price:       49.90,
				Stock:       50,
				Weight:      0.250, // 250 gram
				ImageURL:    "https://images.unsplash.com/photo-1521572267360-ee0c2909d518?w=500&auto=format&fit=crop&q=80",
				IsActive:    true,
			},
			{
				SKU:         "JNS-002",
				Name:        "Slim Fit Denim Jeans",
				Description: "Durable stretch denim with modern comfortable fit.",
				Price:       119.00,
				Stock:       30,
				Weight:      0.650, // 650 gram
				ImageURL:    "https://images.unsplash.com/photo-1542272604-780c96856592?w=500&auto=format&fit=crop&q=80",
				IsActive:    true,
			},
			{
				SKU:         "WLT-003",
				Name:        "Minimalist Leather Wallet",
				Description: "RFID blocking slim genuine leather bi-fold wallet.",
				Price:       69.50,
				Stock:       40,
				Weight:      0.150, // 150 gram
				ImageURL:    "https://images.unsplash.com/photo-1627123424574-724758594e93?w=500&auto=format&fit=crop&q=80",
				IsActive:    true,
			},
		}

		if err := db.Create(&products).Error; err != nil {
			log.Printf("⚠️ Gagal seed produk: %v\n", err)
		} else {
			log.Println("✅ Produk sampel berjaya dimasukkan!")
		}
	}
}
