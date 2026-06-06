package repositories

import (
	"database/sql"
	"errors"
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

func (r *ProductRepository) GetProductByID(id int64) (*product.Product, error) {
	p := &product.Product{}
	errProduct := r.db.QueryRow("SELECT id, name, material, fineness, length, width, price, quantity FROM products WHERE id = $1 AND deleted_at IS NULL", id).Scan(&p.ID, &p.Name, &p.Material, &p.Fineness, &p.Length, &p.Width, &p.Price, &p.Quantity)

	if errProduct != nil {

		if errors.Is(errProduct, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, errProduct
	}

	return p, nil
}
