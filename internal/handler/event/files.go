package eventhandler

import (
	"errors"
	"fmt"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"path"
)

const eventBucket = "ucms-posts-files-dev"
const maxFileSize = 55 * 1024 * 1024  // 55MB
const maxImageSize = 15 * 1024 * 1024 // 15MB

func (h *Handler) UploadFileHandler(c *gin.Context) {
	const op = "handler.event.UploadFileHandler"
	log := h.log.With(slog.String("op", op))

	file, err := utils.GetFileByName(c, "file")
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidFileUpload):
			log.Error("failed to get file from form", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		case errors.Is(err, utils.ErrConvFileToBytes):
			log.Error("failed to copy file into bytes", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		default:
			log.Error("failed to get bytes from file", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size should be less than 55MB"})
		return
	}

	filename := uuid.New().String()

	url, err := h.FileStorage.Upload(c, file.Bytes, filename, eventBucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *Handler) UploadFilesHandler(c *gin.Context) {
	const op = "handler.event.UploadFilesHandler"
	log := h.log.With(slog.String("op", op))

	files, err := utils.GetFilesByName(c, "files", 5)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrMaxFilesCount):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, utils.ErrInvalidFileUpload):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		case errors.Is(err, utils.ErrConvFileToBytes):
			c.AbortWithStatus(http.StatusInternalServerError)
		default:
			log.Error("failed to get bytes from file", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	type FileResult struct {
		File domain.EventFile
		Err  error
	}

	filesData := make([]domain.EventFile, 0, len(files))
	errs := make([]string, 0)

	results := make(chan FileResult, len(files))

	for _, file := range files {
		if file.Size > maxFileSize {
			errs = append(errs, fmt.Sprintf("File %s size should be less than 55MB", file.Name))
			continue
		}

		go func(file domain.File) {
			filename := uuid.New().String()
			url, err := h.FileStorage.Upload(c, file.Bytes, filename, eventBucket)
			if err != nil {
				results <- FileResult{Err: fmt.Errorf("failed to upload file %s", file.Name)}
				return
			}

			results <- FileResult{
				File: domain.EventFile{
					Name: file.Name,
					Type: file.Type,
					Url:  url,
				},
			}
		}(file)
	}

	for range files {
		result := <-results
		if result.Err != nil {
			errs = append(errs, result.Err.Error())
			continue
		}
		filesData = append(filesData, result.File)
	}

	close(results)

	c.JSON(http.StatusOK, gin.H{"files": filesData, "errors": errs})
}

func (h *Handler) UploadImagesHandler(c *gin.Context) {
	const op = "handler.event.UploadImagesHandler"
	log := h.log.With(slog.String("op", op))

	files, err := utils.GetFilesByName(c, "images", 5)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrMaxFilesCount):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	var errs []string
	images := make([]domain.EventFile, 0, len(files))

	for _, file := range files {
		if file.Size > maxImageSize {
			errs = append(errs, fmt.Sprintf("Image %s size should be less than 15MB", file.Name))
		}

		filename := uuid.New().String()

		url, err := h.imageStorage.Upload(c, file.Bytes, filename, eventBucket)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}

		images = append(images, domain.EventFile{
			Name: file.Name,
			Type: file.Type,
			Url:  url,
		})

	}

	c.JSON(http.StatusOK, gin.H{"images": images, "errors": errs})
}

func (h *Handler) UploadImageHandler(c *gin.Context) {
	const op = "handler.event.UploadImageHandler"
	log := h.log.With(slog.String("op", op))

	file, err := utils.GetFileByName(c, "image")
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image size should be less than 15MB"})
		return
	}

	filename := uuid.New().String()

	url, err := h.imageStorage.Upload(c, file.Bytes, filename, eventBucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *Handler) DeleteFileHandler(c *gin.Context) {
	const op = "handler.event.DeleteFileHandler"
	log := h.log.With(slog.String("op", op))

	url := c.Query("url")
	if url == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "url parameter must be provided"})
		return
	}

	err := h.FileStorage.Delete(c, path.Base(url), eventBucket)
	if err != nil {
		log.Error("failed to delete file", logger.Err(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func (h *Handler) DeleteFileOldHandler(c *gin.Context) {
	const op = "handler.event.DeleteFileOldHandler"
	log := h.log.With(slog.String("op", op))

	var input struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error("failed to bind json", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if input.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url must be provided"})
		return
	}

	err := h.FileStorage.Delete(c, path.Base(input.URL), eventBucket)
	if err != nil {
		log.Error("failed to delete file", logger.Err(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
