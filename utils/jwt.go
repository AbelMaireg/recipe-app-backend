package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret is the secret key for signing JWTs (in production, store in env)
const JWTSecret = "my-secret-key"

// Claims defines the JWT claims structure
type Claims struct {
	UserID       uint `json:"user_id"`
	HasuraClaims struct {
		XHasuraUserId       string   `json:"x-hasura-user-id"`
		XHasuraDefaultRole  string   `json:"x-hasura-default-role"`
		XHasuraAllowedRoles []string `json:"x-hasura-allowed-roles"`
	} `json:"https://hasura.io/jwt/claims"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a JWT for a user
func GenerateJWT(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		HasuraClaims: struct {
			XHasuraUserId       string   `json:"x-hasura-user-id"`
			XHasuraDefaultRole  string   `json:"x-hasura-default-role"`
			XHasuraAllowedRoles []string `json:"x-hasura-allowed-roles"`
		}{
			XHasuraUserId:       fmt.Sprintf("%d", userID),
			XHasuraDefaultRole:  "user",
			XHasuraAllowedRoles: []string{"user"},
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// ParseJWT verifies a JWT and returns the claims
func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
