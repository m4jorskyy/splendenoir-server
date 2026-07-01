package services

import (
	"context"
	"errors"
	"splendenoir-server/internal/models/data/cart"
)

type CartRepository interface {
	SetCart(ctx context.Context, profileID int64, cart *cart.Cart) error
	GetCart(ctx context.Context, profileID int64) (*cart.Cart, error)
}

type CartService struct {
	repository CartRepository
}

func NewCartService(repo CartRepository) *CartService {
	return &CartService{repository: repo}
}

func (r *CartService) SetCart(ctx context.Context, profileID int64, cart *cart.Cart) error {
	errSet := r.repository.SetCart(ctx, profileID, cart)

	if errSet != nil {
		return errSet
	}

	return nil
}

func (r *CartService) GetCart(ctx context.Context, profileID int64) (*cart.Cart, error) {
	c, errCart := r.repository.GetCart(ctx, profileID)

	if errCart != nil {
		return nil, errCart
	}

	return c, nil
}

func (r *CartService) AddToCart(ctx context.Context, profileID int64, newItem *cart.CartItem) error {
	c, errCart := r.repository.GetCart(ctx, profileID)

	if errCart != nil {
		return errCart
	}

	for i := 0; i < len(c.Items); i++ {
		item := c.Items[i]

		if item.ID == newItem.ID {
			item.Quantity += newItem.Quantity
			errSet := r.repository.SetCart(ctx, profileID, c)

			if errSet != nil {
				return errSet
			}

			return nil
		}
	}

	c.Items = append(c.Items, newItem)
	errSet := r.repository.SetCart(ctx, profileID, c)

	if errSet != nil {
		return errSet
	}

	return nil
}

func (r *CartService) RemoveFromCart(ctx context.Context, profileID int64, itemToRemove *cart.CartItem) error {
	c, errGetCart := r.repository.GetCart(ctx, profileID)

	if errGetCart != nil {
		return errGetCart
	}

	if c == nil {
		return errors.New("no cart found")
	}

	for i := 0; i < len(c.Items); i++ {
		item := c.Items[i]

		if item.ID == itemToRemove.ID {
			if (item.Quantity - itemToRemove.Quantity) < 0 {
				return errors.New("cannot remove more than it is")
			}

			item.Quantity = item.Quantity - itemToRemove.Quantity

			if item.Quantity == 0 {
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
			}

			errSetCart := r.repository.SetCart(ctx, profileID, c)
			if errSetCart != nil {
				return errSetCart
			}

			return nil
		}
	}

	return nil

}
