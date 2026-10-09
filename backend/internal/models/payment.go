package models

import "time"

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "successful"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	OrderID       uint          `gorm:"uniqueIndex;not null" json:"orderId"`
	Provider      string        `gorm:"size:50;default:'toyyibpay'" json:"provider"`
	BillCode      string        `gorm:"size:100;index" json:"billCode"`
	TransactionID string        `gorm:"size:100;index" json:"transactionId"`
	Amount        float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentURL    string        `gorm:"size:500" json:"paymentUrl"`
	Status        PaymentStatus `gorm:"size:30;default:'pending'" json:"status"`
	PaidAt        *time.Time    `json:"paidAt"`
	RawPayload    string        `gorm:"type:text" json:"rawPayload"` // Audit data dari webhook ToyyibPay
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}
