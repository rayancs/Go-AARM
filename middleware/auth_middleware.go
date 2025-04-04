package middleware

import (
	"app/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Mock function to validate JWT token

// JWT Middleware
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		// Extract token (assuming "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := services.VerifyToken(token)
		if err != nil || claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		// Token is valid, continue
		c.Next()
	}
}

// example use case with a protected end point
// func main() {
// 	r := gin.Default()

// 	// Apply JWT middleware to protected routes
// 	protected := r.Group("/protected")
// 	protected.Use(JWTMiddleware())
// 	{
// 		protected.GET("/data", func(c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{"message": "You have access!"})
// 		})
// 	}
// 	r.Run(":8080") // Start the server on port 8080
// }
