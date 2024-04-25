package jwt

import (
	"fmt"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"os"
	"time"
)

func GetUserID(tokenString, secret string) (int64, bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret used to sign the token
		return []byte(secret), nil
	})

	if err != nil {
		return 0, true, err
	}

	// Check if the token is valid
	if !token.Valid {
		return 0, true, domain.ErrTokenIsNotValid
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, true, domain.ErrInvalidTokenClaims
	}

	// Extract user_id from claims
	expTime, ok := claims["exp"].(float64)
	slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).Debug(fmt.Sprintf("%v", expTime))
	if !ok {
		return 0, true, domain.ErrInvalidTokenClaims
	}

	if time.Now().Unix() > int64(expTime) {
		return 0, true, nil
	}

	// Extract user_id from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, true, domain.ErrUserIDClaimNotFound
	}

	return int64(userID), false, nil
}
