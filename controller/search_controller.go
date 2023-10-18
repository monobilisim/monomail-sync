package controller

import (
	"imap-sync/internal"

	"github.com/gin-gonic/gin"
)

func HandleSearch(ctx *gin.Context) {
	searchQuery := ctx.PostForm("search-input")
	exact := ctx.Query("exact")
	var data internal.PageData

	if exact == "true" {

		sourceCreds := internal.Credentials{
			Server:  ctx.Query("source_server"),
			Account: ctx.Query("source_account"),
		}
		destCreds := internal.Credentials{
			Server:  ctx.Query("destination_server"),
			Account: ctx.Query("destination_account"),
		}

		data = internal.GetSearchData("", true, sourceCreds, destCreds)
	} else {

		data = internal.GetSearchData(searchQuery, false, internal.Credentials{}, internal.Credentials{})
	}

	ctx.HTML(200, "tbody.html", data)
}
