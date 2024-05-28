package event

import (
	"fmt"
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

func (h *Handler) canHandleEventManagement(c *gin.Context, clubID, userID int64) error {
	const op = "EventHandler.canHandleEventManagement"
	log := h.log.With(slog.String("op", op))

	havePermResp, err := h.clubClient.HavePermissionTo(c, &clubv1.HavePermissionToRequest{
		ClubId:     clubID,
		UserId:     userID,
		Permission: domain.ManageEvents,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return err
	}

	if !havePermResp.GetHasPermission() {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you don't have permission to create event"})
		return fmt.Errorf("user %d doesn't have permission to create event", userID)
	}

	return nil
}

func (h *Handler) getUser(c *gin.Context, userID int64) (*eventv1.UserObject, error) {
	const op = "EventHandler.getUser"
	log := h.log.With(slog.String("op", op))

	userResp, err := h.userClient.GetUser(c, &userv1.GetUserRequest{UserId: userID})
	if err != nil {
		switch {
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return nil, err
	}

	user := eventv1.UserObject{
		Id:        userResp.GetUserId(),
		FirstName: userResp.GetFirstName(),
		LastName:  userResp.GetLastName(),
		Barcode:   userResp.GetBarcode(),
		AvatarUrl: userResp.GetAvatarUrl(),
	}

	return &user, nil
}

func (h *Handler) buildUpdatePaths(c *gin.Context, updateRequest *eventv1.UpdateEventRequest) ([]string, error) {
	var input struct {
		Title                 *string             `json:"title,omitempty"`
		Description           *string             `json:"description,omitempty"`
		StartDate             *string             `json:"start_date,omitempty"`
		EndDate               *string             `json:"end_date,omitempty"`
		Tags                  []string            `json:"tags,omitempty"`
		Type                  *string             `json:"type,omitempty"`
		MaxParticipants       *int32              `json:"max_participants,omitempty"`
		LocationUni           *string             `json:"location_uni,omitempty"`
		LocationLink          *string             `json:"location_link,omitempty"`
		CoverImage            []domain.CoverImage `json:"cover_images,omitempty"`
		AttachedFiles         []domain.EventFile  `json:"attached_files,omitempty"`
		AttachedImages        []domain.EventFile  `json:"attached_images,omitempty"`
		IsHiddenForNonMembers *bool               `json:"is_hidden_for_non_members"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	var paths []string

	if input.Title != nil {
		paths = append(paths, "title")
		updateRequest.Title = *input.Title
	}
	if input.Description != nil {
		paths = append(paths, "description")
		updateRequest.Description = *input.Description
	}
	if input.StartDate != nil {
		paths = append(paths, "start_date")
		updateRequest.StartDate = *input.StartDate
	}
	if input.EndDate != nil {
		paths = append(paths, "end_date")
		updateRequest.EndDate = *input.EndDate
	}
	if input.Type != nil {
		paths = append(paths, "type")
		updateRequest.Type = *input.Type
	}
	if input.MaxParticipants != nil {
		paths = append(paths, "max_participants")
		updateRequest.MaxParticipants = *input.MaxParticipants
	}
	if input.LocationUni != nil {
		paths = append(paths, "location_university")
		updateRequest.LocationUniversity = *input.LocationUni
	}
	if input.LocationLink != nil {
		paths = append(paths, "location_link")
		updateRequest.LocationLink = *input.LocationLink
	}
	if input.IsHiddenForNonMembers != nil {
		paths = append(paths, "is_hidden_for_non_members")
		updateRequest.IsHiddenForNonMembers = *input.IsHiddenForNonMembers
	}
	if input.Tags != nil && len(input.Tags) > 0 {
		paths = append(paths, "tags")
		updateRequest.Tags = input.Tags
	}
	if input.CoverImage != nil && len(input.CoverImage) > 0 {
		paths = append(paths, "cover_images")
		updateRequest.CoverImages = domain.CoverImageToProtoArr(input.CoverImage)
	}
	if input.AttachedFiles != nil && len(input.AttachedFiles) > 0 {
		paths = append(paths, "attached_files")
		updateRequest.AttachedFiles = domain.EventFileToProtoArr(input.AttachedFiles)
	}
	if input.AttachedImages != nil && len(input.AttachedImages) > 0 {
		paths = append(paths, "attached_images")
		updateRequest.AttachedImages = domain.EventFileToProtoArr(input.AttachedImages)
	}

	return paths, nil
}

// handleEventError handles errors from the UpdateEvent call
func (h *Handler) handleEventError(c *gin.Context, log *slog.Logger, err error) {
	switch status.Code(err) {
	case codes.InvalidArgument:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
	case codes.NotFound:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
	case codes.PermissionDenied:
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
	case codes.FailedPrecondition:
		c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": status.Convert(err).Message()})
	case codes.AlreadyExists, codes.Aborted:
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
	default:
		log.Error("internal error", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
