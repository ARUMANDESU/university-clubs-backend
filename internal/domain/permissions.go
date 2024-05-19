package domain

const (
	Administrator uint64 = 1 << iota
	ManageClub
	ManageMembership
	KickMember
	BanMember
	ManageRoles
	ManageEvents
	ALL = Administrator | ManageRoles | ManageMembership | KickMember | BanMember | ManageClub | ManageEvents
)

var Names = map[uint64]string{
	Administrator:    "Administrator",
	ManageClub:       "ManageClub",
	ManageMembership: "ManageMembership",
	KickMember:       "KickMember",
	BanMember:        "BanMember",
	ManageRoles:      "ManageRoles",
	ManageEvents:     "ManageEvents",
}

var Values = map[string]uint64{
	"Administrator":    Administrator,
	"ManageClub":       ManageClub,
	"ManageMembership": ManageMembership,
	"KickMember":       KickMember,
	"BanMember":        BanMember,
	"ManageRoles":      ManageRoles,
	"ManageEvents":     ManageEvents,
}

var PermissionList = []string{
	"Administrator",
	"ManageClub",
	"ManageMembership",
	"KickMember",
	"BanMember",
	"ManageRoles",
	"ManageEvents",
}
