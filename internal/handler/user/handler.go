package user

import (
	"fmt"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

type Handler struct {
	usrClient *user.Client
	log       *slog.Logger
}

// New creates and returns a new User Handler instance
// Parameters:
//   - client: A *user.Client which is a gRPC client for the user service.
//   - log: A *slog.Logger used for logging messages and errors.
//
// Returns:
//   - A Handler struct that encapsulates the provided user service client and logger.
func New(client *user.Client, log *slog.Logger) Handler {
	return Handler{
		usrClient: client,
		log:       log,
	}
}

func (h *Handler) SessionAuthMiddleware() gin.HandlerFunc {
	const op = "SessionAuthMiddleware"

	log := h.log.With(slog.String("op", op))

	return func(c *gin.Context) {
		sessionToken, err := c.Cookie(SessionTokenName)
		if err != nil {
			log.Warn("cookie not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s cookie not found", SessionTokenName)})
			return
		}

		res, err := h.usrClient.Authenticate(c, &userv1.AuthenticateRequest{
			SessionToken: sessionToken,
		})
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

		c.Set("userID", res.GetUserId())

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
