package user

import (
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
	"strings"
)

type Handler struct {
	usrClient           *user.Client
	log                 *slog.Logger
	confClient          *confidential.Client
	notificationService NotificationService
	jwtSecret           string
	config.MicrosoftOIDC
}

type NotificationService interface {
	AddNewConnection(userId int64) chan domain.Notification
	RemoveConnection(userId int64)
	GetConnection(userId int64) (<-chan domain.Notification, error)
	//HandleNotification()
}

// New creates and returns a new User Handler instance
// Parameters:
//   - client: A *user.Client which is a gRPC client for the user service.
//   - log: A *slog.Logger used for logging messages and errors.
//
// Returns:
//   - A Handler struct that encapsulates the provided user service client and logger.
func New(
	client *user.Client,
	log *slog.Logger,
	cfg *config.Config,
	confClient confidential.Client,
	notificationService NotificationService,
) Handler {

	return Handler{
		usrClient:           client,
		log:                 log,
		jwtSecret:           cfg.JwtSecret,
		confClient:          &confClient,
		MicrosoftOIDC:       cfg.MicrosoftOIDC,
		notificationService: notificationService,
	}
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	const op = "AuthMiddleware"
	log := h.log.With(slog.String("op", op))

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		// Split the authorization header to retrieve the token part
		authParts := strings.Split(authHeader, " ")
		log.Info(fmt.Sprintf("%v", authParts))
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			return
		}

		jwtToken := authParts[1]
		log.Debug(jwtToken)

		userID, expired, err := jwt.GetUserID(jwtToken, h.jwtSecret)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrTokenIsNotValid),
				errors.Is(err, domain.ErrInvalidTokenClaims),
				errors.Is(err, domain.ErrUserIDClaimNotFound):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				log.Error("failed to get user id from jwt token", logger.Err(err))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		if expired {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func (h *Handler) RoleAuthMiddleware(roles []userv1.Role) gin.HandlerFunc {
	const op = "RoleAuthMiddleware"

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
