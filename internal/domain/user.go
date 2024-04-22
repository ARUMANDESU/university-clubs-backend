package domain

import (
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	AvatarURL string    `json:"avatar_url"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
	Barcode   string    `json:"barcode"`
	Major     string    `json:"major"`
	GroupName string    `json:"group_name"`
	Year      int       `json:"year"`
}

type Notification struct {
	UserID  int64  `json:"userID"`
	Message string `json:"message"`
}

func MapRoleStringToEnum(role string) userv1.Role {
	switch role {
	case "GUEST":
		return userv1.Role_GUEST
	case "USER":
		return userv1.Role_USER
	case "MODER":
		return userv1.Role_MODER
	case "ADMIN":
		return userv1.Role_ADMIN
	case "DSVR":
		return userv1.Role_DSVR
	default:
		return userv1.Role_GUEST // or any default value
	}
}

func UserObjectToDomain(user *userv1.UserObject) User {
	return User{
		ID:        user.GetUserId(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		AvatarURL: user.GetAvatarUrl(),
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
