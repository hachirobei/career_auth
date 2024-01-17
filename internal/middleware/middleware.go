package middleware

import (
    "fmt"
    "net/http"
    "strings"
    "career.com/auth/models" 

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

// This should be a secret key unique to your application
var jwtKey = []byte("your_secret_key") // Store this securely

// JWTAuthentication is a middleware function for validating JWT tokens
func JWTAuthentication() gin.HandlerFunc {
    return func(c *gin.Context) {
			// Get the JWT token from the 'Authorization' header
			tokenString := c.GetHeader("Authorization")
			if tokenString == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
				return
			}

			// The token should have a 'Bearer ' prefix
			if !strings.HasPrefix(tokenString, "Bearer ") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
				return
			}

			// Strip 'Bearer ' from the token string
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			// Parse the JWT token
			token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Make sure token's signature algorithm is what you expect
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtKey, nil
			})

			// Handle any errors from parsing
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Signature invalid"})
					return
				}
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
				return
			}

			// Check if the token is valid
			if !token.Valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
				return
			}

			// If the token is valid, get the user's role from the claims
			if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
				userRole
			// Set the user information in the context
			c.Set("userID", claims.ID)
			c.Set("userRole", userRole)

			c.Next() // proceed to the next handler
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			return
		}
	}
}