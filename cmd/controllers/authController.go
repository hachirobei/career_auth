package controllers

import (
    "net/http"
    "time"
    "os"

    "career.com/auth/database"
    "career.com/auth/models"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
)

var secretKey = os.Getenv("SECRET_KEY"); // should be in an environment variable

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func GenerateJWT(username string, role models.Role) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "role":     role,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func SignUp(c *gin.Context) {
    var newUser models.User
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := HashPassword(newUser.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }
    newUser.Password = hashedPassword

    db := database.ConnectToDB()
    result := db.Create(&newUser)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
    var loginInfo struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    var user models.User

    if err := c.ShouldBindJSON(&loginInfo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db := database.ConnectToDB() // Ensure this is your function to connect to the database
    result := db.Where("username = ?", loginInfo.Username).First(&user)
    if result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
        return
    }

    if !CheckPasswordHash(loginInfo.Password, user.Password) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
        return
    }

    accessToken, err := GenerateJWT(user.Username, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
        return
    }

    refreshToken, err := GenerateRefreshToken(user.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
        return
    }

    // Save the refresh token in user_token table
    userToken := models.UserToken{
        UserID:    user.ID,
        Token:     refreshToken,
        ExpiresAt: time.Now().Add(72 * time.Hour), // Set expiration for the refresh token
    }
    
    db.Save(&userToken)

    c.JSON(http.StatusOK, gin.H{
        "accessToken":  accessToken,
        "refreshToken": refreshToken,
    })
}