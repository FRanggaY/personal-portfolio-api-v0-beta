package config

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = os.Getenv("JWT_KEY")

var JWT_KEY = []byte(jwtKey)

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}
