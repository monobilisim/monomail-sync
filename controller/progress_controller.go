package controller

import (
	"imap-sync/internal"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleTaskProgress(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	progress := internal.GetTaskProgress(id)
	if progress == nil {
		ctx.JSON(http.StatusOK, gin.H{"active": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"active":         true,
		"percent":        progress.Percent,
		"current_folder": progress.CurrentFolder,
		"folders_done":   progress.FoldersDone,
		"folders_total":  progress.FoldersTotal,
		"msgs_done":      progress.MsgsDone,
		"msgs_total":     progress.MsgsTotal,
		"speed":          progress.Speed,
		"bytes_copied":   progress.BytesCopied,
		"eta_seconds":    progress.EtaSeconds,
	})
}
