package domain

type Student struct {
	User
	Barcode     string `json:"barcode"`
	PhoneNumber string `json:"phoneNumber"`
	Major       string `json:"major"`
	Group       string `json:"group"`
	Year        int    `json:"year"`
}
