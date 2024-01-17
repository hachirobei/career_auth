package main

import (
    "os"
    "career.com/auth/internal/controllers" 
    "career.com/auth/internal/database" 
    "career.com/auth/internal/models" 

    "github.com/gin-gonic/gin"
)

func main() {
    db := database.Connect() // Make sure to implement this
    db.AutoMigrate(&models.Users{}, &models.Role{},  &models.UserToken{})

    func createRoleIfNotExists(db *gorm.DB) {
        var count int64
        db.Model(&models.Role{}).Where("description = ?", "admin").Count(&count)
        if count == 0 {
            roleInfo := models.Users{
                Description: "admin"
                Status:   1,
            }
            result := db.Create(&roleInfo)
            if result.Error != nil {
                fmt.Println("Failed to create role admin:", result.Error)
                return
            }
            RoleId :=  result->id
            fmt.Println("Role admin created successfully")
        }
    }

    func createAdminUserIfNotExists(db *gorm.DB) {
        var count int64
        db.Model(&models.Users{}).Where("username = ?", "admin").Count(&count)
        if count == 0 {
            adminUser := models.Users{
                FullName: "Admin User",
                Email:    "admin@example.com",
                Phone:    "1234567890",
                Username: "admin",
                Password: "adminpassword", // You should hash the password in a real application
                RoleId:   1, // Assuming 1 is the admin role ID
                Status:   1,
            }
            result := db.Create(&adminUser)
            if result.Error != nil {
                fmt.Println("Failed to create admin user:", result.Error)
                return
            }
            UserId :=  result->id
            fmt.Println("Admin user created successfully")
        }
    }

    func createAccesssRoleIfNotExists(db *gorm.DB) {
        var count int64
        db.Model(&models.AccessRole{}).Where("userid = ?",UserId).Where("roleid = ?",RoleId).Count(&count)
        if count == 0 {
            AccessRoleUser := models.AccessRole{
                FullName: "Admin User",
                Email:    "admin@example.com",
                Phone:    "1234567890",
                Username: "admin",
                Password: "adminpassword", // You should hash the password in a real application
                Status:   1,
            }
            result := db.Create(&AccessRoleUser)
            if result.Error != nil {
                fmt.Println("Failed to create access role admin:", result.Error)
                return
            }
            fmt.Println("Access role for admin created successfully")
        }
    }

    router := SetupRouter() // Set up routes from routes.go

    // Get port from environment variable with a fallback to a default
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port if not specified
    }

    // Start the server on the specified port
    router.Run(":" + port)
}