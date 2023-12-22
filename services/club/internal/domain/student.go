package domain

import "time"

type Student struct {
	ID         int64     `json:"ID"`
	FirstName  string    `json:"firstName"`
	SecondName string    `json:"secondName"`
	Email      string    `json:"email"`
	Activated  bool      `json:"activated"`
	CreatedAt  time.Time `json:"created_at"`
	Role       string    `json:"role"`
	Major      string    `json:"major"`
	Group      string    `json:"group"`
	Year       int       `json:"year"`
}
