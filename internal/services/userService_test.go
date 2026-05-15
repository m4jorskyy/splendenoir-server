package services

import (
	"splendenoir-server/internal/models/user"
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

func TestRegisterNewUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{}
	svc := NewUserService("secret", mockRepo)
	registerUserID, errRegister := svc.RegisterNewUser("John", "Doe", "john@doe.com", "johndoe", "johndoepassword")
	if errRegister != nil {
		t.Errorf("errRegister: %s", errRegister)
	}

	if registerUserID != 1 {
		t.Errorf("registerUserID is not 1")
	}

}

func TestLoginUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{}
	svc := NewUserService("secret", mockRepo)
	jwtToken, returnedFirstName, errLogin := svc.LoginUser("johndoe", "johndoepassword")

	if errLogin != nil {
		t.Errorf("errLogin: %s", errLogin)
	}

	if returnedFirstName != "John" {
		t.Errorf("returnedFirstName is not John")
	}

	if jwtToken == "" {
		t.Errorf("jwtToken is empty")
	}
}

func TestLoginUser_WrongPassword(t *testing.T) {
	mockRepo := &MockUserRepository{}
	svc := NewUserService("secret", mockRepo)
	jwtToken, returnedFirstName, errLogin := svc.LoginUser("johndoe", "johndoewrongpassword")

	if errLogin == nil {
		t.Errorf("logging in went through")
	}

	if returnedFirstName != "" {
		t.Errorf("returnedFirstName is not empty")
	}

	if jwtToken != "" {
		t.Errorf("jwtToken is not empty")
	}
}
