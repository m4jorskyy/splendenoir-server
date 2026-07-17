package services

import (
	"context"
	"math"
	"splendenoir-server/internal/models/data/cart"
	"strconv"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, profileID int64, addressID int64, items []*cart.CartItem) (int64, float64, error)
	PaymentStatus(orderID int64, status string) error
}

type OrderService struct {
	repository OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repository: repo}
}

func (r *OrderService) CreateOrder(ctx context.Context, profileID int64, addressID int64, items []*cart.CartItem) (string, int64, int64, error) {
	orderID, amount, errCreateOrder := r.repository.CreateOrder(ctx, profileID, addressID, items)

	if errCreateOrder != nil {
		return "", -1, -1, errCreateOrder
	}

	amountInCents := int64(math.Round(amount * 100))

	params := &stripe.PaymentIntentParams{Amount: stripe.Int64(amountInCents), Currency: stripe.String("PLN"), Metadata: map[string]string{"orderID": strconv.Itoa(int(orderID))}}

	intent, errIntent := paymentintent.New(params)

	if errIntent != nil {
		return "", -1, -1, errIntent
	}

	return intent.ClientSecret, orderID, amountInCents, nil
}

func (r *OrderService) PaymentStatus(orderID int64, status string) error {
	errStatus := r.repository.PaymentStatus(orderID, status)

	if errStatus != nil {
		return errStatus
	}

	return nil
}
