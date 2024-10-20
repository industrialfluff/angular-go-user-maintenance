package models

type User struct {
	User_id     int    `json:"user_id"`
	User_name   string `json:"user_name"`
	First_name  string `json:"first_name"`
	Last_name   string `json:"last_name"`
	Email       string `json:"email"`
	User_status string `json:"user_status"`
	Department  string `json:"department"`
}
