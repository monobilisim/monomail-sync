package internal

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
)

func InitServer() {
	SetupLogger(log)

	err := InitDb()
	if err != nil {
		log.Error(err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(ginsession.New())

	router.LoadHTMLGlob("templates/*")

	router.Static("/static", "./static/")

	router.GET("/", handleRoot)
	router.GET("/admin", handleAdmin)
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.File("favicon.ico")
	})
	router.GET("/login", handleLogin)

	go initQueue()
	// API endpoints
	router.GET("/api/queue", handleQueue)
	router.GET("/api/queuepoll", handleQueuePolling)
	router.GET("/api/pagination", handlePagination)
	router.POST("/api/validate", handleValidate)
	router.POST("/api/search", handleSearch)
	router.POST("/auth/login", login)

	log.Info("Server starting on http://localhost:" + *port)

	if err := router.Run(":" + *port); err != nil {
		log.Fatal(err)
	}
}
