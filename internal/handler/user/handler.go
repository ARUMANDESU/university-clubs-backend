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

const SessionTokenName = "session_token"

func (h *Handler) SignUp(c *gin.Context) {
	const op = "sign-up"

	log := h.log.With(slog.String("op", op))

	//request struct
	usr := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Barcode   string `json:"barcode"`
		Major     string `json:"major"`
		GroupName string `json:"group_name"`
		Year      int    `json:"year"`
	}{}
	err := c.ShouldBindJSON(&usr)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.usrClient.Register(c, &userv1.RegisterRequest{
		Email:     usr.Email,
		Password:  usr.Password,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Barcode:   usr.Barcode,
		Major:     usr.Major,
		GroupName: usr.GroupName,
		Year:      int32(usr.Year),
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.AlreadyExists:
			log.Warn("user already exists", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{"userID": res.GetUserId()})
}

func (h *Handler) SignIn(c *gin.Context) {
	const op = "sign-in"

	log := h.log.With(slog.String("op", op))

	usr := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := c.ShouldBindJSON(&usr)
	if err != nil {
		log.Error("decoding err", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usrClient.Login(c, &userv1.LoginRequest{
		Email:    usr.Email,
		Password: usr.Password,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	// if https only then secure: true.
	// todo: deal with cookie domain
	c.SetCookie(SessionTokenName, res.GetSessionToken(), -1, "/", "localhost", false, true)

	c.Status(http.StatusOK)
}

func (h *Handler) Logout(c *gin.Context) {
	const op = "logout"

	log := h.log.With(slog.String("op", op))

	cookie, err := c.Cookie(SessionTokenName)
	if err != nil {
		log.Warn("cookie not found", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s cookie not found", SessionTokenName)})
		return
	}

	_, err = h.usrClient.Logout(c, &userv1.LogoutRequest{SessionToken: cookie})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	// if https only then secure: true.
	// todo: deal with cookie domain
	c.SetCookie(SessionTokenName, "", -1, "/", "localhost", false, true)

	c.Status(http.StatusOK)
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

		_, err = h.usrClient.Authenticate(c, &userv1.AuthenticateRequest{
			SessionToken: sessionToken,
		})
		if err != nil {
			switch {
			case status.Code(err) == codes.InvalidArgument:
				log.Warn("invalid arguments", logger.Err(err))
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
			case status.Code(err) == codes.NotFound:
				log.Warn("session not found", logger.Err(err))
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
			default:
				log.Error("internal", logger.Err(err))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		c.Next()
	}
}
