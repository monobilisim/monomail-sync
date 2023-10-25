package controller

import (
	"imap-sync/internal"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleGetLog(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		ctx.HTML(http.StatusOK, "log_window.html", gin.H{"log": "invalid id"})
		return
	}
	task := internal.GetTaskFromID(id)
	if task == nil {
		ctx.HTML(http.StatusOK, "log_window.html", gin.H{"log": "invalid id"})
		return
	}
	logfile, err := internal.GetLogFromTask(task)
	if err != nil {
		ctx.HTML(http.StatusOK, "log_window.html", gin.H{"log": "failed to get log"})
		return
	}

	startedAt := time.Unix(task.StartedAt, 0)
	endedAt := time.Unix(task.EndedAt, 0)
	ctx.HTML(200, "log_window.html", gin.H{"log": logfile, "start": startedAt.Local(), "end": endedAt.Local()})
}
