package utils

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(claims jwt.MapClaims, privKeyPEM []byte) (string, error) {
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyPEM)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(tokenString string, pubKeyPEM []byte) (jwt.MapClaims, error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyPEM)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
