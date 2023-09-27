package internal

import "github.com/gin-gonic/gin"

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
