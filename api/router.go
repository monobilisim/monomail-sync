package api

import (
	"imap-sync/config"
	"imap-sync/controller"
	"imap-sync/internal"
	"imap-sync/logger"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
)

var log = logger.Log
var port string

func InitServer() {
	port = config.Conf.Port
	logger.SetupLogger()
	err := internal.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	err = internal.InitSettingsTable()
	if err != nil {
		log.Error(err)
	}

	internal.InitLocalizer()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(ginsession.New())

	router.LoadHTMLGlob("templates/*")

	router.Static("/static", "./static/")

	router.GET("/", controller.HandleRoot)
	router.GET("/admin", controller.HandleAdmin)
	router.GET("/login", controller.HandleLogin)
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.File("favicon.ico")
	})
	router.GET("/health", func(ctx *gin.Context) {
		dbStatus := "healthy"
		if err := internal.CheckDB(); err != nil {
			dbStatus = "unhealthy"
		}

		ctx.JSON(200, gin.H{
			"status":    "healthy",
			"version":   "1.0.0",
			"db_status": dbStatus,
			"uptime":    internal.GetUptime(),
		})
	})
	go internal.InitQueue()

	authenticatedAPI := router.Group("/api")
	authenticatedAPI.Use(requireSession())
	// API endpoints
	authenticatedAPI.GET("/queue", controller.HandleQueue)
	authenticatedAPI.GET("/queuepoll", controller.HandleQueuePolling)
	authenticatedAPI.GET("/pagination", controller.HandlePagination)
	authenticatedAPI.GET("/details", controller.HandleGetLog)
	authenticatedAPI.GET("/sync", controller.HandleSync)
	authenticatedAPI.GET("/settings", controller.HandleGetSettings)
	authenticatedAPI.PUT("/settings", controller.HandleUpdateSettings)
	authenticatedAPI.POST("/bulk", controller.HandleBulkMigration)
	authenticatedAPI.GET("/bulk/status", controller.HandleBulkMigrationStatus)
	authenticatedAPI.GET("/stats", controller.HandleGetStats)
	authenticatedAPI.GET("/system", controller.HandleGetSystemInfo)
	authenticatedAPI.GET("/audit", controller.HandleGetAuditLog)
	authenticatedAPI.GET("/sessions", controller.HandleGetSessions)
	authenticatedAPI.POST("/sessions/:id/terminate", controller.HandleTerminateSession)
	authenticatedAPI.POST("/sessions/terminate-all", controller.HandleTerminateAllSessions)
	authenticatedAPI.POST("/validate", controller.HandleValidate)
	authenticatedAPI.POST("/search", controller.HandleSearch)

	router.POST("/auth/login", controller.Login)

	log.Info("Server starting on http://localhost:" + port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
