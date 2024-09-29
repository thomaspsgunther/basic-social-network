package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Secret key being used to sign tokens
var (
	SecretKey = []byte(os.Getenv("TOKEN_KEY"))
)

// GenerateToken generates a jwt token and assign a id to its claims and return it
func GenerateToken(id uuid.UUID) (string, error) {
	idStr := id.String()
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["id"] = idStr
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenStr, err := token.SignedString(SecretKey)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return tokenStr, nil
}

// ParseToken parses a jwt token and returns the id in its claims
func ParseToken(tokenStr string) (uuid.UUID, error) {
	tokenStr = tokenStr[len("Bearer "):]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idStr := claims["id"].(string)
		id, err := uuid.Parse(idStr)
		if err != nil {
			return uuid.UUID{}, err
		}

		return id, nil
	} else {
		return uuid.UUID{}, err
	}
}
