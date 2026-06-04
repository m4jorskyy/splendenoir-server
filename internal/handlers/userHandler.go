package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"splendenoir-server/internal/models/data"
	"splendenoir-server/internal/services"
	"time"
)

type UserHandler struct {
	jwtSecret string
	service   *services.UserService
}

func NewUserHandler(jwtSecret string, service *services.UserService) *UserHandler {
	return &UserHandler{jwtSecret: jwtSecret, service: service}
}

type contextKey string

const userContextKey contextKey = "profileID"

func (s *UserHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, errCookie := r.Cookie("Token")

		if errCookie != nil {
			http.Error(w, "Cookie not found", http.StatusUnauthorized)
			return
		}

		returnedID, errVerify := s.service.VerifyToken(token.Value)

		if errVerify != nil {
			http.Error(w, "Error verifying", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, returnedID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	}
}

func (s *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	var request data.RegisterRequest
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&request)

	if errDecode != nil {
		http.Error(w, fmt.Sprintf("Error decoding: %s", errDecode), http.StatusBadRequest)
		return
	}

	returnedID, errRegister := s.service.RegisterNewUser(request.FirstName, request.LastName, request.Email,
		request.Login, request.Password)

	if errRegister != nil {
		http.Error(w, "Error registering", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	errEncoding := encoder.Encode(map[string]int64{"id": returnedID})

	if errEncoding != nil {
		http.Error(w, fmt.Sprintf("Error encoding: %s", errEncoding), http.StatusInternalServerError)
		return
	}
}

func (s *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	var request data.LoginRequest
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&request)

	if errDecode != nil {
		http.Error(w, fmt.Sprintf("Error decoding: %s", errDecode), http.StatusBadRequest)
		return
	}

	jwtToken, firstName, errLogin := s.service.LoginUser(request.Login, request.Password)

	if errLogin != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{Name: "Token", Value: jwtToken, Expires: time.Now().Add(15 * time.Minute), HttpOnly: true}
	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	errEncoding := encoder.Encode(map[string]string{"first_name": firstName})

	if errEncoding != nil {
		http.Error(w, fmt.Sprintf("Error encoding: %s", errEncoding), http.StatusInternalServerError)
		return
	}
}
