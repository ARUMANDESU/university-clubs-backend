package handler

import (
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	clubgrpc "github.com/ARUMANDESU/university-clubs-backend/internal/clients/club"
	eventgrpc "github.com/ARUMANDESU/university-clubs-backend/internal/clients/event"
	usergrpc "github.com/ARUMANDESU/university-clubs-backend/internal/clients/user"
	"github.com/ARUMANDESU/university-clubs-backend/internal/config"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/club"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/event"
	"github.com/ARUMANDESU/university-clubs-backend/internal/handler/user"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	cfg          *config.Config
	UsrHandler   user.Handler
	ClubHandler  club.Handler
	EventHandler event.Handler
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	usrClient *usergrpc.Client,
	clubClient *clubgrpc.Client,
	confClient confidential.Client,
	imageStorage user.ImageStorage,
	fileStorage event.FileStorage,
	eventClient *eventgrpc.Client,
) *Handler {

	return &Handler{
		cfg:          cfg,
		UsrHandler:   user.New(cfg, log, usrClient, confClient, imageStorage),
		ClubHandler:  club.New(log, clubClient, imageStorage),
		EventHandler: event.New(log, eventClient, clubClient, usrClient, imageStorage, fileStorage),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Cors
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOriginFunc = func(origin string) bool {
		for _, allowedOrigin := range h.cfg.AllowedOrigins {
			if origin == allowedOrigin {
				return true
			}
		}
		return false
	}
	corsCfg.AllowHeaders = []string{
		"Accept", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token",
		"Token", "session", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "X-Requested-With"}
	corsCfg.AllowCredentials = true

	// Middlewares
	router.Use(gin.Logger(), cors.New(corsCfg), gin.Recovery())

	// ALL Routes

	// Rest API
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.UsrHandler.SignUp)
		auth.POST("/sign-in", h.UsrHandler.SignIn)
		auth.POST("/logout", h.UsrHandler.Logout)
		auth.POST("/activate", h.UsrHandler.Activate)
		auth.POST("/refresh", h.UsrHandler.RefreshTokenHandler)
		auth.POST("/microsoft/login", h.UsrHandler.MicrosoftOIDCLogin)
		auth.GET("/microsoft/callback", h.UsrHandler.MicrosoftOIDCCallback)
	}

	userPath := router.Group("/users")
	{
		userPath.GET("/:id", h.UsrHandler.GetUser)
		userPath.GET("/:id/clubs", h.ClubHandler.GetUserClubsHandler)
		userPath.GET("/search", h.UsrHandler.SearchUsers)

		userPathAuth := userPath.Group("")
		{
			userPathAuth.Use(h.UsrHandler.AuthMiddleware())

			userPathAuth.PATCH("/:id", h.UsrHandler.UpdateUser)
			userPathAuth.PATCH("/:id/avatar", h.UsrHandler.UpdateAvatar)
			userPathAuth.PATCH("/:id/roles", h.UsrHandler.RoleAuthMiddleware([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN}), h.UsrHandler.ChangeUserRole)

			userPathAuth.DELETE("/:id", h.UsrHandler.DeleteUser)

			userPathAuth.GET("/events/invites", h.EventHandler.GetUserInvites)
			userPathAuth.POST("/invites/:invite_id/handle", h.EventHandler.HandleUserInvite)
		}

	}

	clubPath := router.Group("/clubs")
	{
		clubPath.GET("/", h.ClubHandler.ListClubsHandler)
		clubPath.GET("/:id/members", h.ClubHandler.ListClubMembersHandler)
		clubPath.GET("/:id/members/:member_id", h.ClubHandler.GetClubMember)
		clubPath.GET("/:id/members/:member_id/roles", h.ClubHandler.GetMemberRoles)
		clubPath.GET("/:id", h.ClubHandler.GetClubHandler)

		clubPath.GET("/:id/events", h.EventHandler.ListClubEventsHandler)

		clubPathAuth := clubPath.Group("")
		{
			clubPathAuth.Use(h.UsrHandler.AuthMiddleware())
			clubPathAuth.POST("/:id", h.UsrHandler.RoleAuthMiddleware([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN}), h.ClubHandler.NewClubHandler)
			clubPathAuth.GET("/pending", h.UsrHandler.RoleAuthMiddleware([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN}), h.ClubHandler.ListNewClubRequestsHandler)
			clubPathAuth.DELETE("/:id", h.UsrHandler.HasRoles([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN}), h.ClubHandler.DeleteClubHandler)

			clubPathAuth.POST("/:id/members", h.ClubHandler.HandleJoinRequestHandler)
			clubPathAuth.DELETE("/:id/members", h.ClubHandler.LeaveClubHandler)
			clubPathAuth.DELETE("/:id/members/:member_id", h.ClubHandler.KickMemberHandler)
			clubPathAuth.GET("/:id/join", h.ClubHandler.ListJoinRequestsHandler)
			clubPathAuth.GET("/:id/join/status", h.ClubHandler.GetUserJoinStatus)
			clubPathAuth.POST("/:id/join", h.ClubHandler.JoinRequestHandler)
			clubPathAuth.POST("", h.ClubHandler.CreateClubHandler)

			clubPathAuth.PATCH("/:id", h.ClubHandler.UpdateClubHandler)
			clubPathAuth.PATCH("/:id/logo", h.ClubHandler.UpdateLogoHandler)
			clubPathAuth.PATCH("/:id/banner", h.ClubHandler.UpdateBannerHandler)

			clubPathAuth.POST("/:id/roles", h.ClubHandler.CreateRoleHandler)
			clubPathAuth.PATCH("/:id/roles", h.ClubHandler.UpdateRolesPositionHandler)
			clubPathAuth.DELETE("/:id/roles/:role_id", h.ClubHandler.DeleteRoleHandler)
			clubPathAuth.PATCH("/:id/roles/:role_id", h.ClubHandler.UpdateRoleHandler)
			clubPathAuth.POST("/:id/roles/:role_id/members", h.ClubHandler.AddRoleMembersHandler)
			clubPathAuth.DELETE("/:id/roles/:role_id/members", h.ClubHandler.RemoveRoleMembersHandler)

			clubPathAuth.POST("/:id/members/:member_id/ban", h.ClubHandler.BanMemberHandler)
			clubPathAuth.DELETE("/:id/members/:member_id/ban", h.ClubHandler.UnbanMemberHandler)
			clubPathAuth.GET("/:id/bans", h.ClubHandler.ListBannedMembersHandler)
			clubPathAuth.PATCH("/:id/ownership/:member_id", h.ClubHandler.TransferOwnershipHandler)

			clubPathAuth.POST("/:id/events", h.EventHandler.CreateEventHandler)
			clubPathAuth.GET("/:id/events/manage", h.EventHandler.ListRestrictedClubEvents)

			clubPathAuth.POST("/:id/invites/:invite_id/handle", h.EventHandler.HandleCollaboratorRequestHandler)
			clubPathAuth.GET("/:id/invites", h.EventHandler.ListCollaboratorRequestsHandler)
		}
	}

	eventPath := router.Group("/events")
	{
		eventPath.GET("/:id", h.UsrHandler.GetUserIDMiddleware(), h.EventHandler.GetEventHandler)
		eventPath.GET("", h.EventHandler.ListPublishedEventsHandler)

		eventPathAuth := eventPath.Group("")
		{
			eventPathAuth.Use(h.UsrHandler.AuthMiddleware())

			eventPathAuth.PATCH("/:id", h.EventHandler.UpdateEventHandler)
			eventPathAuth.DELETE("/:id", h.UsrHandler.HasRoles([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN, userv1.Role_MODER}), h.EventHandler.DeleteEventHandler)

			eventPathAuth.PATCH("/:id/review", h.EventHandler.SendEventForReviewHandler)
			eventPathAuth.DELETE("/:id/review", h.EventHandler.CancelEventReviewHandler)

			eventPathAuth.POST("/:id/upload/files", h.EventHandler.UploadFilesHandler)
			eventPathAuth.POST("/:id/upload/images", h.EventHandler.UploadImagesHandler)

			eventPathAuth.DELETE("/:id/upload/files", h.EventHandler.DeleteFileHandler)

			eventPathAuth.POST("/:id/collaborators", h.EventHandler.AddCollaboratorHandler)
			eventPathAuth.DELETE("/:id/collaborators/:collaborator_id", h.EventHandler.RemoveCollaboratorHandler)
			eventPathAuth.DELETE("/invites/:invite_id/collaborators", h.EventHandler.CancelCollaboratorRequestHandler)
			eventPathAuth.GET("/:id/invites/collaborators", h.EventHandler.ListCollaboratorInvitesHandler)

			eventPathAuth.POST("/:id/organizers", h.EventHandler.AddOrganizerHandler)
			eventPathAuth.DELETE("/:id/organizers/:organizer_id", h.EventHandler.RemoveOrganizerHandler)
			eventPathAuth.DELETE("/invites/:invite_id/organizers", h.EventHandler.CancelOrganizerRequestHandler)
			eventPathAuth.GET("/:id/invites/organizers", h.EventHandler.ListOrganizerInvitesHandler)

			eventPathAuth.PATCH("/:id/publish", h.EventHandler.PublishEventHandler)
			eventPathAuth.PATCH("/:id/unpublish", h.EventHandler.UnpublishEventHandler)

			eventPathAuth.POST("/:id/participants", h.EventHandler.AddParticipantHandler)
			eventPathAuth.DELETE("/:id/participants", h.EventHandler.LeaveEventHandler)
			eventPathAuth.DELETE("/:id/participants/:participant_id", h.EventHandler.RemoveParticipantHandler)
			eventPathAuth.POST("/:id/participants/:participant_id/ban", h.EventHandler.BanParticipantHandler)

			eventPathAuthAdmin := eventPathAuth.Group("")
			{
				eventPathAuthAdmin.Use(h.UsrHandler.RoleAuthMiddleware([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN}))

				eventPathAuth.GET("/admin", h.EventHandler.ListEventsHandler)
				eventPathAuth.PATCH("/:id/approve", h.EventHandler.ApproveEventHandler)
				eventPathAuth.PATCH("/:id/reject", h.EventHandler.RejectEventHandler)
			}
		}
	}

	return router
}
