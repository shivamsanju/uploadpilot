package auth

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
	"github.com/uploadpilot/uploadpilot/internal/msg"
)

var (
	secretKey []byte
)

type Claims struct {
	UserID string
	Email  string
	Name   string
	jwt.StandardClaims
}

func GenerateToken(w http.ResponseWriter, user *goth.User, expiry time.Duration) (string, error) {
	if secretKey == nil {
		return "", errors.New(msg.JWTSecretKeyNotSet)
	}
	claims := &Claims{
		UserID: user.UserID,
		Email:  user.Email,
		Name:   user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(expiry).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateToken(jwtToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return claims, errors.New(msg.InvalidToken)
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return claims, errors.New(msg.TokenExpired)
	}
	return claims, nil
}
