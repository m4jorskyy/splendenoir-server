package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"splendenoir-server/internal/models/data/order"
	"splendenoir-server/internal/services"
	"strconv"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (s *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if http.MethodPost != r.Method {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	profileID := r.Context().Value(userContextKey).(int64)

	var request order.OrderRequest
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&request)

	if errDecode != nil {
		http.Error(w, fmt.Sprintf("Error decoding: %s", errDecode), http.StatusBadRequest)
		return
	}

	clientSecret, orderID, amountInCents, errCreateOrder := s.service.CreateOrder(r.Context(), profileID, request.AddressID, request.Items)

	if errCreateOrder != nil || amountInCents == -1 || orderID == -1 || clientSecret == "" {
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	response := &order.OrderResponse{ClientSecret: clientSecret, OrderID: orderID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	errEncode := encoder.Encode(response)

	if errEncode != nil {
		http.Error(w, "Error encoding", http.StatusInternalServerError)
		return
	}

}

func (s *OrderHandler) PaymentConfirmation(w http.ResponseWriter, r *http.Request) {
	if http.MethodPost != r.Method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 65536)

	payload, errRead := io.ReadAll(r.Body)
	if errRead != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %s", errRead), http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("Stripe-Signature")

	event, errEvent := webhook.ConstructEvent(payload, signature, os.Getenv("STRIPE_WEBHOOK_SECRET"))

	if errEvent != nil {
		http.Error(w, fmt.Sprintf("Error constructing event: %s", errEvent), http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		errUnmarshall := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if errUnmarshall != nil {
			http.Error(w, fmt.Sprintf("Error unmarshall: %s", errUnmarshall), http.StatusBadRequest)
			return
		}

		orderID, errConversion := strconv.ParseInt(paymentIntent.Metadata["orderID"], 10, 64)
		if errConversion != nil {
			http.Error(w, fmt.Sprintf("Error conversion: %s", errConversion), http.StatusInternalServerError)
			return
		}

		errStatus := s.service.PaymentStatus(orderID, "SUCCESS")
		if errStatus != nil {
			http.Error(w, fmt.Sprintf("Error status: %s", errStatus), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		break
	case "payment_intent.payment_failed":
		//TODO
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusOK)
		return
	}

}
