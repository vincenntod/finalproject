package auth

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("dbceria")

type JWT struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
