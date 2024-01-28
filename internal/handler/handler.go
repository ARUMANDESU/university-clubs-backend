package handler

import (
	usergrpc "github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/user"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Handler struct {
	UsrHandler user.Handler
}

func New(log *slog.Logger, usrClient *usergrpc.Client) *Handler {

	return &Handler{
		UsrHandler: user.New(usrClient, log),
	}
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

	userPath := router.Group("/user")
	{
		userPath.Use(h.UsrHandler.SessionAuthMiddleware())

		//todo: remove this later
		userPath.POST("/lol", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"lol": "kek"})
		})
	}

	//TODO: implement other  endpoints

	return router
}
