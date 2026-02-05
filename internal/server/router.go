// Package server 提供 HTTP 服务器路由配置
package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/yeegeek/uyou-go-api-starter/internal/auth"
	"github.com/yeegeek/uyou-go-api-starter/internal/config"
	"github.com/yeegeek/uyou-go-api-starter/internal/errors"
	"github.com/yeegeek/uyou-go-api-starter/internal/health"
	"github.com/yeegeek/uyou-go-api-starter/internal/middleware"
	"github.com/yeegeek/uyou-go-api-starter/internal/friend"
	"github.com/yeegeek/uyou-go-api-starter/internal/message"
	"github.com/yeegeek/uyou-go-api-starter/internal/notification"
	"github.com/yeegeek/uyou-go-api-starter/internal/payment"
	"github.com/yeegeek/uyou-go-api-starter/internal/post"
	"github.com/yeegeek/uyou-go-api-starter/internal/user"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(userHandler *user.Handler, friendHandler *friend.Handler, messageHandler *message.Handler, postHandler *post.Handler, notificationHandler *notification.Handler, paymentHandler *payment.Handler, authService auth.Service, cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := gin.New()

	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	skipPaths := config.GetSkipPaths(cfg.App.Environment)
	loggerConfig := middleware.NewLoggerConfig(
		cfg.Logging.GetLogLevel(),
		skipPaths,
	)
	router.Use(middleware.Logger(loggerConfig))
	router.Use(errors.ErrorHandler())
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	router.Use(cors.New(corsConfig))

	var checkers []health.Checker
	if cfg.Health.DatabaseCheckEnabled {
		dbChecker := health.NewDatabaseChecker(db)
		checkers = append(checkers, dbChecker)
	}
	healthService := health.NewService(checkers, cfg.App.Version, cfg.App.Environment)
	healthHandler := health.NewHandler(healthService)

	router.GET("/health", healthHandler.Health)
	router.GET("/health/live", healthHandler.Live)
	router.GET("/health/ready", healthHandler.Ready)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rlCfg := cfg.Ratelimit
	if rlCfg.Enabled {
		router.Use(
			middleware.NewRateLimitMiddleware(
				rlCfg.Window,
				rlCfg.Requests,
				func(c *gin.Context) string {
					ip := c.ClientIP()
					if ip == "" {
						ip = c.GetHeader("X-Forwarded-For")
						if ip == "" {
							ip = c.GetHeader("X-Real-IP")
						}
						if ip == "" {
							ip = "unknown"
						}
					}
					return ip
				},
				nil,
			),
		)
	}

	v1 := router.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", userHandler.Register)
			authGroup.POST("/login", userHandler.Login)
			authGroup.POST("/refresh", userHandler.RefreshToken)
			authGroup.POST("/logout", auth.AuthMiddleware(authService), userHandler.Logout)
			authGroup.GET("/me", auth.AuthMiddleware(authService), userHandler.GetMe)
		}

		// User endpoints - authenticated users can access their own resources
		usersGroup := v1.Group("/users")
		usersGroup.Use(auth.AuthMiddleware(authService))
		{
			usersGroup.GET("/:id", userHandler.GetUser)
			usersGroup.PUT("/:id", userHandler.UpdateUser)
			usersGroup.DELETE("/:id", userHandler.DeleteUser)
		}

		// Admin endpoints - admin role required, following REST best practices
		adminGroup := v1.Group("/admin")
		adminGroup.Use(auth.AuthMiddleware(authService), middleware.RequireAdmin())
		{
			// User management endpoints
			adminGroup.GET("/users", userHandler.ListUsers)
			adminGroup.GET("/users/:id", userHandler.GetUser)
			adminGroup.PUT("/users/:id", userHandler.UpdateUser)
			adminGroup.DELETE("/users/:id", userHandler.DeleteUser)
		}

		// Friend endpoints
		friendsGroup := v1.Group("/friends")
		friendsGroup.Use(auth.AuthMiddleware(authService))
		{
			friendsGroup.GET("", friendHandler.GetFriendsList)
			friendsGroup.POST("/request", friendHandler.SendFriendRequest)
			friendsGroup.GET("/requests", friendHandler.GetFriendRequests)
			friendsGroup.POST("/requests/:id/accept", friendHandler.AcceptFriendRequest)
			friendsGroup.POST("/requests/:id/reject", friendHandler.RejectFriendRequest)
			friendsGroup.DELETE("/:id", friendHandler.DeleteFriend)
			friendsGroup.PUT("/:id/remark", friendHandler.UpdateFriendRemark)
			friendsGroup.PUT("/:id/group", friendHandler.UpdateFriendGroup)
			friendsGroup.POST("/:id/block", friendHandler.BlockUser)
			friendsGroup.DELETE("/:id/unblock", friendHandler.UnblockUser)
			friendsGroup.GET("/blocked", friendHandler.GetBlockedUsers)
		}

		// Message endpoints
		messagesGroup := v1.Group("/messages")
		messagesGroup.Use(auth.AuthMiddleware(authService))
		{
			messagesGroup.GET("/conversations", messageHandler.GetConversations)
			messagesGroup.GET("/conversations/:conversation_id/messages", messageHandler.GetMessages)
			messagesGroup.POST("/send", messageHandler.SendMessage)
			messagesGroup.POST("/:message_id/read", messageHandler.MarkAsRead)
			messagesGroup.POST("/conversations/:conversation_id/read", messageHandler.MarkConversationAsRead)
			messagesGroup.POST("/:message_id/recall", messageHandler.RecallMessage)
			messagesGroup.POST("/drafts", messageHandler.SaveDraft)
			messagesGroup.GET("/drafts", messageHandler.GetDrafts)
			messagesGroup.DELETE("/drafts/:draft_id", messageHandler.DeleteDraft)
			messagesGroup.DELETE("/conversations/:conversation_id", messageHandler.DeleteConversation)
		}

		// Post endpoints
		postsGroup := v1.Group("/posts")
		postsGroup.Use(auth.AuthMiddleware(authService))
		{
			postsGroup.POST("", postHandler.CreatePost)
			postsGroup.GET("", postHandler.GetPosts)
			postsGroup.GET("/:post_id", postHandler.GetPost)
			postsGroup.GET("/user/:user_id", postHandler.GetUserPosts)
			postsGroup.PUT("/:post_id", postHandler.UpdatePost)
			postsGroup.DELETE("/:post_id", postHandler.DeletePost)
			postsGroup.POST("/:post_id/like", postHandler.LikePost)
			postsGroup.POST("/:post_id/unlike", postHandler.UnlikePost)
			postsGroup.POST("/comments", postHandler.CreateComment)
			postsGroup.GET("/:post_id/comments", postHandler.GetComments)
			postsGroup.DELETE("/comments/:comment_id", postHandler.DeleteComment)
		}

		// Notification endpoints
		notificationsGroup := v1.Group("/notifications")
		notificationsGroup.Use(auth.AuthMiddleware(authService))
		{
			notificationsGroup.GET("", notificationHandler.GetNotifications)
			notificationsGroup.GET("/unread-count", notificationHandler.GetUnreadCount)
			notificationsGroup.POST("/:notification_id/read", notificationHandler.MarkAsRead)
			notificationsGroup.POST("/read-all", notificationHandler.MarkAllAsRead)
			notificationsGroup.DELETE("/:notification_id", notificationHandler.DeleteNotification)
		}

		// Payment endpoints
		paymentGroup := v1.Group("/payment")
		paymentGroup.Use(auth.AuthMiddleware(authService))
		{
			paymentGroup.POST("/recharge", paymentHandler.Recharge)
			paymentGroup.POST("/vip-upgrade", paymentHandler.VIPUpgrade)
			paymentGroup.POST("/purchase", paymentHandler.Purchase)
			paymentGroup.GET("/transactions", paymentHandler.GetTransactions)
			paymentGroup.GET("/balance", paymentHandler.GetBalance)
		}
	}

	return router
}
