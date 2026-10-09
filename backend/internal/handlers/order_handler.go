package handlers

import (
	"encoding/json"
	"net/http"

	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder supports both Guest Checkout and Authenticated Member Checkout
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input service.CreateOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"Format pesanan tidak sah"}`, http.StatusBadRequest)
		return
	}

	claims := middleware.GetUserClaims(r.Context()) // Boleh jadi nil sekiranya Guest
	resp, err := h.orderService.CreateOrder(input, claims)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetOrderByNumber allows customers (including guests) to track their order status
func (h *OrderHandler) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orderNumber := r.URL.Query().Get("orderNumber")
	if orderNumber == "" {
		http.Error(w, `{"error":"Parameter orderNumber diperlukan"}`, http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrderByNumber(orderNumber)
	if err != nil {
		http.Error(w, `{"error":"Pesanan tidak ditemui"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

// GetMyOrders fetches orders placed by the currently logged-in customer
func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	claims := middleware.GetUserClaims(r.Context())
	if claims == nil {
		http.Error(w, `{"error":"Sila log masuk untuk melihat rekod tempahan anda"}`, http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetUserOrders(claims.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(orders)
}

// GetAllOrders allows Admin to view all store orders
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(orders)
}

// CheckShippingRates calculates real-time shipping costs via EasyParcel
func (h *OrderHandler) CheckShippingRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		ReceiverPostcode string  `json:"receiverPostcode"`
		WeightKg         float64 `json:"weightKg"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Format data tidak sah"}`, http.StatusBadRequest)
		return
	}

	if req.ReceiverPostcode == "" {
		http.Error(w, `{"error":"Poskod destinasi diperlukan"}`, http.StatusBadRequest)
		return
	}

	rates, err := h.orderService.CheckShippingRates(req.ReceiverPostcode, req.WeightKg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(rates)
}
