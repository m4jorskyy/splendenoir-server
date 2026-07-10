package services

import (
	"context"
	"math"
	"splendenoir-server/internal/models/data/cart"
	"strconv"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, profileID int64, addressID int64, items []*cart.CartItem) (int64, float64, error)
}

type OrderService struct {
	repository OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repository: repo}
}

func (r *OrderService) CreateOrder(ctx context.Context, profileID int64, addressID int64, items []*cart.CartItem) (int64, int64, error) {
	orderID, amount, errCreateOrder := r.repository.CreateOrder(ctx, profileID, addressID, items)

	if errCreateOrder != nil {
		return -1, -1, errCreateOrder
	}

	amountInCents := int64(math.Round(amount * 100))
	strconv.Itoa(int(orderID))

	return orderID, amountInCents, nil
}
