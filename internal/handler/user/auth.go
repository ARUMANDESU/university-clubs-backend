package user

import (
	"fmt"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"time"
)

const SessionTokenName = "session_token"

func (h *Handler) SignUp(c *gin.Context) {
	const op = "UserHandler.SignUp"

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
	const op = "UserHandler.SignIn"

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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	// if https only then secure: true.
	// todo: deal with cookie domain
	/*c.SetCookie(SessionTokenName, res.GetSessionToken(), 3600*24, "/", "localhost:3000", false, true)*/

	t := &http.Cookie{
		Name:     SessionTokenName,
		Value:    res.GetSessionToken(),
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(c.Writer, t)
	c.JSON(http.StatusOK, gin.H{"user": domain.UserObjectToDomain(res.GetUser())})
}

func (h *Handler) Logout(c *gin.Context) {
	const op = "UserHandler.Logout"

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
	c.SetCookie(SessionTokenName, "", -1, "/", "localhost:3000", false, true)

	c.Status(http.StatusOK)
}

func (h *Handler) Activate(c *gin.Context) {
	const op = "UserHandler.Activate"
	log := h.log.With(slog.String("op", op))

	token, ok := c.GetQuery("token")
	if !ok {
		log.Warn("session token was not provided")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "session token was not provided"})
		return
	}

	_, err := h.usrClient.ActivateUser(c, &userv1.ActivateUserRequest{VerificationToken: token})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			log.Warn("invalid arguments", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			log.Warn("not found", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}
