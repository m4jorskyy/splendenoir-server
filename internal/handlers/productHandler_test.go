package handlers

import (
	"net/http"
	"net/http/httptest"
	"splendenoir-server/internal/models/product"
	"splendenoir-server/internal/services"
	"strings"
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

func TestProductHandler_GetAllProducts(t *testing.T) {
	request := httptest.NewRequest("GET", "/api/products/", strings.NewReader(""))

	recorder := httptest.NewRecorder()
	mockRepo := &MockProductRepository{}
	svc := services.NewProductService(mockRepo)
	handler := NewProductHandler(svc)
	handler.GetAllProducts(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetAllProducts recorder code is not 200")
	}
}

func TestProductHandler_GetProductByID(t *testing.T) {
	request := httptest.NewRequest("GET", "/api/products/", strings.NewReader(""))
	request.SetPathValue("id", "1")
	recorder := httptest.NewRecorder()
	mockRepo := &MockProductRepository{}
	svc := services.NewProductService(mockRepo)
	handler := NewProductHandler(svc)
	handler.GetProductByID(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("GetProductByID recorder code is not 200")
	}
}
