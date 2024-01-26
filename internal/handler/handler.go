package handler

import (
	usergrpc "github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/user"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	UsrHandler user.Handler
}

func New(log *slog.Logger, usrClient *usergrpc.Client) *Handler {

	return &Handler{UsrHandler: user.New(usrClient, log)}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.UsrHandler.SignUp)
		auth.POST("/sign-in", h.UsrHandler.SignIn)
		auth.POST("/logout", h.UsrHandler.Logout)
	}

	//TODO: implement other api endpoints

	return router
}
