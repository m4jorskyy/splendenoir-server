package data

import "splendenoir-server/internal/models/data/cart"

type Order struct {
	ProfileID int64            `json:"profile_id"`
	Items     []*cart.CartItem `json:"items"`
}
