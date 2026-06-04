package handlers

import (
	"net/http"
	"net/http/httptest"
	"splendenoir-server/internal/models/user"
	"splendenoir-server/internal/services"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct{}

func (m *MockUserRepository) RegisterUser(profile *user.UserProfile, creds *user.UserCredentials) (int64, error) {
	return 1, nil
}

func (m *MockUserRepository) LoginUser(username string) (int64, string, string, error) {
	password := "johndoepassword"
	passwordHash, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if errHash != nil {
		return 0, "", "", errHash
	}

	return 1, string(passwordHash), "John", nil
}

func TestUserHandler_RegisterUser(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/register/", strings.NewReader("{\n  \"first_name\": \"John\",\n  \"last_name\": \"Doe\",\n  \"email\": \"john.doe@gmail.com\",\n  \"username\": \"johndoe\",\n  \"password\": \"johndoe123\"\n}"))

	recorder := httptest.NewRecorder()

	mockRepo := &MockUserRepository{}
	svc := services.NewUserService("secret", mockRepo)
	handler := NewUserHandler("secret", svc)
	handler.RegisterUser(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Register Code is not 201")
	}
}

func TestUserHandler_LoginUser(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/login/", strings.NewReader("{\n	\"login\": \"johndoe\",\n  \"password\": \"johndoepassword\"\n}"))

	recorder := httptest.NewRecorder()
	mockRepo := &MockUserRepository{}
	svc := services.NewUserService("secret", mockRepo)

	handler := NewUserHandler("secret", svc)
	handler.LoginUser(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Login code is not 200")
	}
}
