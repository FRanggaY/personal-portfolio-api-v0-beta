package middlewares

import (
	"errors"
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/config"
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]string{"message": "Unauthorized"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}
		// take token value
		tokenString := c.Value

		claims := &config.JWTClaim{}
		// parsing token jwt
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				// token invalid
				response := map[string]string{"message": "Token invalid"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				// Invalid signature
				response := map[string]string{"message": "Invalid signature"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
				// Token is either expired or not active yet
				response := map[string]string{"message": "Token is expired or not active yet"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			default:
				// Token is either expired or not active yet
				response := map[string]string{"message": "Couldn't handle this token"}
				helper.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			}
		}

		if !token.Valid {
			response := map[string]string{"message": "Unauthorized"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}
		next.ServeHTTP(w, r)
	})
}
