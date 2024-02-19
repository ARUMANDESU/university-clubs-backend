package domain

import (
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	"time"
)

type Club struct {
	ID           int64
	Name         string
	OwnerID      *int64
	Description  string
	ClubType     string
	LogoURL      string
	BannerURL    string
	NumOFMembers *int64
	CreatedAt    time.Time
	Roles        []Role
}

type Role struct {
	ID          int
	Name        string
	Permissions []string
	Position    int32
	Color       int32
}

type Member struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Barcode   string `json:"barcode"`
	AvatarURL string `json:"avatar_url"`
	Roles     []Role
}

func ClubObjectToClub(clubObject *clubv1.ClubObject) *Club {
	roles := make([]Role, len(clubObject.GetRoles()))
	for i, role := range clubObject.GetRoles() {
		roles[i] = Role{
			Name:        role.GetName(),
			Permissions: role.GetPermissions(),
			Position:    role.GetPosition(),
			Color:       role.GetColor(),
		}
	}

	return &Club{
		ID:          clubObject.GetClubId(),
		Name:        clubObject.GetName(),
		Description: clubObject.GetDescription(),
		ClubType:    clubObject.GetClubType(),
		LogoURL:     clubObject.GetLogoUrl(),
		BannerURL:   clubObject.GetBannerUrl(),
		CreatedAt:   clubObject.GetCreatedAt().AsTime(),
		Roles:       roles,
	}
}

func UserObjectToMember(userObject *clubv1.UserObject) *Member {
	roles := make([]Role, len(userObject.GetRole()))
	for i, role := range userObject.GetRole() {
		roles[i] = Role{
			Name:        role.GetName(),
			Permissions: role.GetPermissions(),
			Position:    role.GetPosition(),
			Color:       role.GetColor(),
		}
	}

	return &Member{
		ID:        userObject.GetUserId(),
		Email:     userObject.GetEmail(),
		FirstName: userObject.GetFirstName(),
		LastName:  userObject.GetLastName(),
		Barcode:   userObject.GetBarcode(),
		AvatarURL: userObject.GetAvatarUrl(),
		Roles:     roles,
	}
}

func MapUserObjArrToMemberArr(ur []*clubv1.UserObject) []*Member {
	members := make([]*Member, len(ur))
	for i, u := range ur {
		members[i] = UserObjectToMember(u)
	}

	return members
}
