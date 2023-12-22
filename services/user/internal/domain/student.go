package domain

type Student struct {
	User
	Major string `json:"major"`
	Group string `json:"group"`
	Year  int    `json:"year"`
}
