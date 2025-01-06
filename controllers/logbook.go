package controllers

import (
	"auth-backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Logbook struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	StartDate string `json:"start_date" binding:"required"` // Start Date
	EndDate   string `json:"end_date" binding:"required"`   // End Date
	Activity  string `json:"activity" binding:"required"`   // Activity
	PIC       string `json:"pic" binding:"required"`        // Person In Charge
	Status    string `json:"status" `     // Status: approve, pending, reject
}

// GetLogbooks retrieves logbooks for the authenticated user
func GetLogbooks(c *gin.Context) {
	userID := c.GetString("user_id") // Retrieved from middleware
	var logbooks []Logbook

	err := utils.SupabaseClient.DB.From("logbooks").
		Select("id,user_id,start_date,end_date,activity,pic,status"). // Tambahkan `id` di sini
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

	input.UserID = userID

	log.Printf("Logbook Input: %+v", input)

	err := utils.SupabaseClient.DB.From("logbooks").Insert(input).Execute(nil)
	if err != nil {
		log.Printf("Error inserting logbook: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add logbook: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Logbook added successfully"})
}

func UpdateLogbook(c *gin.Context) {
    log.Println("UpdateLogbook: Request received")

    userID := c.GetString("user_id") 
    if userID == "" {
        log.Println("User ID missing or invalid")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing user_id"})
        return
    }

    log.Printf("User ID: %s", userID)

    logbookID := c.Param("id")
    if logbookID == "" {
        log.Println("Logbook ID missing")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing logbook ID"})
        return
    }

    log.Printf("Logbook ID: %s", logbookID)

    var input Logbook
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Printf("Error binding input JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    log.Printf("Input JSON: %+v", input)

    if input.StartDate == "" || input.EndDate == "" || input.Activity == "" || input.PIC == "" {
        log.Println("Validation failed for input fields")
        c.JSON(http.StatusBadRequest, gin.H{"error": "All fields except status are required"})
        return
    }

    updates := map[string]interface{}{
        "start_date": input.StartDate,
        "end_date":   input.EndDate,
        "activity":   input.Activity,
        "pic":        input.PIC,
    }

    log.Printf("Update data: %+v",updates)

    log.Println("Executing update query...")
    err := utils.SupabaseClient.DB.From("logbooks").
        Update(updates).
        Eq("id", logbookID).
        Execute(nil)

    if err != nil {
        log.Printf("Supabase Update Error: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update logbook"})
        return
    }

    log.Println("Update query executed successfully")
    c.JSON(http.StatusOK, gin.H{"message": "Logbook updated successfully"})
}

func DeleteLogbook(c *gin.Context) {
	userID := c.GetString("user_id") 

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing user_id"})
		return
	}

	logbookID := c.Param("id")
	if logbookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing logbook ID"})
		return
	}

	var logbooks []Logbook
	err := utils.SupabaseClient.DB.From("logbooks").
		Select("*").
		Eq("id", logbookID).
		Eq("user_id", userID).
		Execute(&logbooks)

	if err != nil || len(logbooks) == 0 {
		log.Printf("Error fetching logbook: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Logbook not found"})
		return
	}

	err = utils.SupabaseClient.DB.From("logbooks").
		Delete().
		Eq("id", logbookID).
		Execute(nil)

	if err != nil {
		log.Printf("Error deleting logbook: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete logbook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logbook deleted successfully"})
}
