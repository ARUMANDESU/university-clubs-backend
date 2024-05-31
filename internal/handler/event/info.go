package event

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handler) GetEventHandler(c *gin.Context) {
	const op = "EventHandler.GetEventHandler"
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

	log.Debug("getting event", slog.String("event_id", eventID), slog.Int64("user_id", userID))

	res, err := h.eventClient.GetEvent(c, &eventv1.GetEventRequest{EventId: eventID, UserId: userID})
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

	c.JSON(http.StatusOK,
		gin.H{
			"event":                domain.ProtoToEvent(res.Event),
			"user_status":          domain.UserEventStatusFromProto(res.UserStatus),
			"participation_status": domain.ParticipationStatusFromProto(res.ParticipantStatus),
		},
	)
}

func (h *Handler) ListEventsHandler(c *gin.Context) {
	const op = "EventHandler.ListEventsHandler"
	log := h.log.With(slog.String("op", op))

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clubID, err := utils.GetIntFromQuery(c, "club_id")
	if err != nil && !strings.Contains(err.Error(), "query parameter must be provided") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := utils.GetIntFromQuery(c, "user_id")
	if err != nil && !strings.Contains(err.Error(), "query parameter must be provided") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tags := strings.Split(c.Query("tags"), ",")
	statuses := strings.Split(c.Query("status"), ",")

	var paths []string
	isHiddenForNonMembers, err := utils.GetBoolFromQuery(c, "is_hidden_for_non_members")
	if err == nil {
		paths = append(paths, "is_hidden_for_non_members")
	} else if !strings.Contains(err.Error(), "query parameter must be provided") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Filter := eventv1.EventFilter{
		ClubId:                int64(clubID),
		UserId:                int64(userID),
		FromDate:              c.Query("from_date"),
		TillDate:              c.Query("till_date"),
		IsHiddenForNonMembers: isHiddenForNonMembers,
	}

	if len(tags) > 0 && tags[0] != "" {
		Filter.Tags = tags
	}
	if len(statuses) > 0 && statuses[0] != "" {
		Filter.Status = statuses
	}

	res, err := h.eventClient.ListEvents(c, &eventv1.ListEventsRequest{
		Query:      c.Query("query"),
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.Query("sort_order"),
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
		Filter:     &Filter,
		FilterMask: &field_mask.FieldMask{Paths: paths},
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": domain.ProtoToEventArr(res.Events), "metadata": res.Metadata})
}

func (h *Handler) ListPublishedEventsHandler(c *gin.Context) {
	const op = "EventHandler.ListPublishedEventsHandler"
	log := h.log.With(slog.String("op", op))

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.eventClient.ListEvents(c, &eventv1.ListEventsRequest{
		Query:      c.Query("query"),
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.Query("sort_order"),
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
		Filter:     &eventv1.EventFilter{Status: []string{"IN_PROGRESS"}},
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": domain.ProtoToEventArr(res.Events), "metadata": res.Metadata})
}

func (h *Handler) ListClubEventsHandler(c *gin.Context) {
	const op = "EventHandler.ListClubEventsHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tags := strings.Split(c.Query("tags"), ",")

	Filter := eventv1.EventFilter{
		ClubId:   clubID,
		Status:   []string{"IN_PROGRESS", "FINISHED", "CANCELED", "ARCHIVED"},
		FromDate: c.Query("from_date"),
		TillDate: c.Query("till_date"),
	}

	if len(tags) > 0 && tags[0] != "" {
		Filter.Tags = tags
	}

	res, err := h.eventClient.ListEvents(c, &eventv1.ListEventsRequest{
		Query:      c.Query("query"),
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.Query("sort_order"),
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
		Filter:     &Filter,
	})

	if err != nil {
		switch {
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": domain.ProtoToEventArr(res.Events), "metadata": res.Metadata})
}

func (h *Handler) ListRestrictedClubEvents(c *gin.Context) {
	const op = "EventHandler.ListRestrictedClubEvents"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	err = h.canHandleEventManagement(c, clubID, userID)
	if err != nil {
		return
	}

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDFilter, err := utils.GetIntFromQuery(c, "user_id")
	if err != nil && !strings.Contains(err.Error(), "query parameter must be provided") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tags := strings.Split(c.Query("tags"), ",")
	statuses := strings.Split(c.Query("status"), ",")

	Filter := eventv1.EventFilter{
		ClubId:   clubID,
		UserId:   int64(userIDFilter),
		FromDate: c.Query("from_date"),
		TillDate: c.Query("till_date"),
	}

	if len(tags) > 0 && tags[0] != "" {
		Filter.Tags = tags
	}
	if len(statuses) > 0 && statuses[0] != "" {
		Filter.Status = statuses
	}

	res, err := h.eventClient.ListEvents(c, &eventv1.ListEventsRequest{
		Query:      c.Query("query"),
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.Query("sort_order"),
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
		Filter:     &Filter,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": domain.ProtoToEventArr(res.Events), "metadata": res.Metadata})
}
func (h *Handler) ListCollaboratorInvitesHandler(c *gin.Context) {
	const op = "EventHandler.ListCollaboratorInvitesHandler"
	log := h.log.With(slog.String("op", op))

	eventID := c.Params.ByName("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "event_id parameter must be provided"})
		return
	}

	/*userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)*/

	res, err := h.eventClient.GetClubInvites(c, &eventv1.GetInvitesRequest{
		EventId: eventID,
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

	c.JSON(http.StatusOK, gin.H{"invites": res.Invites})
}

func (h *Handler) ListCollaboratorRequestsHandler(c *gin.Context) {
	const op = "EventHandler.ListCollaboratorInvitesHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
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

	err = h.canHandleEventManagement(c, clubID, userID)
	if err != nil {
		return
	}

	res, err := h.eventClient.GetClubInvites(c, &eventv1.GetInvitesRequest{
		ClubId: clubID,
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

	c.JSON(http.StatusOK, gin.H{"invites": res.Invites})
}

func (h *Handler) ListOrganizerInvitesHandler(c *gin.Context) {
	const op = "EventHandler.ListOrganizerInvitesHandler"
	log := h.log.With(slog.String("op", op))

	eventID := c.Params.ByName("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "event_id parameter must be provided"})
		return
	}

	/*userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)*/

	res, err := h.eventClient.GetOrganizerInvites(c, &eventv1.GetInvitesRequest{
		EventId: eventID,
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

	c.JSON(http.StatusOK, gin.H{"invites": res.Invites})
}

func (h *Handler) GetUserInvites(c *gin.Context) {
	const op = "EventHandler.GetUserInvites"
	log := h.log.With(slog.String("op", op))

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	res, err := h.eventClient.GetOrganizerInvites(c, &eventv1.GetInvitesRequest{
		UserId: userID,
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

	c.JSON(http.StatusOK, gin.H{"invites": res.GetInvites()})
}

func (h *Handler) ListParticipantsHandler(c *gin.Context) {
	const op = "EventHandler.ListParticipantsHandler"
	log := h.log.With(slog.String("op", op))

	eventID := c.Params.ByName("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "event_id parameter must be provided"})
		return
	}

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil && !strings.Contains(err.Error(), "query parameter must be provided") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil && !strings.Contains(err.Error(), "query parameter must be provided") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.eventClient.ListParticipants(c, &eventv1.ListParticipantsRequest{
		EventId:    eventID,
		Query:      c.Query("query"),
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		h.handleEventError(c, log, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"participants": domain.ProtoToEventUserArr(res.Participants), "metadata": res.Metadata})
}
