package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			// If no token, return 401 Unauthorized
			c.JSON(401, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
		// Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Return the secret key for validation
			return jwtSecret, nil
		})
		// If token is invalid or parsing fails, return 401 Unauthorized
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Parse JWT claims and store in Gin Context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if email, ok := claims["email"]; ok {
				c.Set("email", email)
			}
			if role, ok := claims["role"]; ok {
				c.Set("role", role)
			}
			if name, ok := claims["name"]; ok {
				c.Set("name", name)
			}
		}
		// Token is valid, continue to next handler
		c.Next()
	}
}
