package model

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Random letters.
// It's better to put it in env but who cares lol
var signingString = []byte("SomethingReallyStupid")

type TokenClaims struct {
	Role     UserRole `json:"role"`
	Username string   `json:"username"`
	jwt.RegisteredClaims
}

func DecodeToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return signingString, nil
}

func GenerateToken(user *User) (string, error) {
	claims := &TokenClaims{
		Role:     user.Role,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 48)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ID:        user.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingString)

}
