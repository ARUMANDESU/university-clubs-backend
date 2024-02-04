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
		auth.POST("/activate", h.UsrHandler.Activate)
	}

	userPath := router.Group("/user")
	{
		userPath.GET("/:id", h.UsrHandler.GetUser)
		userPath.GET("/search", h.UsrHandler.SearchUsers)

		userPath.PATCH("/:id", h.UsrHandler.SessionAuthMiddleware(), h.UsrHandler.UpdateUser)
		userPath.DELETE("/:id", h.UsrHandler.SessionAuthMiddleware(), h.UsrHandler.DeleteUser)
	}

	//TODO: implement other  endpoints

	return router
}
