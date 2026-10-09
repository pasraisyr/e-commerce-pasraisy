package handlers

import (
	"encoding/json"
	"net/http"

	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req service.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Format data tidak sah"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Register(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Format data tidak sah"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	claims := middleware.GetUserClaims(r.Context())
	if claims == nil {
		http.Error(w, `{"error":"Tidak disahkan"}`, http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetProfile(claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"Pengguna tidak ditemui"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
