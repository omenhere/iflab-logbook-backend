package models

type User struct {
	ID           string `json:"id"`
	Nim        string `json:"nim"`
	PasswordHash string `json:"password_hash"`
}
