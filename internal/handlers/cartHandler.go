package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"splendenoir-server/internal/models/data/cart"
	"splendenoir-server/internal/services"
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (s *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	profileID := r.Context().Value(userContextKey).(int64)

	var item *cart.CartItem
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&item)

	if errDecode != nil {
		http.Error(w, fmt.Sprintf("Error decoding: %s", errDecode), http.StatusBadRequest)
		return
	}

	errAdd := s.service.AddToCart(r.Context(), profileID, item)

	if errAdd != nil {
		http.Error(w, "Error adding to cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	profileID := r.Context().Value(userContextKey).(int64)

	var itemToRemove *cart.CartItem
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&itemToRemove)

	if errDecode != nil {
		http.Error(w, fmt.Sprintf("Error decoding: %s", errDecode), http.StatusBadRequest)
		return
	}

	errRemove := s.service.RemoveFromCart(r.Context(), profileID, itemToRemove)

	if errRemove != nil {
		http.Error(w, "Error removing from cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	profileID := r.Context().Value(userContextKey).(int64)

	c, errCart := s.service.GetCart(r.Context(), profileID)

	if errCart != nil {
		http.Error(w, "Error getting cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	errEncode := encoder.Encode(c)

	if errEncode != nil {
		http.Error(w, "Error encoding", http.StatusInternalServerError)
		return
	}
}
