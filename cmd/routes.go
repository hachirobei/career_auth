package main

import (
    "career.com/auth/controllers"
    "career.com/auth/middleware"  
    "career.com/auth/models"    
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    roleController := controllers.NewRoleController()

    // Public routes
    router.POST("/login", controllers.Login)
    router.POST("/signup", controllers.SignUp)

    // Admin routes
    adminRoutes := router.Group("/admin").Use(middleware.RoleCheckMiddleware(models.RoleAdmin))
    {
        adminRoutes.POST("/users", roleController.CreateUsers)
        adminRoutes.PUT("/users/:id", roleController.UpdateUsers)
        adminRoutes.DELETE("/users/:id", roleController.DeleteUsers)

        adminRoutes.POST("/roles", roleController.CreateRole)
        adminRoutes.PUT("/roles/:id", roleController.UpdateRole)
        adminRoutes.DELETE("/roles/:id", roleController.DeleteRole)
    }

    // User routes
    userRoutes := router.Group("/user").Use(middleware.RoleCheckMiddleware(models.RoleUser))
    {

    }

    return router
}