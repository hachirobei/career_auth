package middleware

import (
    "net/http"
    "your_project_name/models" // Replace with your models package

    "github.com/gin-gonic/gin"
)

// RoleCheckMiddleware checks the user's role and allows access based on the role
func RoleCheckMiddleware(requiredRole models.Role) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Here you need to retrieve the user's role. This is just an example.
        // Replace it with your actual logic to retrieve the user's role.
        userRole, exists := c.Get("userRole")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }

        if userRole != requiredRole {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
            return
        }

        c.Next()
    }
}