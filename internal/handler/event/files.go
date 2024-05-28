package event

import (
	"errors"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

const eventBucket = "ucms-posts-files-dev"
const maxFileSize = 25 * 1024 * 1024 // 25MB
const maxImageSize = 8 * 1024 * 1024 // 8MB

func (h *Handler) UploadFilesHandler(c *gin.Context) {

	const op = "handler.event.UploadFilesHandler"
	log := h.log.With(slog.String("op", op))

	file, err := utils.GetFileByName(c, "file")

	filename := uuid.New().String()

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

	if file.Size > maxFileSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Image size should be less than 15MB"})
		return
	}

	url, err := h.FileStorage.Upload(c, file.Bytes, filename, eventBucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": file.Name, "type": file.Type, "url": url})
}

func (h *Handler) UploadImagesHandler(c *gin.Context) {

	const op = "handler.event.UploadImagesHandler"
	log := h.log.With(slog.String("op", op))

	file, err := utils.GetFileByName(c, "image")

	filename := uuid.New().String()

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

	if file.Size > maxImageSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Image size should be less than 15MB"})
		return
	}

	url, err := h.imageStorage.Upload(c, file.Bytes, filename, eventBucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": file.Name, "type": file.Type, "url": url})
}
