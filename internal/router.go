package internal

import (
	"github.com/gin-gonic/gin"
)

func InitServer() {
	SetupLogger(log)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/", handleRoot)
	router.GET("/admin", handleAdmin)
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.File("favicon.ico")
	})

	go initQueue()
	// API endpoints
	router.GET("/api/queue", handleQueue)
	router.GET("/api/queuepoll", handleQueuePolling)
	router.GET("/api/pagination", handlePagination)
	router.POST("/api/validate", handleValidate)
	router.POST("/api/search", handleSearch)

	log.Info("Server starting on http://localhost:" + *port)

	if err := router.Run(":" + *port); err != nil {
		log.Fatal(err)
	}
}
