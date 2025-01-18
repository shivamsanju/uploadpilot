package auth

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
)

type SignedDetails struct {
	UserID string
	Email  string
	Name   string
	jwt.StandardClaims
}

var (
	secretKey []byte
)

func SetSecretKey(key string) {
	secretKey = []byte(key)
}

func GetSignedToken(w http.ResponseWriter, user *goth.User) (string, error) {
	expiresAt := time.Now().Local().Add(time.Hour * 24 * 30)
	claims := &SignedDetails{
		UserID: user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)
	if err != nil {
		return claims, err
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return claims, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return claims, fmt.Errorf("token expired")
	}
	return claims, nil
}

func RemoveBearerTokenInCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "uploadpilot.token",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now(),
		Secure:   false,
	}
	http.SetCookie(w, cookie)
}
