package event

import (
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

func (h *Handler) CreateEventHandler(c *gin.Context) {
	const op = "EventHandler.CreateEventHandler"
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
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"event": eventResponse})
}

func (h *Handler) UpdateEventHandler(c *gin.Context) {
	const op = "EventHandler.UpdateEventHandler"
	log := h.log.With(slog.String("op", op))

	eventID := c.Params.ByName("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "event_id parameter must be provided"})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	updateRequest := &eventv1.UpdateEventRequest{
		EventId: eventID,
		UserId:  userID,
	}
	paths, err := h.buildUpdatePaths(c, updateRequest)
	if err != nil {
		return
	}

	updateRequest.UpdateMask = &field_mask.FieldMask{Paths: paths}

	// Call UpdateEvent method and handle response
	eventResponse, err := h.eventClient.UpdateEvent(c, updateRequest)
	if err != nil {
		h.handleEventError(c, log, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(eventResponse)})
}

func (h *Handler) DeleteEventHandler(c *gin.Context) {
	const op = "EventHandler.DeleteEventHandler"
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

	canDelete, ok := c.Get("has_roles")
	if !ok {
		log.Warn("role not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, err := h.eventClient.DeleteEvent(c, &eventv1.DeleteEventRequest{
		EventId: eventID,
		UserId:  userID,
		IsAdmin: canDelete.(bool),
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.FailedPrecondition:
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) SendEventForReviewHandler(c *gin.Context) {
	const op = "EventHandler.SendEventForReviewHandler"
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

	res, err := h.eventClient.SendToReview(c, &eventv1.EventActionRequest{
		EventId: eventID,
		UserId:  userID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.FailedPrecondition:
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(res)})
}

func (h *Handler) CancelEventReviewHandler(c *gin.Context) {
	const op = "EventHandler.CancelEventReviewHandler"
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

	res, err := h.eventClient.RevokeReview(c, &eventv1.EventActionRequest{
		EventId: eventID,
		UserId:  userID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.FailedPrecondition:
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(res)})
}

func (h *Handler) ApproveEventHandler(c *gin.Context) {
	const op = "EventHandler.ApproveEventHandler"
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

	user, err := h.getUser(c, userID)
	if err != nil {
		return
	}

	res, err := h.eventClient.ApproveEvent(c, &eventv1.ApproveEventRequest{
		EventId: eventID,
		User:    user,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.FailedPrecondition:
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(res)})
}

func (h *Handler) RejectEventHandler(c *gin.Context) {
	const op = "EventHandler.RejectEventHandler"
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
		Reason string `json:"reason"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.getUser(c, userID)
	if err != nil {
		return
	}

	res, err := h.eventClient.RejectEvent(c, &eventv1.RejectEventRequest{
		EventId: eventID,
		User:    user,
		Reason:  input.Reason,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.FailedPrecondition:
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(res)})
}

func (h *Handler) PublishEventHandler(c *gin.Context) {
	const op = "EventHandler.PublishEventHandler"
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

	event, err := h.eventClient.PublishEvent(c, &eventv1.EventActionRequest{
		EventId: eventID,
		UserId:  userID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.FailedPrecondition:
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(event)})
}

func (h *Handler) UnpublishEventHandler(c *gin.Context) {
	const op = "EventHandler.UnpublishEventHandler"
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

	event, err := h.eventClient.UnpublishEvent(c, &eventv1.EventActionRequest{
		EventId: eventID,
		UserId:  userID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.FailedPrecondition:
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": domain.ProtoToEvent(event)})
}
