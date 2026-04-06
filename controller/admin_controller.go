package controller

import (
	"imap-sync/internal"
	"net/http"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
)

type TableData struct {
	TableColumnNum int
}

type AdminData struct {
	Settings *internal.Settings
	Data     map[string]string
}

func HandleAdmin(ctx *gin.Context) {
	store := ginsession.FromContext(ctx)
	_, ok := store.Get("user")
	if !ok {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	settings, err := internal.GetSettings()
	if err != nil {
		settings = &internal.Settings{}
	}

	adminData := AdminData{
		Settings: settings,
		Data:     internal.Data["admin"],
	}

	ctx.HTML(200, "admin.html", adminData)
}
