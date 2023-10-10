package controller

import (
	"imap-sync/internal"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleQueue(ctx *gin.Context) {
	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	data := internal.GetQueueData(index)

	ctx.HTML(200, "tbody.html", data)
}
