package middleware

import (
    "net/http"
    "career.com/auth/internal/models"

    "github.com/gin-gonic/gin"
)

package your_package_name

import (
    "net/http"
    "career.com/auth/internal/models"    
    "career.com/auth/internal/database" // replace with the actual path to your database package
    "github.com/gin-gonic/gin"
)

func RoleCheckMiddleware(requiredRoles ...int) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Retrieve the user's ID from the context/session
        userID, exists := c.Get("userID")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }

        // Query the AccessRole table to find the user's roles
        var accessRoles []models.AccessRole
        result := database.GetDB().Where("user_id = ?", userID).Find(&accessRoles)
        if result.Error != nil || len(accessRoles) == 0 {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }

        // Check if the user has any of the required roles
        hasRequiredRole := false
        for _, accessRole := range accessRoles {
            for _, requiredRole := range requiredRoles {
                if int(accessRole.RoleID) == requiredRole {
                    hasRequiredRole = true
                    break
                }
            }
            if hasRequiredRole {
                break
            }
        }

        if !hasRequiredRole {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
            return
        }

        // If user has required role, proceed with the request
        c.Next()
    }
}