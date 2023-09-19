package main

import (
	"github.com/gin-gonic/gin"
)

func initServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/", handleRoot)
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.File("favicon.ico")
	})

	// API endpoints
	//router.POST("/api/transfer", handleTransfer)
	router.POST("/api/validate", handleValidate)

	return router
}

func handleRoot(ctx *gin.Context) {
	sourceDetails := Credentials{
		Server:   *source_server,
		Account:  *source_account,
		Password: *source_password,
	}

	destinationDetails := Credentials{
		Server:   *destination_server,
		Account:  *destination_account,
		Password: *destination_password,
	}

	data := struct {
		SourceDetails      Credentials
		DestinationDetails Credentials
	}{
		SourceDetails:      sourceDetails,
		DestinationDetails: destinationDetails,
	}
	ctx.HTML(200, "index.html", data)
}
