package repository

import (
	"ecommerce-backend/internal/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").Preload("Payment").Preload("Shipment").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").Preload("Payment").Preload("Shipment").
		Where("order_number = ?", orderNumber).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Items").Preload("Payment").Preload("Shipment").
		Where("user_id = ?", userID).Order("id desc").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Items").Preload("Payment").Preload("Shipment").
		Order("id desc").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateStatus(orderID uint, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *OrderRepository) CreatePayment(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *OrderRepository) UpdatePaymentStatus(billCode string, status models.PaymentStatus, rawPayload string) error {
	return r.db.Model(&models.Payment{}).
		Where("bill_code = ?", billCode).
		Updates(map[string]interface{}{
			"status":      status,
			"raw_payload": rawPayload,
		}).Error
}

func (r *OrderRepository) FindPaymentByBillCode(billCode string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("bill_code = ?", billCode).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *OrderRepository) CreateShipment(shipment *models.Shipment) error {
	return r.db.Create(shipment).Error
}
