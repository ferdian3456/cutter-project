package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserId int `json:"userId"`
	jwt.RegisteredClaims
}
