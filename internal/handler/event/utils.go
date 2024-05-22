package event

import (
	"fmt"
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
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
