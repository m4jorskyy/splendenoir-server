package services

import (
	"context"
	"splendenoir-server/internal/models/data/cart"
	"testing"
)

type MockCartRepository struct {
	c map[int64]*cart.Cart
}

func (m *MockCartRepository) SetCart(ctx context.Context, profileID int64, c *cart.Cart) error {

	if m.c == nil {
		m.c = make(map[int64]*cart.Cart)
	}

	m.c[profileID] = c

	return nil
}

func (m *MockCartRepository) GetCart(ctx context.Context, profileID int64) (*cart.Cart, error) {
	mockCart := m.c[profileID]
	return mockCart, nil
}

func TestCartService_SetCart(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCartRepository{}
	svc := NewCartService(mockRepo)
	items := make([]*cart.CartItem, 0)
	items = append(items, &cart.CartItem{ID: 1, Quantity: 10})
	c := &cart.Cart{Items: items}
	errSetCart := svc.SetCart(ctx, 1, c)

	if errSetCart != nil {
		t.Errorf("errSetCart is not nil: %s", errSetCart)
	}
}

func TestCartService_GetCart(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCartRepository{}
	svc := NewCartService(mockRepo)
	items := make([]*cart.CartItem, 0)
	items = append(items, &cart.CartItem{ID: 1, Quantity: 10})
	c := &cart.Cart{Items: items}
	errSetCart := svc.SetCart(ctx, 1, c)

	if errSetCart != nil {
		t.Errorf("errSetCart is not nil: %s", errSetCart)
	}

	mockCart, errGetCart := svc.GetCart(ctx, 1)

	if errGetCart != nil {
		t.Errorf("errGetCart is not nil: %s", errGetCart)
	}

	if mockCart.Items[0].ID != 1 {
		t.Errorf("cartItemID is not 1")
	}
}

func TestCartService_AddToCart(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCartRepository{}
	svc := NewCartService(mockRepo)
	items := make([]*cart.CartItem, 0)
	items = append(items, &cart.CartItem{ID: 1, Quantity: 10})
	c := &cart.Cart{Items: items}
	errSetCart := svc.SetCart(ctx, 1, c)

	if errSetCart != nil {
		t.Errorf("errSetCart is not nil: %s", errSetCart)
	}

	errAddToCart := svc.AddToCart(ctx, 1, &cart.CartItem{ID: 2, Quantity: 5})

	if errAddToCart != nil {
		t.Errorf("errAddToCart is not nil: %s", errAddToCart)
	}

	getCart, errGetCart := svc.GetCart(ctx, 1)

	if errGetCart != nil {
		t.Errorf("errGetCart is not nil: %s", errGetCart)
	}

	if getCart.Items[1].ID != 2 {
		t.Errorf("cartItemID is not 2")
	}
}

func TestCartService_RemoveFromCartSome(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCartRepository{}
	svc := NewCartService(mockRepo)
	items := make([]*cart.CartItem, 0)
	items = append(items, &cart.CartItem{ID: 1, Quantity: 10})
	c := &cart.Cart{Items: items}
	errSetCart := svc.SetCart(ctx, 1, c)

	if errSetCart != nil {
		t.Errorf("errSetCart is not nil: %s", errSetCart)
	}

	errRemoveFromCart := svc.RemoveFromCart(ctx, 1, &cart.CartItem{ID: 1, Quantity: 3})

	if errRemoveFromCart != nil {
		t.Errorf("errRemoveFromCart is not nil: %s", errRemoveFromCart)
	}

	getCart, errGetCart := svc.GetCart(ctx, 1)

	if errGetCart != nil {
		t.Errorf("errGetCart is not nil: %s", errGetCart)
	}

	if getCart.Items[0].Quantity != 7 {
		t.Errorf("cartItemQuantity is not 7")
	}
}

func TestCartService_RemoveFromCartAll(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockCartRepository{}
	svc := NewCartService(mockRepo)
	items := make([]*cart.CartItem, 0)
	items = append(items, &cart.CartItem{ID: 1, Quantity: 10})
	c := &cart.Cart{Items: items}
	errSetCart := svc.SetCart(ctx, 1, c)

	if errSetCart != nil {
		t.Errorf("errSetCart is not nil: %s", errSetCart)
	}

	errRemoveFromCart := svc.RemoveFromCart(ctx, 1, &cart.CartItem{ID: 1, Quantity: 10})

	if errRemoveFromCart != nil {
		t.Errorf("errRemoveFromCart is not nil: %s", errRemoveFromCart)
	}

	getCart, errGetCart := svc.GetCart(ctx, 1)

	if errGetCart != nil {
		t.Errorf("errGetCart is not nil: %s", errGetCart)
	}

	if len(getCart.Items) != 0 {
		t.Errorf("getCart.Items are not nil")
	}
}
