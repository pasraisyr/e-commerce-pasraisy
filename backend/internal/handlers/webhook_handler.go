package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"ecommerce-backend/internal/service"
)

type WebhookHandler struct {
	orderService *service.OrderService
}

func NewWebhookHandler(orderService *service.OrderService) *WebhookHandler {
	return &WebhookHandler{orderService: orderService}
}

// ToyyibPayCallback receives webhook notifications when payment is completed
func (h *WebhookHandler) ToyyibPayCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	rawPayload := string(bodyBytes)

	// ToyyibPay usually sends either x-www-form-urlencoded or multipart
	_ = r.ParseForm()
	billCode := r.FormValue("billcode")
	status := r.FormValue("status_id")
	if status == "" {
		status = r.FormValue("status")
	}

	log.Printf("📩 [ToyyibPay Webhook] Diterima: BillCode=%s, Status=%s\nPayload: %s\n", billCode, status, rawPayload)

	if billCode != "" {
		if err := h.orderService.HandleToyyibPayWebhook(billCode, status, rawPayload); err != nil {
			log.Printf("⚠️ Gagal memproses webhook ToyyibPay: %v\n", err)
		}
	}

	// ToyyibPay requires an "OK" response to confirm receipt
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
