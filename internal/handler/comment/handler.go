package commenthandler

import (
	"errors"
	"log/slog"
	"net/http"

	commentv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/comment"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/comment"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	commentClient *comment.Client
	log           *slog.Logger
}

func New(log *slog.Logger, client *comment.Client) Handler {
	return Handler{
		commentClient: client,
		log:           log,
	}
}

func (h *Handler) GetCommentByID(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(400, gin.H{"error": "comment id is required"})
		return
	}

	response, err := h.commentClient.GetCommentByID(c, &commentv1.GetCommentByIDRequest{Id: commentID})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"comment": response.Comment})
}

func (h *Handler) ListPostComments(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(400, gin.H{"error": "post id is required"})
		return
	}

	requestBody := &commentv1.ListPostCommentsRequest{PostId: postID}

	page, err := utils.GetIntFromQuery(c, "page")
	if err != nil && !errors.Is(err, utils.ErrMustbeProvided) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestBody.Page = int32(page)

	pageSize, err := utils.GetIntFromQuery(c, "page_size")
	if err != nil && !errors.Is(err, utils.ErrMustbeProvided) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestBody.PageSize = int32(pageSize)

	sortBy := c.Query("sort_by")
	switch sortBy {
	case "created_at":
		requestBody.SortBy = commentv1.SortBy_SORT_BY_CREATED_AT
	case "updated_at":
		requestBody.SortBy = commentv1.SortBy_SORT_BY_UPDATED_AT
	default:
		requestBody.SortBy = commentv1.SortBy_SORT_BY_UNSPECIFIED
	}

	sortOrder := c.Query("sort_order")
	switch sortOrder {
	case "asc":
		requestBody.SortOrder = commentv1.SortOrder_SORT_ORDER_ASC
	case "desc":
		requestBody.SortOrder = commentv1.SortOrder_SORT_ORDER_DESC
	default:
		requestBody.SortOrder = commentv1.SortOrder_SORT_ORDER_UNSPECIFIED
	}

	response, err := h.commentClient.ListPostComments(c, requestBody)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": response.Comments, "metadata": response.Metadata})
}

func handleError(c *gin.Context, err error) {
	switch {
	case status.Code(err) == codes.InvalidArgument:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
	case status.Code(err) == codes.NotFound:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
	case status.Code(err) == codes.PermissionDenied:
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
	default:
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
