package controller

import (
	"imap-sync/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetStats(ctx *gin.Context) {
	stats, err := internal.GetDashboardStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, stats)
}
