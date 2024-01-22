package user

import (
	"context"
	uniclubs_user_service_v1_userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Handler struct {
	usrClient *user.Client
	log       *slog.Logger
}

func New(client *user.Client, log *slog.Logger) Handler {
	return Handler{
		usrClient: client,
		log:       log,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	//todo: write properly: error handling & logging

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
		h.log.Error("some err", err)
		// todo : handle error
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	response, err := h.usrClient.Api.Register(context.TODO(), &uniclubs_user_service_v1_userv1.RegisterRequest{
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
		h.log.Error("some err", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"userID": response.GetUserId()})
}

func (h *Handler) SignIn(c *gin.Context) {
	//TODO: Implement
	panic("Implement")
}

func (h *Handler) Logout(c *gin.Context) {
	//TODO: Implement
	panic("Implement")
}
