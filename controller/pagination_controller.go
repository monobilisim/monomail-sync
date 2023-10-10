package controller

import (
	"imap-sync/internal"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandlePagination(ctx *gin.Context) {
	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))
	pages := internal.GetPagination(index)

	ctx.HTML(200, "pagination.html", pages)
}
