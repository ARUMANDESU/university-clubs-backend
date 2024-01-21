package handler

import (
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/user"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	UsrHandler user.Handler
}

func New() *Handler {
	return &Handler{}
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

	return router
}
