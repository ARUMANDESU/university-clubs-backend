package event

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

const eventBucket = "ucms-posts-files-dev"
const maxFileSize = 15 * 1024 * 1024 // 15MB

func (h *Handler) UploadFilesHandler(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
		return
	}

	if len(fileBytes) > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the 15MB limit"})
		return
	}

	url, err := h.s3Client.UploadImage(c, fileBytes, fileHeader.Filename, eventBucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": fileHeader.Filename, "url": url, "type": fileHeader.Header.Get("Content-Type")})
}

func (h *Handler) UploadImagesHandler(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
		return
	}

	if len(fileBytes) > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the 15MB limit"})
		return
	}

	url, err := h.s3Client.UploadImage(c, fileBytes, fileHeader.Filename, eventBucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": fileHeader.Filename, "url": url, "type": fileHeader.Header.Get("Content-Type")})
}
