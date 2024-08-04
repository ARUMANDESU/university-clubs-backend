package userhandler

import (
	"context"
	"errors"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	imageUtils "github.com/ARUMANDESU/university-clubs-backend/pkg/image"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log/slog"
	"net/http"
	"path"
	"time"
)

const userBucket = "ucms-user-profile-images-dev"

func (h *Handler) GetUser(c *gin.Context) {
	const op = "UserHandler.GetUser"
	log := h.log.With(slog.String("op", op))

	userID, err := utils.GetIntFromParams(c.Params, "id")
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

	user := domain.UserObjectToDomain(res)

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	const op = "UserHandler.UpdateUser"
	log := h.log.With(slog.String("op", op))

	userID, err := utils.GetIntFromParams(c.Params, "id")
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

	c.JSON(http.StatusOK, gin.H{"user": domain.UserObjectToDomain(res)})

}

func (h *Handler) DeleteUser(c *gin.Context) {
	const op = "UserHandler.DeleteUser"
	log := h.log.With(slog.String("op", op))

	userID, err := utils.GetIntFromParams(c.Params, "id")
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

	userID, err := utils.GetIntFromParams(c.Params, "id")
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

	file, err := utils.GetFileByName(c, "avatar")
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidFileUpload):
			log.Error("failed to get image file from form", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		case errors.Is(err, utils.ErrConvFileToBytes):
			log.Error("failed to copy image into bytes", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		default:
			log.Error("failed to get bytes from file", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	if file.Size > 5*1024*1024 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Image size should be less than 5MB"})
		return
	}

	compressImage, filename, err := imageUtils.CompressImage(file.Bytes, 75)
	if err != nil {
		switch {
		case errors.Is(err, imageUtils.ErrImageQuality),
			errors.Is(err, imageUtils.ErrImageFormat),
			errors.Is(err, imageUtils.ErrImageIsEmpty):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Error("failed to compress image", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	imageCtx, cancel := context.WithTimeout(c, time.Second*20)
	defer cancel()

	url, err := h.imageStorage.Upload(imageCtx, compressImage, filename, userBucket)
	if err != nil {
		log.Error("failed to upload avatar", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res, err := h.usrClient.UpdateAvatar(c, &userv1.UpdateAvatarRequest{
		UserId:   userID,
		ImageUrl: url,
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

	if res.GetPrevAvatarUrl() != "" {
		go func() {
			deleteCtx, deleteCtxCancel := context.WithTimeout(c, time.Second*45)
			defer deleteCtxCancel()
			err := h.imageStorage.Delete(deleteCtx, path.Base(res.GetPrevAvatarUrl()), userBucket)
			if err != nil {
				log.Error("failed to delete previous avatar", logger.Err(err))
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"user": domain.UserObjectToDomain(res.GetUser())})

}

func (h *Handler) ChangeUserRole(c *gin.Context) {
	const op = "UserHandler.ChangeUserRole"
	log := h.log.With(slog.String("op", op))

	userID, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	targetID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := c.Query("role")
	err = validation.Validate(role, validation.In("DSVR", "ADMIN", "MODER", "USER"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, err = h.usrClient.ChangeUserRole(c, &userv1.ChangeUserRoleRequest{
		UserId:   userID.(int64),
		TargetId: targetID,
		Role:     domain.MapRoleStringToEnum(role),
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	const op = "handler.user.change_password"
	log := h.log.With(slog.String("op", op))

	userID, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.usrClient.ChangeUserPassword(c, &userv1.ChangeUserPasswordRequest{
		UserId:      userID.(int64),
		OldPassword: input.OldPassword,
		NewPassword: input.NewPassword,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	const op = "handler.user.forgot_password"
	log := h.log.With(slog.String("op", op))

	var input struct {
		Email   string `json:"email" binding:"required,email"`
		Barcode string `json:"barcode" binding:"required"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.usrClient.ForgotPassword(c, &userv1.ForgotPasswordRequest{
		Email:   input.Email,
		Barcode: input.Barcode,
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

	c.Status(http.StatusNoContent)
}

func (h *Handler) ResetPassword(c *gin.Context) {
	const op = "handler.user.reset_password"
	log := h.log.With(slog.String("op", op))

	token := c.Query("token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	var input struct {
		NewPassword string `json:"new_password" binding:"required"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.usrClient.ResetPassword(c, &userv1.ResetPasswordRequest{
		VerificationToken: token,
		NewPassword:       input.NewPassword,
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

	c.Status(http.StatusNoContent)
}
