package controllers

import (
    "net/http"
    "career.com/auth/database"
    "career.com/auth/models"
	"career.com/auth/utils"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"

	"strconv"
)

type UserController struct{}

// NewUserController creates a new UserController instance
func NewUserController() *UserController {
	return &UserController{}
}

// CreateUser handles the creation of a new user
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	db := database.GetDB()
	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated,user)
}

// UpdateUser handles updating a user's information
func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.Users
	userID := c.Param("id")

	db := database.GetDB()
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&user)
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles the deletion of a user
func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	db := database.GetDB()
	
	if result := db.Delete(&models.Users{}, userID); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}	



// ForgotPassword handles user password recovery process
func (uc *UserController) ForgotPassword(c *gin.Context) {
	// Placeholder for password recovery logic
	c.JSON(http.StatusOK, gin.H{"message": "Password recovery not implemented"})
	}
	
	// GetDB is a placeholder for your method to get the DB instance
	func GetDB() *gorm.DB {
	// Implement this based on your database setup
	return nil
}

// Index users handles pagination and listing of users
func (uc *UserController) Index(c *gin.Context) {
    pageStr, pageSizeStr := c.DefaultQuery("page", "1"), c.DefaultQuery("pageSize", "10")
    page, err := strconv.Atoi(pageStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }
    pageSize, err := strconv.Atoi(pageSizeStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
        return
    }

    var users []models.Users
    var totalRecords int64
    offset := (page - 1) * pageSize

    db := database.GetDB()
    // Count total records
    db.Model(&models.Users{}).Count(&totalRecords)

    result := db.Offset(offset).Limit(pageSize).Find(&users)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize) // Calculate total pages

    c.JSON(http.StatusOK, gin.H{
        "data":         users,
			"pagination": gin.H{
			"totalRecords": totalRecords,
			"totalPages":   totalPages,
			"currentPage":  page,
		},
    })
}


// Index users handles pagination and listing of users
func (uc *UserController) Index(c *gin.Context) {
    pageStr, pageSizeStr := c.DefaultQuery("page", "1"), c.DefaultQuery("pageSize", "10")
    page, err := strconv.Atoi(pageStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }
    pageSize, err := strconv.Atoi(pageSizeStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
        return
    }

    var users []models.Users
    var totalRecords int64
    offset := (page - 1) * pageSize

    db := database.GetDB()
    // Count total records
    db.Model(&models.Users{}).Count(&totalRecords)

	if result.Error != nil {
        response := utils.NewApiResponse(http.StatusInternalServerError, false, nil, "Error finding users")
        c.JSON(http.StatusInternalServerError, response)
        return
    }

    totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize) // Calculate total pages

	if len(users) == 0 {
        response := utils.NewApiResponse(http.StatusOK, true, nil, "No record found")
        c.JSON(http.StatusOK, response)
        return
    }

	response := utils.NewApiResponse(http.StatusOK, true, gin.H{
        "users":        users,
        "totalRecords": totalRecords,
        "totalPages":   totalPages,
        "currentPage":  page,
    }, "")
    c.JSON(http.StatusOK, response)
}

// Show fetches a single user by ID
func (uc *UserController) Show(c *gin.Context) {
    var user models.Users
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        response := utils.NewApiResponse(http.StatusBadRequest,  false, nil, "error": "Invalid user ID")
		c.JSON(http.StatusInternalServerError, response)
        return
    }

	if len(users) == 0 {
        response := utils.NewApiResponse(http.StatusOK, true, nil, "No record found")
        c.JSON(http.StatusOK, response)
        return
    }

	response := utils.NewApiResponse(http.StatusOK, true,  gin.H{"record": user})
    c.JSON(http.StatusOK, response)
}

func RefreshToken(c *gin.Context) {
    var request struct {
        RefreshToken string `json:"refreshToken"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var userToken models.UserToken
    db := database.ConnectToDB() // Replace with your method to get DB instance
    result := db.Where("token = ?", request.RefreshToken).First(&userToken)

    if result.Error != nil || userToken.ExpiresAt.Before(time.Now()) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
        return
    }

    // Generate a new access token
    newAccessToken, err := GenerateJWT(user.Username, user.Role)
    if err != nil {
		response := utils.NewApiResponse(http.StatusBadRequest,  false, nil,"error": "Error generating access token")
        c.JSON(http.StatusInternalServerError,response)
        return
    }

	response := utils.NewApiResponse(http.StatusOK, true,  gin.H{"accessToken": newAccessToken})
    c.JSON(http.StatusOK, response)
	return
}