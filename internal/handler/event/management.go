package event

import (
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

func (h *Handler) CreateEventHandler(c *gin.Context) {
	const op = "ClubHandler.CreateEventHandler"

	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	var getClubResponse *clubv1.ClubObject
	var getUserResponse *userv1.UserObject

	g, ctx := errgroup.WithContext(c)

	g.Go(func() error {
		var err error
		getClubResponse, err = h.clubClient.GetClub(ctx, &clubv1.GetClubRequest{
			ClubId: clubID,
		})
		return err
	})

	g.Go(func() error {
		var err error
		getUserResponse, err = h.userClient.GetUser(ctx, &userv1.GetUserRequest{
			UserId: userID,
		})
		return err
	})

	if err := g.Wait(); err != nil {
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

	eventResponse, err := h.eventClient.CreateEvent(c, &eventv1.CreateEventRequest{
		Club: &eventv1.ClubObject{
			Id:      getClubResponse.GetClubId(),
			Name:    getClubResponse.GetName(),
			LogoUrl: getClubResponse.GetLogoUrl(),
		},
		User: &eventv1.UserObject{
			Id:        getUserResponse.GetUserId(),
			FirstName: getUserResponse.GetFirstName(),
			LastName:  getUserResponse.GetLastName(),
			Barcode:   getUserResponse.GetBarcode(),
			AvatarUrl: getUserResponse.GetAvatarUrl(),
		},
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		}
	}
	c.JSON(http.StatusCreated, gin.H{"event": eventResponse})
}
