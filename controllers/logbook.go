package controllers

import (
    "auth-backend/utils"
    "net/http"
	"log"

    "github.com/gin-gonic/gin"
)

type Logbook struct {
    UserID    string `json:"user_id"`
    StartDate string `json:"start_date" binding:"required"` // Start Date
    EndDate   string `json:"end_date" binding:"required"`   // End Date
    Activity  string `json:"activity" binding:"required"`   // Activity
    PIC       string `json:"pic" binding:"required"`        // Person In Charge
    Status    string `json:"status" binding:"required"`     // Status: approve, pending, reject
}

// GetLogbooks retrieves logbooks for the authenticated user
func GetLogbooks(c *gin.Context) {
    userID := c.GetString("user_id") // Retrieved from middleware
    var logbooks []Logbook

    err := utils.SupabaseClient.DB.From("logbooks").
        Select("*").
        Eq("user_id", userID).
        Execute(&logbooks)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logbooks"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": logbooks})
}

// AddLogbook adds a new logbook entry for the authenticated user
func AddLogbook(c *gin.Context) {
    userID := c.GetString("user_id") // Retrieved from middleware

    // Validate userID
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing user_id"})
        return
    }
    log.Printf("Adding logbook for user ID: %s", userID) // Debug log

    var input Logbook
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Assign userID to the logbook entry
    input.UserID = userID

    // Log all data being inserted into the database
    log.Printf("Logbook Input: %+v", input)

    // Insert logbook into database
    err := utils.SupabaseClient.DB.From("logbooks").Insert(input).Execute(nil)
    if err != nil {
        log.Printf("Error inserting logbook: %v", err) // Log the error for debugging
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add logbook: " + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Logbook added successfully"})
}

