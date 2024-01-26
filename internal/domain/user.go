package domain

import "time"

type User struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	Role        string    `json:"role"`
	Barcode     string    `json:"barcode"`
	PhoneNumber string    `json:"phone_number"`
	Major       string    `json:"major"`
	GroupName   string    `json:"group_name"`
	Year        int       `json:"year"`
}
