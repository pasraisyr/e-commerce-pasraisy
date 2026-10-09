package models

import "time"

// Product represents an item in the ecommerce store
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SKU         string    `gorm:"size:100;uniqueIndex" json:"sku"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int       `gorm:"default:0" json:"stock"`
	Weight      float64   `gorm:"type:decimal(8,3);default:0.5" json:"weight"` // Berat dalam KG (penting untuk kiraan EasyParcel)
	ImageURL    string    `gorm:"size:500" json:"imageUrl"`
	IsActive    bool      `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
