package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a signed JWT for a given member ID.
func GenerateToken(memberID uint, role string) (string, error) {
	hours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil || hours == 0 {
		hours = 24
	}

	claims := jwt.MapClaims{
		"member_id": memberID,
		"role":      role,
		"exp":       time.Now().Add(time.Duration(hours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
