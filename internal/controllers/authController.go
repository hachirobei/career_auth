package controllers

import (
    "fmt"
    "net/http"
    "os"
    "time"

    "career.com/auth/internal/database"
    "career.com/auth/internal/models"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
)

var secretKey = os.Getenv("SECRET_KEY")

type CustomClaims struct {
    Username string `json:"username"`
    Role     models.Role `json:"role"`
    jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func GenerateJWT(username string, role models.Role) (string, error) {
    claims := CustomClaims{
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func getUserByUsername(username string) (models.User, error) {
    db := database.ConnectToDB()
    var user models.User
    result := db.Where("username = ?", username).First(&user)
    if result.Error != nil {
        return models.User{}, result.Error
    }
    return user, nil
}

func saveRefreshToken(userID int, refreshToken string) error {
    // Implement logic to save the refresh token to the database
    return nil
}

func GenerateRefreshToken(username string) (string, error) {
    // Implement logic to generate a refresh token
    return "refresh-token", nil
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

    if err := c.ShouldBindJSON(&loginInfo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := getUserByUsername(loginInfo.Username)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or
        password"})
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
    
    err = saveRefreshToken(user.ID, refreshToken)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving refresh token"})
        return
    }
        
    c.JSON(http.StatusOK, gin.H{
        "accessToken":  accessToken,
        "refreshToken": refreshToken,
    })
    
}


// TokenValidationHandler checks if the provided JWT token is valid
func TokenValidationHandler(c *gin.Context) {
    tokenString := c.GetHeader("Authorization")if tokenString == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
        return
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(secretKey), nil
    })

    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "message": "Token is valid"})
}