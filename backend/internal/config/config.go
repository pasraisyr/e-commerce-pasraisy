package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	FrontendURL    string
	BackendURL     string

	// ToyyibPay Config
	ToyyibPaySecretKey    string
	ToyyibPayCategoryCode string
	ToyyibPaySandbox      bool

	// EasyParcel Config
	EasyParcelAPIKey  string
	EasyParcelSandbox bool
}

var AppConfig *Config

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ Tiada fail .env ditemui, membaca konfigurasi sistem.")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-ecommerce-jwt-key-change-in-production"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}

	toyyibSandbox := os.Getenv("TOYYIBPAY_SANDBOX") != "false" // default to sandbox for safety
	easyParcelSandbox := os.Getenv("EASYPARCEL_SANDBOX") != "false"

	AppConfig = &Config{
		Port:                  port,
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		JWTSecret:             jwtSecret,
		FrontendURL:           frontendURL,
		BackendURL:            backendURL,
		ToyyibPaySecretKey:    os.Getenv("TOYYIBPAY_SECRET_KEY"),
		ToyyibPayCategoryCode: os.Getenv("TOYYIBPAY_CATEGORY_CODE"),
		ToyyibPaySandbox:      toyyibSandbox,
		EasyParcelAPIKey:      os.Getenv("EASYPARCEL_API_KEY"),
		EasyParcelSandbox:     easyParcelSandbox,
	}

	return AppConfig
}
