package models

import "time"

type OrderStatus string

const (
	StatusPendingPayment OrderStatus = "pending_payment"
	StatusPaid           OrderStatus = "paid"
	StatusProcessing     OrderStatus = "processing"
	StatusShipped        OrderStatus = "shipped"
	StatusCompleted      OrderStatus = "completed"
	StatusCancelled      OrderStatus = "cancelled"
)

type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	OrderNumber string      `gorm:"size:50;uniqueIndex;not null" json:"orderNumber"`
	UserID      *uint       `gorm:"index" json:"userId"` // Boleh NULL bagi Guest Checkout!
	User        *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// Maklumat Pelanggan (Wajib bagi Guest mahupun Pelanggan Berdaftar)
	CustomerName  string `gorm:"size:255;not null" json:"customerName"`
	CustomerEmail string `gorm:"size:255;not null" json:"customerEmail"`
	CustomerPhone string `gorm:"size:50;not null" json:"customerPhone"`

	// Alamat Penghantaran
	ShippingAddress string `gorm:"type:text;not null" json:"shippingAddress"`
	City            string `gorm:"size:100;not null" json:"city"`
	State           string `gorm:"size:100;not null" json:"state"`
	Postcode        string `gorm:"size:20;not null" json:"postcode"`

	// Kewangan & Berat
	Subtotal    float64 `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	ShippingFee float64 `gorm:"type:decimal(10,2);default:0" json:"shippingFee"`
	TotalAmount float64 `gorm:"type:decimal(10,2);not null" json:"totalAmount"`
	TotalWeight float64 `gorm:"type:decimal(8,3);default:0" json:"totalWeight"` // Jumlah KG pesanan

	Status OrderStatus `gorm:"size:30;default:'pending_payment'" json:"status"`
	Notes  string      `gorm:"type:text" json:"notes"`

	Items    []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Payment  *Payment    `gorm:"foreignKey:OrderID" json:"payment,omitempty"`
	Shipment *Shipment   `gorm:"foreignKey:OrderID" json:"shipment,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type OrderItem struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	OrderID     uint    `gorm:"index;not null" json:"orderId"`
	ProductID   uint    `gorm:"not null" json:"productId"`
	ProductName string  `gorm:"size:255;not null" json:"productName"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity    int     `gorm:"not null" json:"quantity"`
	Subtotal    float64 `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Weight      float64 `gorm:"type:decimal(8,3)" json:"weight"`
}
