package main

import "github.com/gin-gonic/gin"

func initServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", handleRoot)
	return r
}

func handleRoot(c *gin.Context) {
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
	c.HTML(200, "index.html", data)
}
