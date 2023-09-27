package internal

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleQueuePolling(ctx *gin.Context) {
	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	if queue.Len() == 0 {
		ctx.String(200, "")
		return
	}

	tasks := getPageByIndex(index)

	data := PageData{
		Index: index,
		Tasks: tasks,
	}

	ctx.HTML(200, "tbody.html", data)
}
