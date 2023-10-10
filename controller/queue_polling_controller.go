package controller

import (
	"imap-sync/internal"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleQueuePolling(ctx *gin.Context) {
	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	data := internal.GetPollingData(index)

	ctx.HTML(200, "table.html", data)
}
