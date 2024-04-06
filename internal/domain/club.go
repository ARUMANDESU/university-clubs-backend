package domain

import (
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	"time"
)

type Club struct {
	ID           int64
	Name         string
	OwnerID      int64
	Description  string
	ClubType     string
	LogoURL      string
	BannerURL    string
	NumOFMembers int64
	CreatedAt    time.Time
	Roles        []*Role
}

type Role struct {
	ID          int64
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
	Roles     []int64
}

func ProtoToRole(role *clubv1.Role) *Role {
	return &Role{
		ID:          role.GetId(),
		Name:        role.GetName(),
		Permissions: role.GetPermissions(),
		Position:    role.GetPosition(),
		Color:       role.GetColor(),
	}
}

func MapProtoToRoleArr(r []*clubv1.Role) []*Role {
	roles := make([]*Role, len(r))
	for i, role := range r {
		roles[i] = ProtoToRole(role)
	}

	return roles
}

func ClubObjectToClub(clubObject *clubv1.ClubObject) *Club {
	roles := make([]*Role, len(clubObject.GetRoles()))
	if clubObject.GetRoles() != nil {
		for i, role := range clubObject.GetRoles() {
			roles[i] = ProtoToRole(role)
		}
	}

	return &Club{
		ID:           clubObject.GetClubId(),
		OwnerID:      clubObject.GetOwnerId(),
		Name:         clubObject.GetName(),
		Description:  clubObject.GetDescription(),
		ClubType:     clubObject.GetClubType(),
		LogoURL:      clubObject.GetLogoUrl(),
		BannerURL:    clubObject.GetBannerUrl(),
		CreatedAt:    clubObject.GetCreatedAt().AsTime(),
		NumOFMembers: clubObject.GetNumberOfMembers(),
		Roles:        roles,
	}
}

func MapClubObjArrToClubArr(clubObjects []*clubv1.ClubObject) []*Club {
	clubs := make([]*Club, len(clubObjects))
	for i, clubObject := range clubObjects {
		clubs[i] = ClubObjectToClub(clubObject)
	}

	return clubs
}

func UserObjectToMember(userObject *clubv1.UserObject) *Member {
	return &Member{
		ID:        userObject.GetUserId(),
		Email:     userObject.GetEmail(),
		FirstName: userObject.GetFirstName(),
		LastName:  userObject.GetLastName(),
		Barcode:   userObject.GetBarcode(),
		AvatarURL: userObject.GetAvatarUrl(),
		Roles:     userObject.GetRoles(),
	}
}

func MapUserObjArrToMemberArr(ur []*clubv1.UserObject) []*Member {
	members := make([]*Member, len(ur))
	for i, u := range ur {
		members[i] = UserObjectToMember(u)
	}

	return members
}
