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

const (
	RefreshTokenName = "rt_token"
	AccessTokenName  = "access_token"
)

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

	refreshTokenCookie := &http.Cookie{
		Name:     RefreshTokenName,
		Value:    res.GetRtToken(),
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		Path:     "/",
	}
	accessTokenCookie := &http.Cookie{
		Name:    AccessTokenName,
		Value:   res.GetJwtToken(),
		Expires: time.Now().Add(time.Minute * 15),
		Path:    "/",
	}
	http.SetCookie(c.Writer, refreshTokenCookie)
	http.SetCookie(c.Writer, accessTokenCookie)
	c.JSON(http.StatusOK, gin.H{"user": domain.UserObjectToDomain(res.GetUser())})
}

func (h *Handler) Logout(c *gin.Context) {
	const op = "UserHandler.Logout"
	log := h.log.With(slog.String("op", op))

	cookie, err := c.Cookie(RefreshTokenName)
	if err != nil {
		log.Warn("cookie not found", logger.Err(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s cookie not found", RefreshTokenName)})
		return
	}

	_, err = h.usrClient.Logout(c, &userv1.LogoutRequest{RtToken: cookie})
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
	c.SetCookie(RefreshTokenName, "", -1, "/", "", false, true)
	c.SetCookie(AccessTokenName, "", -1, "/", "", false, false)

	c.Status(http.StatusOK)
}

func (h *Handler) RefreshTokenHandler(c *gin.Context) {
	const op = "UserHandler.RefreshTokenHandler"
	log := h.log.With(slog.String("op", op))

	accessToken, err := c.Cookie(AccessTokenName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s cookie not found", AccessTokenName)})
		return
	}

	refreshToken, err := c.Cookie(RefreshTokenName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s cookie not found", RefreshTokenName)})
		return
	}

	res, err := h.usrClient.RefreshToken(c, &userv1.RefreshTokenRequest{
		RtToken:  refreshToken,
		JwtToken: accessToken,
	})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": status.Convert(err).Message()})
		case status.Code(err) == codes.NotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": status.Convert(err).Message()})
		default:
			log.Error("internal", logger.Err(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	refreshTokenCookie := &http.Cookie{
		Name:     RefreshTokenName,
		Value:    res.GetRtToken(),
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
		Path:     "/",
	}
	accessTokenCookie := &http.Cookie{
		Name:    AccessTokenName,
		Value:   res.GetJwtToken(),
		Expires: time.Now().Add(time.Minute * 15),
		Path:    "/",
	}
	http.SetCookie(c.Writer, refreshTokenCookie)
	http.SetCookie(c.Writer, accessTokenCookie)
	c.JSON(http.StatusOK, gin.H{"user": domain.UserObjectToDomain(res.GetUser())})
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

func (h *Handler) MicrosoftOIDCLogin(c *gin.Context) {
	const op = "UserHandler.MicrosoftOIDCLogin"
	log := h.log.With(slog.String("op", op))

	authURL, err := h.confClient.AuthCodeURL(c, h.MicrosoftOIDC.ClientID, "http://localhost:5000/auth/microsoft/callback", []string{"openid", "profile", "email"})
	if err != nil {
		log.Error("failed to create auth URL", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create auth URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": authURL})
}

func (h *Handler) MicrosoftOIDCCallback(c *gin.Context) {
	const op = "UserHandler.MicrosoftOIDCCallback"
	log := h.log.With(slog.String("op", op))

	authCode := c.Query("code")
	if authCode == "" {
		log.Error("missing authorization code")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing authorization code"})
		return
	}

	result, err := h.confClient.AcquireTokenByAuthCode(c, authCode, "http://localhost:5000/auth/microsoft/callback", []string{"openid", "profile", "email"})
	if err != nil {
		log.Error("failed to acquire token", logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to acquire token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": result})
}
