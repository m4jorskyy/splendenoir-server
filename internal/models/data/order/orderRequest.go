package order

import "splendenoir-server/internal/models/data/cart"

type OrderRequest struct {
	AddressID int64            `json:"address_id"`
	Items     []*cart.CartItem `json:"items"`
}
