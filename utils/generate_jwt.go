package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SignedDetails struct {
	Id    string
	Email string
	Name  string
	Role  string
	jwt.RegisteredClaims
}

func GenerateJwt(id uuid.UUID, email string, name string, role string) (string, error) {

	jwt_secret_key := os.Getenv("JWT_SECRET_KEY")

	signedDetails := &SignedDetails{
		Id:    id.String(),
		Email: email,
		Name:  name,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    jwt_secret_key,
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, signedDetails)

	return token.SignedString([]byte(jwt_secret_key))

}

func ValidateToken(tokenString string) (*SignedDetails, error) {

	jwt_secret_key := os.Getenv("JWT_SECRET_KEY")

	token, err := jwt.ParseWithClaims(tokenString, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwt_secret_key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SignedDetails); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")

}
