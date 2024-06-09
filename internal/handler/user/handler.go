package userhandler

import (
	"context"
	"errors"
	"fmt"
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
)

type Handler struct {
	log           *slog.Logger
	cfg           *config.Config
	usrClient     *user.Client
	imageStorage  ImageStorage
	confClient    *confidential.Client
	MicrosoftOIDC config.MicrosoftOIDC
}

type ImageStorage interface {
	Upload(ctx context.Context, image []byte, filename string, bucket string) (string, error)
	Delete(ctx context.Context, filename string, bucket string) error
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
		cfg:           cfg,
		confClient:    &confClient,
		MicrosoftOIDC: cfg.MicrosoftOIDC,
		imageStorage:  imageStorage,
	}
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	const op = "handler.user.authMiddleware"
	log := h.log.With(slog.String("op", op))

	return func(c *gin.Context) {
		accessToken, err := c.Cookie(AccessTokenName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s cookie not found", AccessTokenName)})
			return
		}

		userID, err := jwt.GetUserID(accessToken, h.cfg.JwtSecret)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrTokenIsNotValid),
				errors.Is(err, domain.ErrInvalidTokenClaims),
				errors.Is(err, domain.ErrUserIDClaimNotFound),
				errors.Is(err, domain.ErrTokenIsExpired),
				errors.Is(err, domain.ErrUnexpectedSigningMethod):
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

func (h *Handler) GetUserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie(AccessTokenName)
		if err != nil {
			c.Set("userID", int64(0))
		} else {
			h.AuthMiddleware()(c)
		}

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

		res, err := h.usrClient.CheckUserRole(c, &userv1.CheckUserRoleRequest{UserId: userID.(int64), Roles: roles})
		if err != nil {
			switch {
			case status.Code(err) == codes.InvalidArgument:
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
			case status.Code(err) == codes.NotFound:
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

func (h *Handler) HasRoles(roles []userv1.Role) gin.HandlerFunc {
	const op = "handler.user.hasRoles"
	log := h.log.With(slog.String("op", op))

	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			log.Warn("userID not found")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		res, err := h.usrClient.CheckUserRole(c, &userv1.CheckUserRoleRequest{UserId: userID.(int64), Roles: roles})
		if err != nil {
			switch {
			case status.Code(err) == codes.InvalidArgument:
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
			case status.Code(err) == codes.NotFound:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": status.Convert(err).Message()})
			default:
				log.Error("internal", logger.Err(err))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		c.Set("has_roles", res.GetHasRole())

		c.Next()
	}
}
