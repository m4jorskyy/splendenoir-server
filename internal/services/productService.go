package services

import (
	"errors"
	"splendenoir-server/internal/models/product"
)

type ProductRepository interface {
	GetAllProducts() ([]*product.Product, error)

	GetProductByID(id int64) (*product.Product, error)
}

type ProductService struct {
	repository ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repository: repo}
}

func (r *ProductService) GetAllProducts() ([]*product.Product, error) {
	products, productsErr := r.repository.GetAllProducts()

	if productsErr != nil {
		return nil, productsErr
	}

	return products, nil
}

func (r *ProductService) GetProductByID(id int64) (*product.Product, error) {
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}

	p, productErr := r.repository.GetProductByID(id)

	if productErr != nil {
		return nil, productErr
	}

	return p, nil
}
