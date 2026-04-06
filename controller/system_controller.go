package controller

import (
	"imap-sync/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetSystemInfo(ctx *gin.Context) {
	info := internal.GetSystemInfo()
	ctx.JSON(http.StatusOK, info)
}
