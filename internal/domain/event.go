package domain

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
)

type Event struct {
	ID                 string       `json:"id"`
	ClubId             int64        `json:"club_id"`
	OwnerId            int64        `json:"owner_id"`
	CollaboratorClubs  []EventClub  `json:"collaborator_clubs"`
	Organizers         []Organizer  `json:"organizers"`
	Title              string       `json:"title,omitempty"`
	Description        string       `json:"description,omitempty"`
	Type               string       `json:"type,omitempty"`
	Status             string       `json:"status,omitempty"`
	Tags               []string     `json:"tags,omitempty"`
	MaxParticipants    uint32       `json:"max_participants,omitempty"`
	ParticipantsCount  uint32       `json:"participants_count,omitempty"`
	LocationLink       string       `json:"location_link,omitempty"`
	LocationUniversity string       `json:"location_university,omitempty"`
	StartDate          string       `json:"start_date"`
	EndDate            string       `json:"end_date"`
	CoverImages        []CoverImage `json:"cover_images,omitempty"`
	AttachedImages     []EventFile  `json:"attached_images,omitempty"`
	AttachedFiles      []EventFile  `json:"attached_files,omitempty"`
	CreatedAt          string       `json:"created_at"`
	UpdatedAt          string       `json:"updated_at"`
	DeletedAt          string       `json:"deleted_at"`
}

type EventFile struct {
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
	Type string `json:"type" bson:"type"`
}

type CoverImage struct {
	EventFile
	Position uint32 `json:"position" bson:"position"`
}
type EventUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Barcode   string `json:"barcode"`
	AvatarURL string `json:"avatar_url"`
}

type Organizer struct {
	EventUser `json:",inline"`
	ClubId    int64 `json:"club_id"`
	ByWhoId   int64 `json:"by_who_id"`
}
type EventClub struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	LogoURL string `json:"logo_url"`
}

func EventProtoToDomain(event *eventv1.EventObject) *Event {
	return &Event{
		ID:                 event.GetId(),
		ClubId:             event.GetClubId(),
		OwnerId:            event.GetOwnerId(),
		CollaboratorClubs:  ProtoToCollaboratorClubArr(event.GetCollaboratorClubs()),
		Organizers:         ProtoToOrganizerArr(event.GetOrganizers()),
		Title:              event.GetTitle(),
		Description:        event.GetDescription(),
		Type:               event.GetType(),
		Status:             event.GetStatus(),
		Tags:               event.GetTags(),
		MaxParticipants:    event.GetMaxParticipants(),
		ParticipantsCount:  event.GetParticipantsCount(),
		LocationLink:       event.GetLocationLink(),
		LocationUniversity: event.GetLocationUniversity(),
		StartDate:          event.GetStartDate(),
		EndDate:            event.GetEndDate(),
		CoverImages:        ProtoToCoverImageArr(event.GetCoverImages()),
		AttachedImages:     ProtoToEventFileArr(event.GetAttachedImages()),
		AttachedFiles:      ProtoToEventFileArr(event.GetAttachedFiles()),
		CreatedAt:          event.GetCreatedAt(),
		UpdatedAt:          event.GetUpdatedAt(),
		DeletedAt:          event.GetDeletedAt(),
	}
}

func ProtoToCoverImage(coverImage *eventv1.CoverImage) *CoverImage {
	return &CoverImage{
		EventFile: EventFile{
			Name: coverImage.GetName(),
			Url:  coverImage.GetUrl(),
			Type: coverImage.GetType(),
		},
		Position: uint32(coverImage.GetPosition()),
	}
}

func ProtoToEventFile(eventFile *eventv1.FileObject) *EventFile {
	return &EventFile{
		Name: eventFile.GetName(),
		Url:  eventFile.GetUrl(),
		Type: eventFile.GetType(),
	}
}

func ProtoToOrganizer(organizer *eventv1.OrganizerObject) *Organizer {
	return &Organizer{
		EventUser: EventUser{
			ID:        organizer.GetId(),
			FirstName: organizer.GetFirstName(),
			LastName:  organizer.GetLastName(),
			Barcode:   organizer.GetBarcode(),
			AvatarURL: organizer.GetAvatarUrl(),
		},
		ClubId: organizer.GetClubId(),
	}
}

func ProtoToEventClub(eventClub *eventv1.ClubObject) *EventClub {
	return &EventClub{
		ID:      eventClub.GetId(),
		Name:    eventClub.GetName(),
		LogoURL: eventClub.GetLogoUrl(),
	}
}

func ProtoToEventUser(eventUser *eventv1.UserObject) *EventUser {
	return &EventUser{
		ID:        eventUser.GetId(),
		FirstName: eventUser.GetFirstName(),
		LastName:  eventUser.GetLastName(),
		Barcode:   eventUser.GetBarcode(),
		AvatarURL: eventUser.GetAvatarUrl(),
	}
}

func ProtoToEventFileArr(eventFiles []*eventv1.FileObject) []EventFile {
	files := make([]EventFile, len(eventFiles))
	for i, file := range eventFiles {
		files[i] = *ProtoToEventFile(file)
	}
	return files
}

func ProtoToCoverImageArr(coverImages []*eventv1.CoverImage) []CoverImage {
	images := make([]CoverImage, len(coverImages))
	for i, image := range coverImages {
		images[i] = *ProtoToCoverImage(image)
	}
	return images
}

func ProtoToOrganizerArr(organizers []*eventv1.OrganizerObject) []Organizer {
	orgs := make([]Organizer, len(organizers))
	for i, organizer := range organizers {
		orgs[i] = *ProtoToOrganizer(organizer)
	}
	return orgs
}

func ProtoToEventClubArr(eventClubs []*eventv1.ClubObject) []EventClub {
	clubs := make([]EventClub, len(eventClubs))
	for i, club := range eventClubs {
		clubs[i] = *ProtoToEventClub(club)
	}
	return clubs

}

func ProtoToEventUserArr(eventUsers []*eventv1.UserObject) []EventUser {
	users := make([]EventUser, len(eventUsers))
	for i, user := range eventUsers {
		users[i] = *ProtoToEventUser(user)
	}
	return users
}

func ProtoToEventArr(events []*eventv1.EventObject) []Event {
	e := make([]Event, len(events))
	for i, event := range events {
		e[i] = *EventProtoToDomain(event)
	}
	return e
}

func EventToProto(event *Event) *eventv1.EventObject {
	return &eventv1.EventObject{
		Id:                 event.ID,
		ClubId:             event.ClubId,
		OwnerId:            event.OwnerId,
		Title:              event.Title,
		Description:        event.Description,
		Type:               event.Type,
		Status:             event.Status,
		Tags:               event.Tags,
		MaxParticipants:    event.MaxParticipants,
		ParticipantsCount:  event.ParticipantsCount,
		LocationLink:       event.LocationLink,
		LocationUniversity: event.LocationUniversity,
		StartDate:          event.StartDate,
		EndDate:            event.EndDate,
		CreatedAt:          event.CreatedAt,
		UpdatedAt:          event.UpdatedAt,
		DeletedAt:          event.DeletedAt,
	}
}

func CoverImageToProto(coverImage *CoverImage) *eventv1.CoverImage {
	return &eventv1.CoverImage{
		Name:     coverImage.Name,
		Url:      coverImage.Url,
		Type:     coverImage.Type,
		Position: int32(coverImage.Position),
	}
}

func EventFileToProto(eventFile *EventFile) *eventv1.FileObject {
	return &eventv1.FileObject{
		Name: eventFile.Name,
		Url:  eventFile.Url,
		Type: eventFile.Type,
	}
}

func ProtoToCollaboratorClub(collaboratorClub *eventv1.ClubObject) *EventClub {
	return &EventClub{
		ID:      collaboratorClub.GetId(),
		Name:    collaboratorClub.GetName(),
		LogoURL: collaboratorClub.GetLogoUrl(),
	}
}

func ProtoToCollaboratorClubArr(collaboratorClubs []*eventv1.ClubObject) []EventClub {
	clubs := make([]EventClub, len(collaboratorClubs))
	for i, club := range collaboratorClubs {
		clubs[i] = *ProtoToCollaboratorClub(club)
	}
	return clubs
}
