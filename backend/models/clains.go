package models

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}
