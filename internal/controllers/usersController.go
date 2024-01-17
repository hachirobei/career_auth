package controllers

import (
    "net/http"
    "career.com/auth/internal/database"
    "career.com/auth/internal/models"
	"career.com/auth/internal/helpers"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"

	"strconv"
)

type UserController struct{}

// NewUserController creates a new UserController instance
func NewUserController() *UserController {
	return &UserController{}
}

// BindUserData binds JSON data to a user model
func BindUserData(c *gin.Context, user *models.Users) error {
    if err := c.ShouldBindJSON(user); err != nil {
        return err
    }
    return nil
}

// HashUserPassword hashes a user's password
func HashUserPassword(user *models.Users) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
    return nil
}

// CreateUser handles the creation of a new user
func (uc *UserController) CreateUser(c *gin.Context) {
    var user models.Users
    if err := BindUserData(c, &user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := HashUserPassword(&user); err != nil {
        response := helpers.NewApiResponse(http.StatusInternalServerError, false, nil, "Error hashing password")
        c.JSON(http.StatusInternalServerError, response)
        return
    }

    db := database.GetDB()
    if result := db.Create(&user); result.Error != nil {
        response := helpers.NewApiResponse(http.StatusInternalServerError, false, nil, result.Error.Error())
        c.JSON(http.StatusInternalServerError, response)
        return
    }

    c.JSON(http.StatusCreated, user)
}

// UpdateUser handles updating a user's information
func (uc *UserController) UpdateUser(c *gin.Context) {
    var updatedUserData models.Users
    userID := c.Param("id")

    db := database.GetDB()
    var existingUser models.Users
    if err := db.Where("id = ?", userID).First(&existingUser).Error; err != nil {
        response := helpers.NewApiResponse(http.StatusNotFound, false, nil, "User not found")
        c.JSON(http.StatusNotFound, response)
        return
    }

    if err := BindUserData(c, &updatedUserData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash the password if it's being updated
    if updatedUserData.Password != "" {
        if err := HashUserPassword(&updatedUserData); err != nil {
            response := helpers.NewApiResponse(http.StatusInternalServerError, false, nil, "Error hashing password")
            c.JSON(http.StatusInternalServerError, response)
            return
        }
        existingUser.Password = updatedUserData.Password
    }

    // Update other user fields as necessary
    existingUser.FullName = updatedUserData.FullName
    existingUser.Email = updatedUserData.Email
    existingUser.Phone = updatedUserData.Phone
    existingUser.Username = updatedUserData.Username
    // ... handle other fields

    db.Save(&existingUser)
    c.JSON(http.StatusOK, existingUser)
}

// Check if email exists
var user models.Users
result := db.Where("email = ?", request.Email).First(&user)
if result.Error != nil {
    c.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent if the email is registered"})
    return
}

func generateSecureToken() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}

func updateUserWithResetToken(db *gorm.DB, user *models.Users, token string) error {
    user.ResetPasswordToken = token
    user.ResetPasswordExpires = time.Now().Add(24 * time.Hour) // token expires in 24 hours
    return db.Save(user).Error
}

func (uc *UserController) ForgotPassword(c *gin.Context) {
    var request struct {
        Email string `json:"email"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Verify Email
    db := database.GetDB() // Make sure you have a way to get your database instance
    var user models.Users
    result := db.Where("email = ?", request.Email).First(&user)
    if result.Error != nil {
        c.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent if the email is registered"})
        return
    }

    // Generate Token
    token, err := generateSecureToken()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating reset token"})
        return
    }

    // Update User with Token
    err = updateUserWithResetToken(db, &user, token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user with reset token"})
        return
    }

    // Send Password Reset Email using helper
    err := helpers.SendPasswordResetEmail(user.Email, token)
    if err != nil {
        log.Println("Failed to send password reset email:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send password reset email"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent if the email is registered"})
}


// DeleteUser handles the deletion of a user
func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	db := database.GetDB()
	
	if result := db.Delete(&models.Users{}, userID); result.Error != nil {
        response := helpers.NewApiResponse(http.StatusInternalServerError, false, nil, result.Error.Error())
        c.JSON(http.StatusInternalServerError, response)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

    response := helpers.NewApiResponse(http.StatusOK, true, result, "User deleted successfully")
    c.JSON(http.StatusOK, response)
}	

// Index users handles pagination and listing of users
func (uc *UserController) Index(c *gin.Context) {
    db := database.GetDB()
    var users []models.Users
    var totalRecords int64

    db.Model(&models.Users{}).Count(&totalRecords)

    pagination, err := helpers.NewPagination(c, totalRecords)
    if err != nil {
        response := helpers.NewApiResponse(http.StatusBadRequest, false, nil, "Invalid pagination parameters")
        c.JSON(http.StatusBadRequest, response)
        return
    }

    offset := (pagination.CurrentPage - 1) * pagination.PageSize
    result := db.Offset(offset).Limit(pagination.PageSize).Find(&users)
    if result.Error != nil {
        response := helpers.NewApiResponse(http.StatusInternalServerError, false, nil, result.Error.Error())
        c.JSON(http.StatusInternalServerError, response)
        return
    }

    responseData := gin.H{
        "users":       users,
        "pagination": pagination,
    }
    response := helpers.NewApiResponse(http.StatusOK, true, responseData, "Users fetched successfully")
    c.JSON(http.StatusOK, response)
}

// Show fetches a single user by ID
func (uc *UserController) Show(c *gin.Context) {
    var user models.Users
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        response := helpers.NewApiResponse(http.StatusBadRequest,  false, nil, "error": "Invalid user ID")
		c.JSON(http.StatusInternalServerError, response)
        return
    }

	if len(users) == 0 {
        response := helpers.NewApiResponse(http.StatusOK, true, nil, "No record found")
        c.JSON(http.StatusOK, response)
        return
    }

	response := helpers.NewApiResponse(http.StatusOK, true,  gin.H{"record": user},  "Users fetched successfully" )
    c.JSON(http.StatusOK, response)
}

func RefreshToken(c *gin.Context) {
    var request struct {
        RefreshToken string `json:"refreshToken"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        response := helpers.NewApiResponse(http.StatusBadRequest, false, nil, gin.H{"error": err.Error()})
        c.JSON(http.StatusBadRequest, response)
        return
    }

    var userToken models.UserToken
    db := database.ConnectToDB() // Replace with your method to get DB instance
    result := db.Where("token = ?", request.RefreshToken).First(&userToken)

    if result.Error != nil || userToken.ExpiresAt.Before(time.Now()) {
        response := helpers.NewApiResponse(http.StatusUnauthorized.StatusOK, true, nil, gin.H{"error": "Invalid or expired refresh token"})
        c.JSON(http.StatusUnauthorized,response)
        return
    }

    // Generate a new access token
    newAccessToken, err := GenerateJWT(user.Username, user.Role)
    if err != nil {
		response := helpers.NewApiResponse(http.StatusBadRequest,  false, nil,"error": "Error generating access token")
        c.JSON(http.StatusInternalServerError,response)
        return
    }

	response := helpers.NewApiResponse(http.StatusOK, true,  gin.H{"accessToken": newAccessToken} , "User refrech token successfully")
    c.JSON(http.StatusOK, response)
	return
}