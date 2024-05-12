package club

import (
	"context"
	"errors"
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	imageUtils "github.com/ARUMANDESU/university-clubs-backend/pkg/image"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"log/slog"
	"net/http"
	"path"
	"time"
)

const clubBucket = "ucms-club-images-dev"

func (h *Handler) CreateClubHandler(c *gin.Context) {
	const op = "ClubHandler.CreateClubHandler"
	log := h.log.With(slog.String("op", op))

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID := userIDFromCtx.(int64)

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ClubType    string `json:"club_type"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.clubClient.CreateClub(c, &clubv1.CreateClubRequest{
		Name:        input.Name,
		Description: input.Description,
		ClubType:    input.ClubType,
		OwnerId:     userID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			log.Warn("club not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusCreated)

}

func (h *Handler) NewClubHandler(c *gin.Context) {
	const op = "ClubHandler.NewClubHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input struct {
		Status string `json:"status"`
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action := clubv1.HandleClubAction_REJECT
	if input.Status == "approved" {
		action = clubv1.HandleClubAction_APPROVE
	}

	_, err = h.clubClient.HandleNewClub(c, &clubv1.HandleNewClubRequest{
		ClubId: clubID,
		Action: action,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			log.Warn("club not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusCreated)

}

func (h *Handler) HandleJoinRequestHandler(c *gin.Context) {
	const op = "ClubHandler.HandleJoinRequestHandler"
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

	var input struct {
		TargetID int64  `json:"user_id"`
		Status   string `json:"status"`
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action := clubv1.HandleClubAction_REJECT
	if input.Status == "approved" {
		action = clubv1.HandleClubAction_APPROVE
	}

	_, err = h.clubClient.HandleJoinClub(c, &clubv1.HandleJoinClubRequest{
		ClubId:   clubID,
		UserId:   input.TargetID,
		MemberId: userID,
		Action:   action,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			log.Warn("club not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) JoinRequestHandler(c *gin.Context) {
	const op = "ClubHandler.JoinRequestHandler"
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

	_, err = h.clubClient.RequestToJoinClub(c, &clubv1.RequestToJoinClubRequest{
		UserId: userID,
		ClubId: clubID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.AlreadyExists:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusCreated)

}

func (h *Handler) UpdateLogoHandler(c *gin.Context) {
	const op = "ClubHandler.UpdateLogoHandler"
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

	fileBytes, fileSize, err := utils.GetFileByName(c, "logo")
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

	if fileSize > 5*1024*1024 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Image size should be less than 5MB"})
		return
	}

	compressImage, filename, err := imageUtils.CompressImage(fileBytes, 75)
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

	url, err := h.imageStorage.UploadImage(imageCtx, compressImage, filename, clubBucket)
	if err != nil {
		log.Error("failed to upload avatar", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res, err := h.clubClient.UpdateLogo(c, &clubv1.UpdateLogoRequest{
		LogoUrl: url,
		UserId:  userID,
		ClubId:  clubID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	if res.GetPrevLogoUrl() != "" {
		go func() {
			deleteCtx, deleteCtxCancel := context.WithTimeout(c, time.Second*45)
			defer deleteCtxCancel()
			err := h.imageStorage.DeleteImage(deleteCtx, path.Base(res.GetPrevLogoUrl()), clubBucket)
			if err != nil {
				log.Error("failed to delete previous logo", logger.Err(err))
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"club": domain.ClubObjectToClub(res.GetClub())})

}

func (h *Handler) UpdateBannerHandler(c *gin.Context) {
	const op = "ClubHandler.UpdateBannerHandler"
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

	fileBytes, fileSize, err := utils.GetFileByName(c, "banner")
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

	if fileSize > 5*1024*1024 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Image size should be less than 5MB"})
		return
	}

	compressImage, filename, err := imageUtils.CompressImage(fileBytes, 75)
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

	url, err := h.imageStorage.UploadImage(imageCtx, compressImage, filename, clubBucket)
	if err != nil {
		log.Error("failed to upload avatar", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res, err := h.clubClient.UpdateBanner(c, &clubv1.UpdateBannerRequest{
		BannerUrl: url,
		UserId:    userID,
		ClubId:    clubID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	if res.GetPrevBannerUrl() != "" {
		go func() {
			deleteCtx, deleteCtxCancel := context.WithTimeout(c, time.Second*45)
			defer deleteCtxCancel()
			err := h.imageStorage.DeleteImage(deleteCtx, path.Base(res.GetPrevBannerUrl()), clubBucket)
			if err != nil {
				log.Error("failed to delete previous banner", logger.Err(err))
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"club": domain.ClubObjectToClub(res.GetClub())})

}

func (h *Handler) UpdateClubHandler(c *gin.Context) {
	const op = "ClubHandler.UpdateBannerHandler"
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

	var input struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		ClubType    string `json:"club_type,omitempty"`
	}

	err = c.ShouldBindJSON(&input)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var paths []string
	if input.Name != "" {
		paths = append(paths, "name")
	}
	if input.Description != "" {
		paths = append(paths, "description")
	}
	if input.ClubType != "" {
		paths = append(paths, "club_type")
	}

	club, err := h.clubClient.UpdateClub(c, &clubv1.UpdateClubRequest{
		ClubId:      clubID,
		UserId:      userID,
		Name:        input.Name,
		Description: input.Description,
		ClubType:    input.ClubType,
		UpdateMask:  &fieldmaskpb.FieldMask{Paths: paths},
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.PermissionDenied:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"club": domain.ClubObjectToClub(club)})

}

func (h *Handler) LeaveClubHandler(c *gin.Context) {
	const op = "ClubHandler.LeaveClubHandler"
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

	_, err = h.clubClient.LeaveClub(c, &clubv1.LeaveClubRequest{
		UserId: userID,
		ClubId: clubID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusNoContent)

}

func (h *Handler) KickMemberHandler(c *gin.Context) {
	const op = "ClubHandler.KickMemberHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memberID, err := utils.GetIntFromParams(c.Params, "member_id")
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

	_, err = h.clubClient.KickMemberFromClub(c, &clubv1.KickMemberFromClubRequest{
		ClubId:   clubID,
		UserId:   userID,
		TargetId: memberID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.Aborted:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
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

func (h *Handler) TransferOwnershipHandler(c *gin.Context) {
	const op = "ClubHandler.TransferOwnershipHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		log.Warn("failed to get id params", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newOwnerID, err := utils.GetIntFromParams(c.Params, "member_id")
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

	_, err = h.clubClient.TransferOwnership(c, &clubv1.TransferOwnershipRequest{
		ClubId:   clubID,
		UserId:   userID,
		TargetId: newOwnerID,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.AlreadyExists:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
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
