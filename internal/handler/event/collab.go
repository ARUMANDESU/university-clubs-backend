package event

import (
	"fmt"
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	"github.com/ARUMANDESU/uniclubs-protos/gen/go/posts"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

func (h *Handler) AddCollaboratorHandler(c *gin.Context) {
	const op = "EventHandler.AddCollaboratorHandler"
	log := h.log.With(slog.String("op", op))

	eventID := c.Params.ByName("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "event_id parameter must be provided"})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	var input struct {
		ClubId int64 `json:"club_id"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	club, err := h.clubClient.GetClub(c, &clubv1.GetClubRequest{
		ClubId: input.ClubId,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	_, err = h.eventClient.AddCollaborator(c, &eventv1.AddCollaboratorRequest{
		EventId: eventID,
		UserId:  userID,
		Club: &posts.ClubObject{
			Id:      club.ClubId,
			Name:    club.Name,
			LogoUrl: club.LogoUrl,
		},
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.AlreadyExists:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) HandleCollaboratorRequestHandler(c *gin.Context) {
	const op = "EventHandler.HandleCollaboratorRequestHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inviteID := c.Params.ByName("invite_id")
	if inviteID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invite_id parameter must be provided"})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	var input struct {
		Action string `json:"action"` // "accept" or "reject"
	}

	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var action eventv1.HandleInvite_Action
	switch input.Action {
	case "accept":
		action = eventv1.HandleInvite_Action_ACCEPT
	case "reject":
		action = eventv1.HandleInvite_Action_REJECT
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "action must be either 'accept' or 'reject'"})
	}

	user, err := h.userClient.GetUser(c, &userv1.GetUserRequest{UserId: userID})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			log.Warn("user not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	res, err := h.eventClient.HandleInviteClub(c, &eventv1.HandleInviteClubRequest{
		InviteId: inviteID,
		User: &eventv1.UserObject{
			Id:        user.UserId,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Barcode:   user.Barcode,
			AvatarUrl: user.AvatarUrl,
		},
		ClubId: clubID,
		Action: action,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.AlreadyExists, status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	if action == eventv1.HandleInvite_Action_REJECT {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(res)})
}

func (h *Handler) RemoveCollaboratorHandler(c *gin.Context) {
	const op = "EventHandler.RemoveCollaboratorHandler"
	log := h.log.With(slog.String("op", op))

	eventID := c.Params.ByName("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "event_id parameter must be provided"})
		return
	}

	collaboratorID, err := utils.GetIntFromParams(c.Params, "collaborator_id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	res, err := h.eventClient.RemoveCollaborator(c, &eventv1.RemoveCollaboratorRequest{
		EventId: eventID,
		UserId:  userID,
		ClubId:  collaboratorID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.AlreadyExists, status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(res)})
}
func (h *Handler) CancelCollaboratorRequestHandler(c *gin.Context) {
	const op = "EventHandler.CancelCollaboratorRequestHandler"
	log := h.log.With(slog.String("op", op))

	inviteID := c.Params.ByName("invite_id")
	if inviteID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invite_id parameter must be provided"})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	_, err := h.eventClient.RevokeInviteClub(c, &eventv1.RevokeInviteRequest{
		InviteId: inviteID,
		UserId:   userID,
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
		return

	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) AddOrganizerHandler(c *gin.Context) {
	const op = "EventHandler.AddOrganizerHandler"
	log := h.log.With(slog.String("op", op))

	eventID := c.Params.ByName("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "event_id parameter must be provided"})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	var input struct {
		UserId int64 `json:"user_id"`
		ClubId int64 `json:"club_id"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target, err := h.clubClient.GetClubMember(c, &clubv1.GetClubMemberRequest{
		ClubId: input.ClubId,
		UserId: input.UserId,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	_, err = h.eventClient.AddOrganizer(c, &eventv1.AddOrganizerRequest{
		EventId: eventID,
		UserId:  userID,
		Target: &eventv1.UserObject{
			Id:        target.GetUserId(),
			FirstName: target.GetFirstName(),
			LastName:  target.GetLastName(),
			Barcode:   target.GetBarcode(),
			AvatarUrl: target.GetAvatarUrl(),
		},
		TargetClubId: input.ClubId,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.AlreadyExists:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusNoContent)
}
func (h *Handler) RemoveOrganizerHandler(c *gin.Context) {
	const op = "EventHandler.RemoveOrganizerHandler"
	log := h.log.With(slog.String("op", op))

	eventID := c.Params.ByName("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "event_id parameter must be provided"})
		return
	}

	organizerID, err := utils.GetIntFromParams(c.Params, "organizer_id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	res, err := h.eventClient.RemoveOrganizer(c, &eventv1.RemoveOrganizerRequest{
		EventId:     eventID,
		UserId:      userID,
		OrganizerId: organizerID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.AlreadyExists, status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(res)})
}

func (h *Handler) CancelOrganizerRequestHandler(c *gin.Context) {
	const op = "EventHandler.CancelOrganizerRequestHandler"
	log := h.log.With(slog.String("op", op))

	inviteID := c.Params.ByName("invite_id")
	if inviteID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invite_id parameter must be provided"})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	_, err := h.eventClient.RevokeInviteUser(c, &eventv1.RevokeInviteRequest{
		InviteId: inviteID,
		UserId:   userID,
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
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) HandleUserInvite(c *gin.Context) {
	const op = "EventHandler.HandleUserInvite"
	log := h.log.With(slog.String("op", op))

	inviteId := c.Params.ByName("invite_id")
	if inviteId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invite_id parameter must be provided"})
		return
	}

	userIdFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userId := userIdFromCtx.(int64)

	var input struct {
		Action string `json:"action"` // accept,  reject
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to bind json: %w", err)})
		return
	}

	var action eventv1.HandleInvite_Action
	switch input.Action {
	case "accept":
		action = eventv1.HandleInvite_Action_ACCEPT
	case "reject":
		action = eventv1.HandleInvite_Action_REJECT
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "action must be either 'accept' or 'reject'"})
		return
	}

	res, err := h.eventClient.HandleInviteUser(c, &eventv1.HandleInviteUserRequest{
		InviteId: inviteId,
		UserId:   userId,
		Action:   action,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.AlreadyExists, status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	if action == eventv1.HandleInvite_Action_REJECT {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(res)})
}
