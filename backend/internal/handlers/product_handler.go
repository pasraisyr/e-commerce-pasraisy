package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func parseID(r *http.Request) (uint, error) {
	val := r.PathValue("id")
	if val == "" {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) > 0 {
			val = parts[len(parts)-1]
		}
	}
	id, err := strconv.ParseUint(val, 10, 32)
	return uint(id), err
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	products, err := h.productService.GetAll(true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := parseID(r)
	if err != nil {
		http.Error(w, `{"error":"ID produk tidak sah"}`, http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetByID(id)
	if err != nil {
		http.Error(w, `{"error":"Produk tidak dijumpai"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, `{"error":"Format data tidak sah"}`, http.StatusBadRequest)
		return
	}

	if err := h.productService.Create(&product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := parseID(r)
	if err != nil {
		http.Error(w, `{"error":"ID produk tidak sah"}`, http.StatusBadRequest)
		return
	}

	var input models.Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"Format data tidak sah"}`, http.StatusBadRequest)
		return
	}

	updated, err := h.productService.Update(id, &input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(updated)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := parseID(r)
	if err != nil {
		http.Error(w, `{"error":"ID produk tidak sah"}`, http.StatusBadRequest)
		return
	}

	if err := h.productService.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Produk berjaya dipadam"})
}
