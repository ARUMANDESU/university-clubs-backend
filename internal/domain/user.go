package domain

import (
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
	Barcode   string    `json:"barcode"`
	Major     string    `json:"major"`
	GroupName string    `json:"group_name"`
	Year      int       `json:"year"`
}

func UserObjectToDomain(user *userv1.UserObject) User {
	return User{
		ID:        user.GetUserId(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Email:     user.GetEmail(),
		CreatedAt: user.GetCreatedAt().AsTime(),
		Role:      user.GetRole().String(),
		Barcode:   user.GetBarcode(),
		Major:     user.GetMajor(),
		GroupName: user.GetGroupName(),
		Year:      int(user.GetYear()),
	}
}

func MapUserObjectArrToDomain(usersObject []*userv1.UserObject) []User {
	users := make([]User, len(usersObject))
	for i, user := range usersObject {
		users[i] = UserObjectToDomain(user)
	}
	return users
}
