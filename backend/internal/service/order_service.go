package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"
	"ecommerce-backend/pkg/easyparcel"
	"ecommerce-backend/pkg/toyyibpay"
)

type OrderService struct {
	orderRepo      *repository.OrderRepository
	productRepo    *repository.ProductRepository
	userRepo       *repository.UserRepository
	toyyibClient   *toyyibpay.Client
	easyClient     *easyparcel.Client
	frontendURL    string
	backendURL     string
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
	userRepo *repository.UserRepository,
	toyyibClient *toyyibpay.Client,
	easyClient *easyparcel.Client,
	frontendURL, backendURL string,
) *OrderService {
	return &OrderService{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		userRepo:     userRepo,
		toyyibClient: toyyibClient,
		easyClient:   easyClient,
		frontendURL:  frontendURL,
		backendURL:   backendURL,
	}
}

type OrderItemInput struct {
	ProductID uint `json:"productId"`
	Quantity  int  `json:"quantity"`
}

type CreateOrderInput struct {
	// Maklumat Pelanggan (Boleh diisi oleh Guest atau auto-fill dari profil)
	CustomerName  string `json:"customerName"`
	CustomerEmail string `json:"customerEmail"`
	CustomerPhone string `json:"customerPhone"`

	// Alamat Penghantaran
	ShippingAddress string `json:"shippingAddress"`
	City            string `json:"city"`
	State           string `json:"state"`
	Postcode        string `json:"postcode"`

	CourierName  string           `json:"courierName"`
	ShippingRate float64          `json:"shippingRate"`
	Notes        string           `json:"notes"`
	Items        []OrderItemInput `json:"items"`
}

type CreateOrderResponse struct {
	Order      *models.Order `json:"order"`
	PaymentURL string        `json:"paymentUrl"`
}

func (s *OrderService) CreateOrder(input CreateOrderInput, claims *middleware.UserClaims) (*CreateOrderResponse, error) {
	if len(input.Items) == 0 {
		return nil, errors.New("pesanan mestilah mempunyai sekurang-kurangnya 1 item")
	}

	var userID *uint
	// Jika user sedang login (bukan guest)
	if claims != nil && claims.UserID > 0 {
		userID = &claims.UserID
	}

	// Validasi input asas
	if strings.TrimSpace(input.CustomerName) == "" || strings.TrimSpace(input.CustomerEmail) == "" || strings.TrimSpace(input.CustomerPhone) == "" {
		return nil, errors.New("nama, emel, dan nombor telefon pelanggan diperlukan")
	}
	if strings.TrimSpace(input.ShippingAddress) == "" || strings.TrimSpace(input.Postcode) == "" {
		return nil, errors.New("alamat penghantaran dan poskod lengkap diperlukan")
	}

	var totalWeight float64
	var subtotal float64
	var orderItems []models.OrderItem

	// Semak setiap produk dalam inventory
	for _, item := range input.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("kuantiti item tidak sah")
		}

		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("produk ID %d tidak ditemui", item.ProductID)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("stok tidak mencukupi untuk %s (baki: %d)", product.Name, product.Stock)
		}

		itemSubtotal := product.Price * float64(item.Quantity)
		subtotal += itemSubtotal
		totalWeight += product.Weight * float64(item.Quantity)

		orderItems = append(orderItems, models.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			Subtotal:    itemSubtotal,
			Weight:      product.Weight,
		})

		// Tolak stok produk
		product.Stock -= item.Quantity
		_ = s.productRepo.Update(product)
	}

	shippingFee := input.ShippingRate
	if shippingFee <= 0 {
		shippingFee = 8.00 // Default kadar asas pos domestik
	}

	totalAmount := subtotal + shippingFee
	orderNum := fmt.Sprintf("ORD-%s-%04d", time.Now().Format("20060102"), rand.Intn(9999))

	order := &models.Order{
		OrderNumber:     orderNum,
		UserID:          userID,
		CustomerName:    strings.TrimSpace(input.CustomerName),
		CustomerEmail:   strings.TrimSpace(input.CustomerEmail),
		CustomerPhone:   strings.TrimSpace(input.CustomerPhone),
		ShippingAddress: strings.TrimSpace(input.ShippingAddress),
		City:            strings.TrimSpace(input.City),
		State:           strings.TrimSpace(input.State),
		Postcode:        strings.TrimSpace(input.Postcode),
		Subtotal:        subtotal,
		ShippingFee:     shippingFee,
		TotalAmount:     totalAmount,
		TotalWeight:     totalWeight,
		Status:          models.StatusPendingPayment,
		Notes:           input.Notes,
		Items:           orderItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("gagal mencipta pesanan: %w", err)
	}

	// Rekod maklumat penghantaran awal (EasyParcel)
	courier := input.CourierName
	if courier == "" {
		courier = "J&T Express"
	}
	shipment := &models.Shipment{
		OrderID:      order.ID,
		Provider:     "easyparcel",
		CourierName:  courier,
		ShippingRate: shippingFee,
		Status:       models.ShipmentPending,
	}
	_ = s.orderRepo.CreateShipment(shipment)

	// Bina bil ToyyibPay
	returnURL := fmt.Sprintf("%s/order/status?orderNumber=%s", s.frontendURL, order.OrderNumber)
	callbackURL := fmt.Sprintf("%s/api/webhooks/toyyibpay", s.backendURL)

	billCode, payURL, err := s.toyyibClient.CreateBill(toyyibpay.CreateBillParams{
		BillName:            fmt.Sprintf("Order %s", order.OrderNumber),
		BillDescription:     fmt.Sprintf("Pembayaran Pesanan %s di Pas Raisy Store", order.OrderNumber),
		AmountRM:            totalAmount,
		ExternalReferenceNo: order.OrderNumber,
		CustomerName:        order.CustomerName,
		CustomerEmail:       order.CustomerEmail,
		CustomerPhone:       order.CustomerPhone,
		ReturnURL:           returnURL,
		CallbackURL:         callbackURL,
	})

	if err != nil {
		// Sekiranya toyyibpay keys belum diisi di .env, kita sediakan mock payment URL untuk testing
		billCode = fmt.Sprintf("MOCK-%d", time.Now().Unix())
		payURL = fmt.Sprintf("%s/mock-payment?billCode=%s&order=%s", s.frontendURL, billCode, order.OrderNumber)
	}

	payment := &models.Payment{
		OrderID:    order.ID,
		Provider:   "toyyibpay",
		BillCode:   billCode,
		Amount:     totalAmount,
		PaymentURL: payURL,
		Status:     models.PaymentStatusPending,
	}
	_ = s.orderRepo.CreatePayment(payment)

	return &CreateOrderResponse{
		Order:      order,
		PaymentURL: payURL,
	}, nil
}

func (s *OrderService) GetOrderByNumber(orderNumber string) (*models.Order, error) {
	return s.orderRepo.FindByOrderNumber(orderNumber)
}

func (s *OrderService) GetUserOrders(userID uint) ([]models.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

func (s *OrderService) CheckShippingRates(receiverPostcode string, weightKg float64) ([]easyparcel.CourierRate, error) {
	senderPostcode := "50450" // Poskod gudang/kedai Pas Raisy
	return s.easyClient.CheckRates(easyparcel.RateCheckRequest{
		SenderPostcode:   senderPostcode,
		ReceiverPostcode: receiverPostcode,
		WeightKg:         weightKg,
	})
}

func (s *OrderService) HandleToyyibPayWebhook(billCode string, statusStr string, rawPayload string) error {
	payment, err := s.orderRepo.FindPaymentByBillCode(billCode)
	if err != nil {
		return err
	}

	var paymentStatus models.PaymentStatus
	var orderStatus models.OrderStatus

	// ToyyibPay status: 1 = success, 2 = pending, 3 = fail
	if statusStr == "1" {
		paymentStatus = models.PaymentStatusSuccess
		orderStatus = models.StatusPaid
		now := time.Now()
		payment.PaidAt = &now
	} else if statusStr == "3" {
		paymentStatus = models.PaymentStatusFailed
		orderStatus = models.StatusCancelled
	} else {
		paymentStatus = models.PaymentStatusPending
		orderStatus = models.StatusPendingPayment
	}

	_ = s.orderRepo.UpdatePaymentStatus(billCode, paymentStatus, rawPayload)
	_ = s.orderRepo.UpdateStatus(payment.OrderID, orderStatus)
	return nil
}
