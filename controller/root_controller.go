package controller

import (
	"imap-sync/internal"
	"imap-sync/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
)

var log = logger.Log

func HandleRoot(ctx *gin.Context) {
	store := ginsession.FromContext(ctx)
	if _, ok := store.Get("user"); !ok {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	settings, err := internal.GetSettings()
	if err != nil {
		log.Error(err)
		settings = &internal.Settings{}
	}

	sourceDetails := internal.Credentials{
		Server:  settings.SourceServer,
		Account: settings.SourceAccountPrefix,
	}

	destinationDetails := internal.Credentials{
		Server:  settings.DestinationServer,
		Account: settings.DestinationAccountPrefix,
	}

	data := struct {
		SourceDetails      internal.Credentials
		DestinationDetails internal.Credentials
		Text               map[string]string
		Table              map[string]string
		Validation         map[string]string
	}{
		SourceDetails:      sourceDetails,
		DestinationDetails: destinationDetails,
		Text:               internal.Data["index"],
		Table:              internal.Data["table"],
		Validation:         internal.Data["validation"],
	}
	ctx.HTML(200, "index.html", data)
}
