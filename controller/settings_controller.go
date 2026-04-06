package controller

import (
	"imap-sync/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetSettings(ctx *gin.Context) {
	settings, err := internal.GetSettings()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	providers := internal.Providers

	ctx.JSON(http.StatusOK, gin.H{
		"settings":  settings,
		"providers": providers,
	})
}

func HandleUpdateSettings(ctx *gin.Context) {
	var settings internal.Settings
	if err := ctx.ShouldBindJSON(&settings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := internal.UpdateSettings(&settings); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
