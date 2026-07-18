package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a signed JWT for a given member ID.
func GenerateToken(memberID uint) (string, error) {
	// How long until the token expires (default 24h)
	hours, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil || hours == 0 {
		hours = 24
	}

	// "Claims" = the data we store inside the token
	claims := jwt.MapClaims{
		"member_id": memberID,
		"exp":       time.Now().Add(time.Duration(hours) * time.Hour).Unix(),
	}

	// Create the token, then sign it with our secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
