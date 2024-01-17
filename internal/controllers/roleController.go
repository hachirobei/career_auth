package controllers

import (
    "net/http"
    "career.com/auth/internal/database"
    "career.com/auth/internal/models"  
	
    "github.com/gin-gonic/gin"
)

type RoleController struct{}

func NewRoleController() *RoleController {
    return &RoleController{}
}

// CreateRole handles the creation of a new role
func (rc *RoleController) CreateRole(c *gin.Context) {
    var role models.Role
    if err := c.ShouldBindJSON(&role); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db := database.GetDB()
    if result := db.Create(&role); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusCreated, role)
}

// UpdateRole handles updating a role's information
func (rc *RoleController) UpdateRole(c *gin.Context) {
    var role models.Role
    roleID := c.Param("id")

    db := database.GetDB()
    if err := db.Where("id = ?", roleID).First(&role).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
        return
    }

    if err := c.ShouldBindJSON(&role); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.Save(&role)
    c.JSON(http.StatusOK, role)
}

// DeleteRole handles the deletion of a role
func (rc *RoleController) DeleteRole(c *gin.Context) {
    roleID := c.Param("id")
    db := database.GetDB()

    if result := db.Delete(&models.Role{}, roleID); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}