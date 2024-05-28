package event

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (h *Handler) AddParticipantHandler(c *gin.Context) {
	const op = "event.Handler.AddParticipantHandler"
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

	_, err := h.eventClient.ParticipateEvent(c, &eventv1.EventActionRequest{
		EventId: eventID,
		UserId:  userID,
	})
	if err != nil {
		h.handleEventError(c, log, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) LeaveEventHandler(c *gin.Context) {
	const op = "event.Handler.RemoveParticipantHandler"
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

	_, err := h.eventClient.CancelParticipation(c, &eventv1.EventActionRequest{
		EventId: eventID,
		UserId:  userID,
	})
	if err != nil {
		h.handleEventError(c, log, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) RemoveParticipantHandler(c *gin.Context) {
	const op = "event.Handler.RemoveParticipantHandler"
	log := h.log.With(slog.String("op", op))

	eventID := c.Params.ByName("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "event_id parameter must be provided"})
		return
	}

	participantId, err := utils.GetIntFromParams(c.Params, "participant_id")
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

	_, err = h.eventClient.KickParticipant(c, &eventv1.KickParticipantRequest{
		EventId:       eventID,
		UserId:        userID,
		ParticipantId: participantId,
	})
	if err != nil {
		h.handleEventError(c, log, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) BanParticipantHandler(c *gin.Context) {
	const op = "event.Handler.BanParticipantHandler"
	log := h.log.With(slog.String("op", op))

	eventID := c.Params.ByName("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "event_id parameter must be provided"})
		return
	}

	participantId, err := utils.GetIntFromParams(c.Params, "participant_id")
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

	var input struct {
		Reason string `json:"reason,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.eventClient.BanParticipant(c, &eventv1.BanParticipantRequest{
		EventId:       eventID,
		UserId:        userID,
		ParticipantId: participantId,
		Reason:        input.Reason,
	})
	if err != nil {
		h.handleEventError(c, log, err)
		return
	}

	c.Status(http.StatusNoContent)
}
