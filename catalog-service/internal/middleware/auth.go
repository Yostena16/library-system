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
// Catalog only ever verifies tokens issued by loan-service — it never
// issues its own, so this is a straight validator.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Read the "Authorization" header
		authHeader := c.GetHeader("Authorization")
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

		// 3. Parse + verify the token using our shared secret
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
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

		// 4. Pull the claims out of the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// 5. Catalog only cares about role (for the librarian check) —
		// member_id isn't needed here, but we set it too in case a
		// handler ever wants to know who made the request.
		if memberIDFloat, ok := claims["member_id"].(float64); ok {
			c.Set("member_id", uint(memberIDFloat))
		}

		role, _ := claims["role"].(string)
		c.Set("role", role)

		// 6. All good — let the request continue
		c.Next()
	}
}
