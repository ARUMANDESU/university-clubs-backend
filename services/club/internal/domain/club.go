package domain

type Club struct {
	ID             int64
	Name           string
	ClubOwner      Student
	ClubModerators []Student
	Members        []Student
}
