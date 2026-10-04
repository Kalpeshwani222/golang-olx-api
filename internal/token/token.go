package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type Claims struct {
	UserID string `json:"user_id"` 
	jwt.RegisteredClaims        //standard JWT fields likes ExpiresAt, IssuedAt, Subject
}

// Generate creates and signs a new JWT token for the given userID
func Generate(userID string, secret []byte) (string, error) {
	expiresAt := time.Now().Add(30 * time.Minute) // token valid for 30m

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),  //
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := t.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("token.Generate: %w", err)
	}

	return signed, nil
}

// Parse validates a JWT token
func Parse(tokenStr string, secret []byte) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token is expired")
		}
		return nil, fmt.Errorf("token.Parse: %w", err)
	}

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
