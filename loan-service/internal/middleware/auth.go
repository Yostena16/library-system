package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware protects routes by requiring a valid JWT token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Read the "Authorization" header
		authHeader := c.GetHeader("Authorization") //from the request
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// 2. It should look like "Bearer <token>" — split it
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "format must be 'Bearer <token>'"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 3. Parse + verify the token using our secret
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// make sure it was signed with HMAC (HS256), not something else
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// 4. Pull the claims (the data) out of the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// 5. member_id comes back as float64 (JSON numbers) — convert to uint
		memberIDFloat, ok := claims["member_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			c.Abort()
			return
		}

		// 6. Save who's logged in, so handlers can use it
		c.Set("member_id", uint(memberIDFloat))

		role, _ := claims["role"].(string)
		c.Set("role", role)

		// 7. All good — let the request continue
		c.Next()
	}
}
