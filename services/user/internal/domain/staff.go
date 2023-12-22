package domain

type Staff struct {
	User
	Position string `json:"position"`
}
