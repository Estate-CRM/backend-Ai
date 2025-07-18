package middlewares

import (
	"errors"
	"net/http"
	"strings"
)

func GetVerifiedJWTClaims(r *http.Request) (*JWTClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, errors.New("invalid Authorization format")
	}

	tokenStr := parts[1]

	claims, err := VerifyAccessToken(tokenStr)
	if err != nil {
		return nil, errors.New("invalid token: " + err.Error())
	}

	return claims, nil
}
