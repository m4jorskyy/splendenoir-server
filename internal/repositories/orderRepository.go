package repositories

import (
	"context"
	"database/sql"
	"errors"
	"splendenoir-server/internal/models/data"
	"splendenoir-server/internal/models/data/cart"
	"strconv"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type OrderRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewOrderRepository(db *sql.DB, rdb *redis.Client) *OrderRepository {
	return &OrderRepository{db: db, rdb: rdb}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, profileID int64, addressID int64, items []*cart.CartItem) (int64, float64, error) {
	tx, errTx := r.db.BeginTx(ctx, nil)
	defer func(tx *sql.Tx) {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return
		}
	}(tx)

	if errTx != nil {
		return -1, -1, errTx
	}

	var sum float64
	var itemsID []int64
	var orderID int64
	var products []*data.Product

	itemsMap := make(map[int64]int64)

	for i := 0; i < len(items); i++ {
		itemsID = append(itemsID, items[i].ID)
		itemsMap[items[i].ID] = items[i].Quantity
	}

	rowItemsDB, errRowItemsDB := tx.QueryContext(ctx,
		"SELECT id, name, material, fineness, length, width, price, quantity FROM products WHERE id = ANY($1) AND deleted_at IS NULL", pq.Array(itemsID))

	if errRowItemsDB != nil {
		if errors.Is(errRowItemsDB, sql.ErrNoRows) {
			return -1, -1, sql.ErrNoRows
		}

		return -1, -1, errRowItemsDB
	}

	for rowItemsDB.Next() {
		itemDB := &data.Product{}
		errScan := rowItemsDB.Scan(&itemDB.ID, &itemDB.Name, &itemDB.Material, &itemDB.Fineness, &itemDB.Length, &itemDB.Width,
			&itemDB.Price, &itemDB.Quantity)

		if errScan != nil {
			return -1, -1, errScan
		}

		products = append(products, itemDB)

		sum += itemDB.Price * float64(itemsMap[itemDB.ID])
	}

	if len(products) != len(itemsMap) {
		return -1, -1, errors.New("some products are unavailable")
	}

	errOrderCreate := tx.QueryRowContext(ctx, "INSERT INTO orders (profile_id, address_id, amount) VALUES ($1, $2, $3) RETURNING id", profileID, addressID, sum).Scan(&orderID)

	if errOrderCreate != nil {
		return -1, -1, errOrderCreate
	}

	for _, p := range products {
		_, errInsertOrderProduct := tx.ExecContext(ctx, "INSERT INTO orders_products (order_id, product_id, product_price, quantity) VALUES ($1, $2, $3, $4)", orderID, p.ID, p.Price, itemsMap[p.ID])

		if errInsertOrderProduct != nil {
			return -1, -1, errInsertOrderProduct
		}
	}

	errCommit := tx.Commit()

	if errCommit != nil {
		return -1, -1, errCommit
	}

	errRedisDel := r.rdb.Del(ctx, "cart:"+strconv.FormatInt(profileID, 10)).Err()

	if errRedisDel != nil {
		return -1, -1, errRedisDel
	}

	return orderID, sum, nil
}
