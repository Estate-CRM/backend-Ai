package middlewares

import (
	"errors"

	"github.com/Estate-CRM/backend-go/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyAccessToken(tokenStr string) (*JWTClaims, error) {
	cfg, _ := config.LoadConfig() // Load your secret
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Access_token), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}
	return claims, nil
}
