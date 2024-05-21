package event

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
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

	c.JSON(http.StatusOK, gin.H{"club": domain.EventProtoToDomain(res)})
}

func (h *Handler) ListEventsHandler(c *gin.Context) {
	const op = "EventHandler.ListEventsHandler"
	log := h.log.With(slog.String("op", op))

	query := c.Query("query")

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		log.Warn("failed to get page query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		log.Warn("failed to get page_size query parameter", logger.Err(err))
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

	Filter := eventv1.EventFilter{
		ClubId: int64(clubID),
		UserId: int64(userID),
		Status: c.Query("status"),
	}

	if len(tags) > 0 && tags[0] != "" {
		Filter.Tags = tags
	}

	res, err := h.eventClient.ListEvents(c, &eventv1.ListEventsRequest{
		Query:      query,
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

func (h *Handler) ListPublishedEventsHandler(c *gin.Context) {
	const op = "EventHandler.ListPublishedEventsHandler"
	log := h.log.With(slog.String("op", op))

	query := c.Query("query")

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		log.Warn("failed to get page query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		log.Warn("failed to get page_size query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.eventClient.ListEvents(c, &eventv1.ListEventsRequest{
		Query:      query,
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
		Filter:     &eventv1.EventFilter{Status: "IN_PROGRESS"},
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
