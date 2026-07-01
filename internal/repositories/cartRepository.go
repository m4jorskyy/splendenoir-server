package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"splendenoir-server/internal/models/data/cart"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type CartRepository struct {
	rdb *redis.Client
}

func NewCartRepository(rdb *redis.Client) *CartRepository {
	return &CartRepository{rdb: rdb}
}

func (r *CartRepository) SetCart(ctx context.Context, profileID int64, cart *cart.Cart) error {
	marshal, errMarshal := json.Marshal(cart)
	if errMarshal != nil {
		return errMarshal
	}

	set := r.rdb.Set(ctx, "cart:"+strconv.FormatInt(profileID, 10), marshal, 15*time.Minute)

	if set.Err() != nil {
		return set.Err()
	}

	return nil
}

func (r *CartRepository) GetCart(ctx context.Context, profileID int64) (*cart.Cart, error) {
	cartRdb := r.rdb.Get(ctx, "cart:"+strconv.FormatInt(profileID, 10))

	if cartRdb.Err() != nil {
		if errors.Is(cartRdb.Err(), redis.Nil) {
			return &cart.Cart{}, nil
		}

		return nil, cartRdb.Err()
	}

	c := &cart.Cart{}
	cartBytes, errCartBytes := cartRdb.Bytes()
	if errCartBytes != nil {
		return nil, errCartBytes
	}

	errUnmarshal := json.Unmarshal(cartBytes, c)

	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return c, nil
}
