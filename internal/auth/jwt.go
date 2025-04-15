package auth

import (
	"log"
	"fmt"
	"github.com/google/uuid"
	"golang-jwt/jwt/v5"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now()) + expiresIn,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}

	//make a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//sign token
	signedToken := token.SignedString(tokenSecret)



	return 
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	return
}