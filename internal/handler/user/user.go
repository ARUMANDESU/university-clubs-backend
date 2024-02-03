package user

import (
	"errors"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) GetUser(c *gin.Context) {
	const op = "UserHandler.GetUser"
	log := h.log.With(slog.String("op", op))

	userID, err := getUserIdFromParams(c.Params)
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (h *Handler) UpdateUser(c *gin.Context) {
	const op = "UserHandler.UpdateUser"
	log := h.log.With(slog.String("op", op))

	userID, err := getUserIdFromParams(c.Params)
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDFronCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if userID != userIDFronCtx.(int64) {
		log.Warn("not account owner")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var input struct {
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		Major     string `json:"major,omitempty"`
		GroupName string `json:"group_name,omitempty"`
		Year      int    `json:"year,omitempty"`
	}

	err = c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var paths []string
	if input.FirstName != "" {
		paths = append(paths, "first_name")
	}
	if input.LastName != "" {
		paths = append(paths, "last_name")
	}
	if input.Major != "" {
		paths = append(paths, "major")
	}
	if input.GroupName != "" {
		paths = append(paths, "group_name")
	}
	if input.Year != 0 {
		paths = append(paths, "year")
	}

	res, err := h.usrClient.UpdateUser(c, &userv1.UpdateUserRequest{
		UserId:     userID,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Major:      input.Major,
		GroupName:  input.GroupName,
		Year:       int32(input.Year),
		UpdateMask: &fieldmaskpb.FieldMask{Paths: paths},
	})
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

	c.JSON(http.StatusOK, gin.H{"user_id": res.GetUserId()})

}

func (h *Handler) DeleteUser(c *gin.Context) {
	const op = "UserHandler.DeleteUser"
	log := h.log.With(slog.String("op", op))

	userID, err := getUserIdFromParams(c.Params)
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDFronCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if userID != userIDFronCtx.(int64) {
		log.Warn("not account owner")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	_, err = h.usrClient.DeleteUser(c, &userv1.DeleteUserRequest{UserId: userID})
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

	c.Status(http.StatusOK)
}

func getUserIdFromParams(p gin.Params) (int64, error) {
	id := p.ByName("id")
	if id == "" {
		return 0, errors.New("'id' parameter must be provided")
	}

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, errors.New("'id' parameter must be integer")
	}

	return userID, nil
}
