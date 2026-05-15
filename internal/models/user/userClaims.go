package user

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	ProfileID int64 `json:"profile_id"`
	jwt.RegisteredClaims
}
