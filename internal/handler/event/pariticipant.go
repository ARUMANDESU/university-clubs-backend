package event

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
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
