package posthandler

import (
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/post"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

type Handler struct {
	postsClient *post.Client
	log         *slog.Logger
}

func New(
	log *slog.Logger,
	postsClient *post.Client,
) Handler {
	return Handler{
		postsClient: postsClient,
		log:         log,
	}
}

func handleErrors(c *gin.Context, log *slog.Logger, msg string, err error) {
	switch {
	case status.Code(err) == codes.InvalidArgument:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
	case status.Code(err) == codes.NotFound:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
	case status.Code(err) == codes.PermissionDenied:
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": status.Convert(err).Message()})
	case status.Code(err) == codes.AlreadyExists:
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
	default:
		log.Error(msg, logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
