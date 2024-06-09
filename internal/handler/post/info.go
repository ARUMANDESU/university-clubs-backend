package posthandler

import (
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handler) GetPostHandler(c *gin.Context) {
	const op = "handler.post.getPostHandler"
	log := h.log.With(slog.String("op", op))

	postID := c.Params.ByName("id")
	if postID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "post id parameter must be provided"})
		return
	}

	post, err := h.postsClient.GetPost(c, &postv1.GetPostRequest{
		Id:     postID,
		UserId: 0,
	})
	if err != nil {
		handleErrors(c, log, "failed to update post", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": domain.PostFromPb(post)})
}

func (h *Handler) ListPublishedPostsHandler(c *gin.Context) {
	const op = "handler.post.listPublishedPostsHandler"
	log := h.log.With(slog.String("op", op))

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request := &postv1.ListPostsRequest{
		Page:      int32(page),
		PageSize:  int32(pageSize),
		Query:     c.Query("query"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	tags := strings.Split(c.Query("tags"), ",")
	if len(tags) > 0 && tags[0] != "" {
		request.Tags = tags
	}

	res, err := h.postsClient.ListPosts(c, request)
	if err != nil {
		handleErrors(c, log, "failed to list published posts", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": domain.PostsFromPb(res.GetPosts()), "metadata": res.GetMetadata()})
}

func (h *Handler) ListClubPostsHandler(c *gin.Context) {
	const op = "handler.post.listClubPostsHandler"
	log := h.log.With(slog.String("op", op))

	clubID, err := utils.GetIntFromParams(c.Params, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request := &postv1.ListPostsRequest{
		ClubId:    clubID,
		Page:      int32(page),
		PageSize:  int32(pageSize),
		Query:     c.Query("query"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	tags := strings.Split(c.Query("tags"), ",")
	if len(tags) > 0 && tags[0] != "" {
		request.Tags = tags
	}

	res, err := h.postsClient.ListPosts(c, request)
	if err != nil {
		handleErrors(c, log, "failed to list club posts", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": domain.PostsFromPb(res.GetPosts()), "metadata": res.GetMetadata()})
}
