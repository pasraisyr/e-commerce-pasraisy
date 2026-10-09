package models

import "time"

type ShipmentStatus string

const (
	ShipmentPending   ShipmentStatus = "pending"
	ShipmentBooked    ShipmentStatus = "booked"
	ShipmentInTransit ShipmentStatus = "in_transit"
	ShipmentDelivered ShipmentStatus = "delivered"
	ShipmentCancelled ShipmentStatus = "cancelled"
)

type Shipment struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	OrderID           uint           `gorm:"uniqueIndex;not null" json:"orderId"`
	Provider          string         `gorm:"size:50;default:'easyparcel'" json:"provider"`
	CourierName       string         `gorm:"size:100" json:"courierName"`       // Cth: "J&T Express", "Pos Laju", "DHL"
	ServiceCode       string         `gorm:"size:100" json:"serviceCode"`       // Kod servis EasyParcel
	EasyParcelOrderID string         `gorm:"size:100;index" json:"easyParcelOrderId"`
	TrackingNumber    string         `gorm:"size:100;index" json:"trackingNumber"`
	AirwayBillURL     string         `gorm:"size:500" json:"airwayBillUrl"`    // URL PDF salinan consignment note/AWB
	ShippingRate      float64        `gorm:"type:decimal(10,2)" json:"shippingRate"`
	Status            ShipmentStatus `gorm:"size:30;default:'pending'" json:"status"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}
