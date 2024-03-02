package helper

import (
	"errors"
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/config"
	"github.com/golang-jwt/jwt/v5"
)

func GetJWTClaim(r *http.Request) (config.JWTClaim, error) {
	c, _ := r.Cookie("token")

	tokenString := c.Value

	claims := &config.JWTClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})
	if err != nil {
		return config.JWTClaim{}, err
	}

	if !token.Valid {
		return config.JWTClaim{}, errors.New("invalid token")
	}

	return *claims, nil
}
