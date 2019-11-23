package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}
