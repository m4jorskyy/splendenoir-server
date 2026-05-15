package services

import (
	"errors"
	"splendenoir-server/internal/models/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	RegisterUser(profile *user.UserProfile, creds *user.UserCredentials) (int64, error)
	LoginUser(username string) (int64, string, string, error)
}

type UserService struct {
	jwtSecret  string
	repository UserRepository
}

func NewUserService(jwtSecret string, repo UserRepository) *UserService {
	return &UserService{jwtSecret: jwtSecret, repository: repo}
}

func GenerateToken(profileID int64, secret string) (string, error) {
	registeredClaims := jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute))}
	claims := &user.UserClaims{ProfileID: profileID, RegisteredClaims: registeredClaims}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, errSignedString := token.SignedString([]byte(secret))

	if errSignedString != nil {
		return "", errSignedString
	}

	return signedToken, nil
}

func (r *UserService) VerifyToken(jwtToken string) (int64, error) {
	claims := &user.UserClaims{}

	parsedClaims, errParseClaims := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.jwtSecret), nil
	})

	if errParseClaims != nil {
		return 0, errParseClaims
	}

	if !parsedClaims.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.ProfileID, nil

}

func (r *UserService) RegisterNewUser(firstName string, lastName string, email string, login string, password string) (int64, error) {
	passwordHash, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if errHash != nil {
		return 0, errHash
	}

	profile := &user.UserProfile{FirstName: firstName, LastName: lastName, Email: email}
	creds := &user.UserCredentials{Username: login, PasswordHash: string(passwordHash)}

	registerUserID, errRegister := r.repository.RegisterUser(profile, creds)
	if errRegister != nil {
		return 0, errRegister
	}

	return registerUserID, nil
}

func (r *UserService) LoginUser(username string, password string) (string, string, error) {
	var returnedUserID int64
	var returnedPasswordHash string
	var returnedFirstName string
	var errReturn error

	returnedUserID, returnedPasswordHash, returnedFirstName, errReturn = r.repository.LoginUser(username)

	if errReturn != nil {
		return "", "", errReturn
	}

	errCompare := bcrypt.CompareHashAndPassword([]byte(returnedPasswordHash), []byte(password))

	if errCompare != nil {
		return "", "", errCompare
	}

	jwtToken, errToken := GenerateToken(returnedUserID, r.jwtSecret)

	if errToken != nil {
		return "", "", errToken
	}

	return jwtToken, returnedFirstName, nil
}
