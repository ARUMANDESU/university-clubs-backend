package jwt

import (
	"fmt"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"os"
	"strings"
	"time"
)

func GetUserID(tokenString, secret string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret used to sign the token
		return []byte(secret), nil
	})

	if err != nil {
		errorMessage := err.Error()
		if strings.Contains(errorMessage, "token is expired") {
			return 0, domain.ErrTokenIsExpired
		}
		return 0, err
	}

	// Check if the token is valid
	if !token.Valid {
		return 0, domain.ErrTokenIsNotValid
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, domain.ErrInvalidTokenClaims
	}

	// Extract user_id from claims
	expTime, ok := claims["exp"].(float64)
	slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).Debug(fmt.Sprintf("%v", expTime))
	if !ok {
		return 0, domain.ErrInvalidTokenClaims
	}

	if time.Now().Unix() > int64(expTime) {
		return 0, domain.ErrTokenIsExpired
	}

	// Extract user_id from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, domain.ErrUserIDClaimNotFound
	}

	return int64(userID), nil
}
