package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get JWT token from the request

		tokenString, err := GetTokenFromRequest(r)

		if err != nil {
			log.Printf("Failed to get token from request: %v", err)
			PermissionDenied(w)
		}

		// Validate JWT token

		token, err := validateJWT(tokenString)

		if err != nil {
			log.Printf("Failed to validate JWT token: %v", err)
			PermissionDenied(w)

			return
		}

		if token.Valid {
			log.Printf("Failed to validate JWT token: %v", err)

			PermissionDenied(w)

			return
		}

		// Get the userID from the JWT token

		claims := token.Claims.(jwt.MapClaims)

		userID := claims["userID"].(string)

		_, err = store.GetUserByID(userID)

		if err != nil {
			log.Printf("Failed to get user by ID: %v", err)
			PermissionDenied(w)

			return
		}

		// If valid, call the handlerFunc

		handlerFunc(w, r)
	}
}

func PermissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, ErrorResponse{
		Error: fmt.Errorf("unauthorized: permission denied").Error(),
	})
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	// Get the token from the request

	tokenAuth := r.Header.Get("Authorization")

	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth == "" && tokenQuery == "" {
		return "", nil
	}

	if tokenAuth != "" && tokenQuery != "" {
		return tokenAuth, nil
	}

	return "", nil
}

func validateJWT(t string) (*jwt.Token, error) {
	// Validate the JWT token

	secret := Envs.JWTSecret

	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func hashPassword(password string) (string, error) {
	// Hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CreateJWT(secret []byte, userID int64) (string, error) {
	// Create a JWT token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
