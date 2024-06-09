package posthandler

import (
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/genproto/protobuf/field_mask"
	"log/slog"
	"net/http"
)

func (h *Handler) CreatePostHandler(c *gin.Context) {
	const op = "handler.post.createPostHandler"
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

	var input struct {
		Title         string              `json:"title" binding:"required"`
		Description   string              `json:"description,omitempty"`
		Tags          []string            `json:"tags,omitempty"`
		CoverImages   []domain.CoverImage `json:"cover_images,omitempty"`
		AttachedFiles []domain.EventFile  `json:"attached_files,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var paths []string
	if input.Title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	} else {
		paths = append(paths, "title")
	}
	if input.Description != "" {
		paths = append(paths, "description")
	}
	if len(input.Tags) > 0 {
		paths = append(paths, "tags")
	}
	if len(input.CoverImages) > 0 {
		paths = append(paths, "cover_images")
	}
	if len(input.AttachedFiles) > 0 {
		paths = append(paths, "attached_files")
	}

	post, err := h.postsClient.CreatePost(c, &postv1.CreatePostRequest{
		ClubId:        clubID,
		UserId:        userID,
		Title:         input.Title,
		Description:   input.Description,
		Tags:          input.Tags,
		CoverImages:   domain.CoverImageToProtoArr(input.CoverImages),
		AttachedFiles: domain.EventFileToProtoArr(input.AttachedFiles),
		CreateMask:    &field_mask.FieldMask{Paths: paths},
	})
	if err != nil {
		handleErrors(c, log, "failed to create post", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"post": domain.PostFromPb(post)})
}

func (h *Handler) UpdatePostHandler(c *gin.Context) {
	const op = "handler.post.updatePostHandler"
	log := h.log.With(slog.String("op", op))

	postID := c.Params.ByName("id")
	if postID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "post id parameter must be provided"})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	var input struct {
		Title         string              `json:"title" binding:"required"`
		Description   string              `json:"description,omitempty"`
		Tags          []string            `json:"tags,omitempty"`
		CoverImages   []domain.CoverImage `json:"cover_images,omitempty"`
		AttachedFiles []domain.EventFile  `json:"attached_files,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var paths []string
	request := &postv1.UpdatePostRequest{
		Id:     postID,
		UserId: userID,
	}

	if input.Title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	} else {
		paths = append(paths, "title")
		request.Title = input.Title
	}
	if input.Description != "" {
		paths = append(paths, "description")
		request.Description = input.Description
	}
	if len(input.Tags) > 0 {
		paths = append(paths, "tags")
		request.Tags = input.Tags
	}
	if len(input.CoverImages) > 0 {
		paths = append(paths, "cover_images")
		request.CoverImages = domain.CoverImageToProtoArr(input.CoverImages)
	}
	if len(input.AttachedFiles) > 0 {
		paths = append(paths, "attached_files")
		request.AttachedFiles = domain.EventFileToProtoArr(input.AttachedFiles)
	}

	request.UpdateMask = &field_mask.FieldMask{Paths: paths}

	post, err := h.postsClient.UpdatePost(c, request)
	if err != nil {
		handleErrors(c, log, "failed to update post", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": domain.PostFromPb(post)})
}

func (h *Handler) DeletePostHandler(c *gin.Context) {
	const op = "handler.post.deletePostHandler"
	log := h.log.With(slog.String("op", op))

	postID := c.Params.ByName("id")
	if postID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "post id parameter must be provided"})
		return
	}

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := userIDFromCtx.(int64)

	post, err := h.postsClient.DeletePost(c, &postv1.ActionRequest{
		Id:     postID,
		UserId: userID,
	})
	if err != nil {
		handleErrors(c, log, "failed to delete post", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": domain.PostFromPb(post)})
}
