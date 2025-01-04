package controllers

import (
	"auth-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Nim      string `json:"nim" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Nim      string `json:"nim" binding:"required"`
	Password string `json:"password" binding:"required"`
}
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
    }

    var result interface{}
    err = utils.SupabaseClient.DB.From("users").Insert(userData).Execute(&result)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
        return
    }

    token, err := utils.GenerateToken(input.Nim)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}


func Login(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

    passwordHash := users[0]["password_hash"].(string)
    if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
        return
    }

    token, err := utils.GenerateToken(input.Nim)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}
