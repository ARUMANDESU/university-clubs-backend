package handler

import (
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	clubgrpc "github.com/ARUMANDESU/university-clubs-backend/internal/clients/club"
	usergrpc "github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/config"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/club"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/user"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	UsrHandler  user.Handler
	ClubHandler club.Handler
}

func New(log *slog.Logger, microsoftOIDC config.MicrosoftOIDC, usrClient *usergrpc.Client, clubClient *clubgrpc.Client, confClient confidential.Client) *Handler {

	return &Handler{
		UsrHandler:  user.New(usrClient, log, confClient, microsoftOIDC),
		ClubHandler: club.New(clubClient, log),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"http://localhost:3000"}
	corsCfg.AllowCredentials = true

	router.Use(gin.Logger(), cors.New(corsCfg), gin.Recovery())

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.UsrHandler.SignUp)
		auth.POST("/sign-in", h.UsrHandler.SignIn)
		auth.POST("/logout", h.UsrHandler.Logout)
		auth.POST("/activate", h.UsrHandler.Activate)
		auth.POST("/microsoft/login", h.UsrHandler.MicrosoftOIDCLogin)
		auth.GET("/microsoft/callback", h.UsrHandler.MicrosoftOIDCCallback)
	}

	userPath := router.Group("/user")
	{
		userPath.GET("/:id", h.UsrHandler.GetUser)
		userPath.GET("/search", h.UsrHandler.SearchUsers)

		userPathAuth := userPath.Group("")
		{
			userPathAuth.Use(h.UsrHandler.SessionAuthMiddleware())

			userPathAuth.PATCH("/:id", h.UsrHandler.UpdateUser)
			userPathAuth.PATCH("/:id/avatar", h.UsrHandler.UpdateAvatar)
			userPathAuth.PATCH("/:id/roles", h.UsrHandler.RoleAuthMiddleware([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN}), h.UsrHandler.ChangeUserRole)

			userPathAuth.DELETE("/:id", h.UsrHandler.DeleteUser)
		}

	}

	clubPath := router.Group("/clubs")
	{
		clubPath.GET("/", h.ClubHandler.ListClubsHandler)
		clubPath.GET("/:id/members", h.ClubHandler.ListClubMembersHandler)
		clubPath.GET("/:id", h.ClubHandler.GetClubHandler)

		clubPathAuth := clubPath.Group("")
		{
			clubPathAuth.Use(h.UsrHandler.SessionAuthMiddleware())
			clubPathAuth.POST("/:id", h.UsrHandler.RoleAuthMiddleware([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN}), h.ClubHandler.NewClubHandler)
			clubPathAuth.GET("/pending", h.UsrHandler.RoleAuthMiddleware([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN}), h.ClubHandler.ListNewClubRequestsHandler)

			clubPathAuth.POST("/:id/members", h.ClubHandler.HandleJoinRequestHandler)
			clubPathAuth.GET("/:id/join", h.ClubHandler.ListJoinRequestsHandler)
			clubPathAuth.POST("/:id/join", h.ClubHandler.JoinRequestHandler)
			clubPathAuth.POST("", h.ClubHandler.CreateClubHandler)

			clubPathAuth.PATCH("/:id", h.ClubHandler.UpdateClubHandler)
			clubPathAuth.PATCH("/:id/logo", h.ClubHandler.UpdateLogoHandler)
			clubPathAuth.PATCH("/:id/banner", h.ClubHandler.UpdateBannerHandler)

			clubPathAuth.POST("/:id/roles", h.ClubHandler.CreateRoleHandler)
			clubPathAuth.PATCH("/:id/roles", h.ClubHandler.UpdateRolesPositionHandler)
			clubPathAuth.DELETE("/:id/roles/:role_id", h.ClubHandler.DeleteRoleHandler)
			clubPathAuth.PATCH("/:id/roles/:role_id", h.ClubHandler.UpdateRoleHandler)
			clubPathAuth.PATCH("/:id/roles/:role_id/members", h.ClubHandler.AddRoleMembersHandler)

		}
	}

	//TODO: implement other  endpoints

	return router
}
