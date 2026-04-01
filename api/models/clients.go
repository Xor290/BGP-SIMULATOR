package models

import "github.com/golang-jwt/jwt/v5"

type Client struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClientClaims struct {
	ClientID int    `json:"client_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
