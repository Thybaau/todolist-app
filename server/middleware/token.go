package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("secretpassword")

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token valid for 24 hours

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Validate the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get JWT from header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			NewHTTPError(w, "Missing authorization token", http.StatusUnauthorized, nil)
			return
		}

		// Check that token is prefixed by "Bearer "
		tokenSplitted := strings.Split(tokenString, " ")
		if len(tokenSplitted) != 2 || tokenSplitted[0] != "Bearer" {
			NewHTTPError(w, "Invalid authentication token, no prefix Bearer", http.StatusUnauthorized, nil)
			return
		}

		claims, err := VerifyToken(tokenSplitted[1])
		if err != nil {
			NewHTTPError(w, "Invalid authentication token", http.StatusUnauthorized, err)
			return
		}

		username, ok := claims["username"]
		if !ok {
			NewHTTPError(w, "Username not found in token payload", http.StatusUnauthorized, nil)
			return
		}

		ctx := context.WithValue(r.Context(), "username", username)
		r = r.WithContext((ctx))
		next.ServeHTTP(w, r)
	})
}
