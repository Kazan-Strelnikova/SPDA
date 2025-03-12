package user

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserClaims struct {
	Payload User `json:"payload"`
	jwt.RegisteredClaims
}
