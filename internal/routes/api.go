package main

import (
    "career.com/auth/internal/controllers"
    "career.com/auth/internal/middleware"  
    "career.com/auth/internal/models"    
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    roleController := controllers.NewRoleController()
    userController := controllers.NewUserController()
    accessRoleController := controllers.NewAccessRoleController()

    // Public routes
    router.POST("/login", controllers.Login)
    router.POST("/signup", controllers.SignUp)

    // Admin routes
    adminRoutes := router.Group("/admin").Use(middleware.RoleCheckMiddleware(models.RoleAdmin))
    {
        adminRoutes.GET("/users/:id", userController.ShowUsers)
        adminRoutes.GET("/users", userController.IndexUsers)
        adminRoutes.POST("/users", userController.CreateUsers)
        adminRoutes.PUT("/users/:id", userController.UpdateUsers)
        adminRoutes.DELETE("/users/:id", userController.DeleteUsers)

        adminRoutes.GET("/roles/:id", roleController.ShowRole)
        adminRoutes.GET("/roles", roleController.IndexRole)
        adminRoutes.POST("/roles", roleController.CreateRole)
        adminRoutes.PUT("/roles/:id", roleController.UpdateRole)
        adminRoutes.DELETE("/roles/:id", roleController.DeleteRole)


        accessRoleRoutes.GET("/accessrole/:id", accessRoleController.ShowRole)
        accessRoleRoutes.GET("/accessrole", accessRoleController.IndexRole)
        accessRoleRoutes.GET("/accessrole/user/:userId", accessRoleController.IndexRole)
        accessRoleRoutes.POST("/accessrole", accessRoleController.CreateRole)
        accessRoleRoutes.PUT("/accessrole/:id", accessRoleController.UpdateRole)
        accessRoleRoutes.DELETE("/accessrole/:id", roleController.DeleteRole)
    }

    // User routes
    userRoutes := router.Group("/user").Use(middleware.RoleCheckMiddleware(models.RoleUser))
    {

    }

    router.GET("/validate-token", controllers.TokenValidationHandler)

    return router
}