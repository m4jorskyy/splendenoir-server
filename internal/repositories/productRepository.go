package repositories

import (
	"database/sql"
	"splendenoir-server/internal/models/product"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllProducts() ([]*product.Product, error) {
	var products []*product.Product
	rows, rowsErr := r.db.Query("SELECT id, name, material, fineness, length, width, price, quantity FROM products WHERE deleted_at IS NULL")

	if rowsErr != nil {
		return nil, rowsErr
	}

	defer func(rows *sql.Rows) {
		errCloseRows := rows.Close()
		if errCloseRows != nil {
			return
		}
	}(rows)

	for rows.Next() {
		p := &product.Product{}

		errScan := rows.Scan(&p.ID, &p.Name, &p.Material, &p.Fineness, &p.Length, &p.Width, &p.Price, &p.Quantity)

		if errScan != nil {
			return nil, errScan
		}

		products = append(products, p)
	}

	if errRows := rows.Err(); errRows != nil {
		return nil, errRows
	}

	return products, nil
}
