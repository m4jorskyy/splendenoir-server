package handlers

import (
	"encoding/json"
	"net/http"
	"splendenoir-server/internal/services"
	"strconv"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (s *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	products, errProducts := s.service.GetAllProducts()

	if errProducts != nil {
		http.Error(w, "Error getting products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	errEncoding := encoder.Encode(products)

	if errEncoding != nil {
		http.Error(w, "Error encoding", http.StatusInternalServerError)
		return
	}
}

func (s *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	ID, errID := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if errID != nil {
		http.Error(w, "Error parsing ID", http.StatusInternalServerError)
		return
	}

	p, errProduct := s.service.GetProductByID(ID)

	if errProduct != nil {
		http.Error(w, "Error getting product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	errEncoding := encoder.Encode(p)

	if errEncoding != nil {
		http.Error(w, "Error encoding", http.StatusInternalServerError)
		return
	}
}
