package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID      int
	account string
	server  string
	status  string
}

var queue []Task

func handleQueue(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	log.Infof("Queueing task")

	task := Task{
		ID:      len(queue) + 1,
		account: "imap.gmail.com",
		server:  "jomo",
		status:  "In progress",
	}
	queue = append(queue, task)

	row := fmt.Sprintf(`<tr>
					<th>%d</th>
					<td>%s</td>
					<td>%s</td>
					<td>
						<span class="badge badge-outline-warning">%s</span>
					</td>
				</tr>`, task.ID, task.server, task.account, task.status)

	ctx.String(200, row)
}
