package user

import (
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) GetUser(c *gin.Context) {
	const op = "UserHandler.GetUser"
	log := h.log.With(slog.String("op", op))

	id := c.Param("id")
	if id == "" {
		log.Warn("id not provided")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "'id' parameter must be provided"})
		return
	}

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Warn("failed to parse into int64")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "'id' parameter must be integer"})
		return
	}

	res, err := h.usrClient.GetUser(c, &userv1.GetUserRequest{UserId: userID})
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

	user := &domain.User{
		ID:        res.GetUserId(),
		FirstName: res.GetFirstName(),
		LastName:  res.GetLastName(),
		Email:     res.GetEmail(),
		CreatedAt: res.GetCreatedAt().AsTime(),
		Role:      res.GetRole().String(),
		Barcode:   res.GetBarcode(),
		Major:     res.GetMajor(),
		GroupName: res.GetGroupName(),
		Year:      int(res.GetYear()),
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
