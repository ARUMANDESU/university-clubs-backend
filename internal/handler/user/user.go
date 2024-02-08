package user

import (
	"bytes"
	"fmt"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) GetUser(c *gin.Context) {
	const op = "UserHandler.GetUser"
	log := h.log.With(slog.String("op", op))

	userID, err := getIntFromParams(c.Params, "id")
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
		AvatarURL: res.GetAvatarUrl(),
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

	userID, err := getIntFromParams(c.Params, "id")
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

	if userID != userIDFromCtx.(int64) {
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

	userID, err := getIntFromParams(c.Params, "id")
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

	if userID != userIDFromCtx.(int64) {
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

func (h *Handler) SearchUsers(c *gin.Context) {
	const op = "UserHandler.SearchUsers"
	log := h.log.With(slog.String("op", op))

	query := c.Query("query")
	page, err := getIntFromQuery(c, "page")
	if err != nil {
		log.Warn("failed to get page query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := getIntFromQuery(c, "page_size")
	if err != nil {
		log.Warn("failed to get page_size query parameter", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usrClient.SearchUsers(c, &userv1.SearchUsersRequest{
		Query:      query,
		PageNumber: int32(page),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	users := domain.MapUserObjectArrToDomain(res.Users)

	c.JSON(http.StatusOK, gin.H{"users": users, "metadata": res.Metadata})
}

func (h *Handler) UpdateAvatar(c *gin.Context) {
	const op = "UserHandler.UpdateAvatar"
	log := h.log.With(slog.String("op", op))

	userID, err := getIntFromParams(c.Params, "id")
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

	if userID != userIDFromCtx.(int64) {
		log.Warn("not account owner")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		log.Error("failed to get image file from form", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Error("failed to open file", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}

	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Error("failed to copy image into bytes", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = h.usrClient.UpdateAvatar(c, &userv1.UpdateAvatarRequest{
		UserId: userID,
		Image:  buf.Bytes(),
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)

}

func getIntFromParams(c gin.Params, param string) (int64, error) {
	p := c.ByName(param)
	if p == "" {
		return 0, fmt.Errorf("%s parameter must be provided", param)
	}

	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s parameter must be an integer", param)
	}

	return i, nil
}

func getIntFromQuery(c *gin.Context, query string) (int, error) {
	q, ok := c.GetQuery(query)
	if !ok {
		return 0, fmt.Errorf("%s query parameter must be provided", query)
	}
	res, err := strconv.Atoi(q)
	if err != nil {
		return 0, fmt.Errorf("%s query parameter must be an integer", query)
	}

	return res, nil
}
