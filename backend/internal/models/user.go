package models

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Email        string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"` // Disembunyikan daripada respon JSON
	Phone        string    `gorm:"size:50" json:"phone"`
	Role         Role      `gorm:"size:20;default:'customer'" json:"role"`
	Address      string    `gorm:"type:text" json:"address"`
	City         string    `gorm:"size:100" json:"city"`
	State        string    `gorm:"size:100" json:"state"`
	Postcode     string    `gorm:"size:20" json:"postcode"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
