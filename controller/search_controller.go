package controller

import (
	"imap-sync/internal"

	"github.com/gin-gonic/gin"
)

func HandleSearch(ctx *gin.Context) {
	searchQuery := ctx.PostForm("search-input")

	if searchQuery == "" {
		HandleQueue(ctx)
		return
	}

	data := internal.GetSearchData(searchQuery)

	ctx.HTML(200, "tbody.html", data)
}
