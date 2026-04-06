package controller

import (
	"fmt"
	"imap-sync/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetAuditLog(ctx *gin.Context) {
	limit := 100
	entries, err := internal.GetAuditLog(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, entries)
}

func HandleGetSessions(ctx *gin.Context) {
	sessions, err := internal.GetActiveSessions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, sessions)
}

func HandleTerminateSession(ctx *gin.Context) {
	sessionID := ctx.Param("id")
	var id int
	_, err := fmt.Sscanf(sessionID, "%d", &id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	if err := internal.TerminateSession(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func HandleTerminateAllSessions(ctx *gin.Context) {
	if err := internal.TerminateAllSessions(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	internal.AddAuditLog("admin", "session_terminate_all", "All sessions terminated", ctx.ClientIP())

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
