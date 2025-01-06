package controllers

import (
    "auth-backend/utils"
    "net/http"
    "time"
    "log"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
    Nim      string `json:"nim" binding:"required"`
    Password string `json:"password" binding:"required"`
    Name     string `json:"name" binding:"required"`
    Prodi    string `json:"prodi" binding:"required"`
    Mentor   string `json:"mentor" binding:"required"`
}

type User struct {
    Nim   string `json:"nim"`
    Name  string `json:"name"`
    Prodi string `json:"prodi"`
    Mentor string `json:"mentor"`
}


type LoginInput struct {
    Nim      string `json:"nim" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// Register handles user registration
func Register(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    userData := map[string]interface{}{
        "nim":           input.Nim,
        "password_hash": string(hashedPassword),
        "name":          input.Name,
        "prodi":         input.Prodi,
        "mentor":        input.Mentor,
    }

    var result map[string]interface{}
    err = utils.SupabaseClient.DB.From("users").Insert(userData).Execute(&result)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// Login handles user authentication
func Login(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    var users []map[string]interface{}
    err := utils.SupabaseClient.DB.From("users").
        Select("*").
        Eq("nim", input.Nim).
        Execute(&users)

    if err != nil || len(users) == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "NIM not found"})
        return
    }

    passwordHash, ok := users[0]["password_hash"].(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
        return
    }

    userID, ok := users[0]["id"].(string)
    userName, _ := users[0]["name"].(string) // Ambil nama jika ada
    userNim, _ := users[0]["nim"].(string)  // Ambil nim jika ada
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
        return
    }

    token, err := utils.GenerateToken(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    utils.SetTokenCookie(c, token)

    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "user": gin.H{
            "id":   userID,
            "name": userName,
            "nim":  userNim,
        },
    })
}


// Logout handles user logout by clearing the JWT cookie
func Logout(c *gin.Context) {
    cookie := &http.Cookie{
        Name:     "jwt",
        Value:    "",
        Path:     "/",
        Expires:  time.Now().Add(-time.Hour), // Expire the cookie
        HttpOnly: true,
        SameSite: http.SameSiteNoneMode, 
        Secure:   false, // Change to true in production
    }
    http.SetCookie(c.Writer, cookie)

    c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func GetUser(c *gin.Context) {
    userID := c.GetString("user_id") // Ambil user_id dari middleware
    log.Printf("Authenticated user ID: %s", userID) // Log untuk debugging
    
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Buat slice untuk menampung hasil query
    var users []User

    // Jalankan query untuk mendapatkan pengguna berdasarkan ID
    err := utils.SupabaseClient.DB.From("users").
        Select("nim,name,prodi,mentor").
        Eq("id", userID).
        Execute(&users) // Gunakan Execute untuk menempatkan hasil ke slice
    if err != nil || len(users) == 0 {
        log.Printf("Error fetching user: %v", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Ambil pengguna pertama dari hasil query
    user := users[0]

    c.JSON(http.StatusOK, gin.H{"user": user})
}


