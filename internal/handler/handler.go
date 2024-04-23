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
	socketio "github.com/googollee/go-socket.io"
	"log"
	"log/slog"
)

type Handler struct {
	UsrHandler          user.Handler
	ClubHandler         club.Handler
	notificationService user.NotificationService
}

func New(log *slog.Logger, microsoftOIDC config.MicrosoftOIDC, usrClient *usergrpc.Client, clubClient *clubgrpc.Client, confClient confidential.Client, notificationService user.NotificationService) *Handler {

	return &Handler{
		UsrHandler:          user.New(usrClient, log, confClient, microsoftOIDC, notificationService),
		ClubHandler:         club.New(clubClient, log),
		notificationService: notificationService,
	}
}

func (h *Handler) InitRoutes() (*gin.Engine, *socketio.Server) {
	router := gin.New()

	// Socket IO
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Printf("connected: %s, URL: %s", s.ID(), s.URL())
		h.notificationService.AddNewConnection(1)
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/", "notification", func(s socketio.Conn) {
		log.Println("notification:")
		go func() {
			notifChan, err := h.notificationService.GetConnection(1)
			if err != nil {
				return
			}

			for {
				if notification, ok := <-notifChan; ok {
					s.Emit("notification", notification)
				}
			}
		}()
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Printf("disconnected: %s, reason: %s", s.ID(), msg)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

	// Cors
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"http://localhost:3000"}
	corsCfg.AllowHeaders = []string{
		"Accept", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token",
		"Token", "session", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "X-Requested-With"}
	corsCfg.AllowCredentials = true

	// Middlewares
	router.Use(gin.Logger(), cors.New(corsCfg), gin.Recovery())

	// ALL Routes
	// Socket IO routes
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

	// Rest API
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
		userPath.GET("/:id/clubs", h.ClubHandler.GetUserClubsHandler)
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
		clubPath.GET("/:id/members/:member_id", h.ClubHandler.GetClubMember)
		clubPath.GET("/:id/members/:member_id/roles", h.ClubHandler.GetMemberRoles)
		clubPath.GET("/:id", h.ClubHandler.GetClubHandler)

		clubPathAuth := clubPath.Group("")
		{
			clubPathAuth.Use(h.UsrHandler.SessionAuthMiddleware())
			clubPathAuth.POST("/:id", h.UsrHandler.RoleAuthMiddleware([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN}), h.ClubHandler.NewClubHandler)
			clubPathAuth.GET("/pending", h.UsrHandler.RoleAuthMiddleware([]userv1.Role{userv1.Role_DSVR, userv1.Role_ADMIN}), h.ClubHandler.ListNewClubRequestsHandler)

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

		}
	}

	//TODO: implement other  endpoints

	return router, server
}
