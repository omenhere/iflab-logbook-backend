package models

type Logbook struct {
    ID        string `json:"id,omitempty"` // UUID
    UserID    string `json:"user_id"`      // UUID
    StartDate string `json:"start_date" binding:"required"`
    EndDate   string `json:"end_date" binding:"required"`
    Activity  string `json:"activity" binding:"required"`
    PIC       string `json:"pic" binding:"required"`
    Status    string `json:"status" binding:"required"`
}
