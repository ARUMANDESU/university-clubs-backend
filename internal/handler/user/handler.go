package user

import (
	"context"
	"errors"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/config"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/jwt"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strings"
)

type Handler struct {
	usrClient     *user.Client
	imageStorage  ImageStorage
	log           *slog.Logger
	confClient    *confidential.Client
	jwtSecret     string
	MicrosoftOIDC config.MicrosoftOIDC
}

type ImageStorage interface {
	UploadImage(ctx context.Context, image []byte, filename string, bucket string) (string, error)
	DeleteImage(ctx context.Context, filename string, bucket string) error
}

// New creates and returns a new User Handler instance
// Parameters:
//   - client: A *user.Client which is a gRPC client for the user service.
//   - log: A *slog.Logger used for logging messages and errors.
//
// Returns:
//   - A Handler struct that encapsulates the provided user service client and logger.
func New(
	cfg *config.Config,
	log *slog.Logger,
	client *user.Client,
	confClient confidential.Client,
	imageStorage ImageStorage,
) Handler {

	return Handler{
		usrClient:     client,
		log:           log,
		jwtSecret:     cfg.JwtSecret,
		confClient:    &confClient,
		MicrosoftOIDC: cfg.MicrosoftOIDC,
		imageStorage:  imageStorage,
	}
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	const op = "handler.user.authMiddleware"
	log := h.log.With(slog.String("op", op))

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		// Split the authorization header to retrieve the token part
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			return
		}

		jwtToken := authParts[1]

		userID, err := jwt.GetUserID(jwtToken, h.jwtSecret)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrTokenIsNotValid),
				errors.Is(err, domain.ErrInvalidTokenClaims),
				errors.Is(err, domain.ErrUserIDClaimNotFound):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.Is(err, domain.ErrTokenIsExpired):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			default:
				log.Error("failed to get user id from jwt token", logger.Err(err))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func (h *Handler) RoleAuthMiddleware(roles []userv1.Role) gin.HandlerFunc {
	const op = "handler.user.roleAuthMiddleware"

	log := h.log.With(slog.String("op", op))

	return func(c *gin.Context) {

		userID, ok := c.Get("userID")
		if !ok {
			log.Warn("userID not found")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		/*		uid, err := strconv.ParseInt(userID.(int64), 10, 64)
				if err != nil {
					log.Error("userID cannot convert into int64", logger.Err(err))
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}*/

		res, err := h.usrClient.CheckUserRole(c, &userv1.CheckUserRoleRequest{UserId: userID.(int64), Roles: roles})
		if err != nil {
			switch {
			case status.Code(err) == codes.InvalidArgument:
				log.Warn("invalid arguments", logger.Err(err))
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
			case status.Code(err) == codes.NotFound:
				log.Warn("session not found", logger.Err(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
			default:
				log.Error("internal", logger.Err(err))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		if !res.GetHasRole() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()

	}
}
