package services

import (
	"splendenoir-server/internal/models/product"
	"testing"
)

type MockProductRepository struct{}

func (m *MockProductRepository) GetAllProducts() ([]*product.Product, error) {
	p := &product.Product{ID: 1, Name: "Test Product", Material: "Silver", Fineness: "925", Length: 10, Width: 10, Price: 1000, Quantity: 10}

	products := make([]*product.Product, 0)
	products = append(products, p)

	return products, nil
}

func (m *MockProductRepository) GetProductByID(ID int64) (*product.Product, error) {
	p := &product.Product{ID: 1, Name: "Test Product", Material: "Silver", Fineness: "925", Length: 10, Width: 10, Price: 1000, Quantity: 10}

	return p, nil
}

func TestProductService_GetAllProducts(t *testing.T) {
	mockRepo := &MockProductRepository{}
	svc := NewProductService(mockRepo)
	products, errProducts := svc.GetAllProducts()

	if errProducts != nil {
		t.Errorf("errProducts: %s", errProducts)
	}

	if products[0].ID != 1 {
		t.Errorf("productID is not 1")
	}
}

func TestProductService_GetProductByID(t *testing.T) {
	mockRepo := &MockProductRepository{}
	svc := NewProductService(mockRepo)
	p, errProduct := svc.GetProductByID(1)

	if errProduct != nil {
		t.Errorf("errProducts: %s", errProduct)
	}

	if p.ID != 1 {
		t.Errorf("productID is not 1")
	}
}
