package main

import (
    "os"
    "career.com/auth/controllers" 
    "career.com/auth/database" 

    "github.com/gin-gonic/gin"
)

func main() {
    db := database.Connect() // Make sure to implement this
    db.AutoMigrate(&models.Users{}, &models.Role{},  &models.UserToken{})

    router := SetupRouter() // Set up routes from routes.go

    // Get port from environment variable with a fallback to a default
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port if not specified
    }

    // Start the server on the specified port
    router.Run(":" + port)
}