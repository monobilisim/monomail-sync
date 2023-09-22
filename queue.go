package main

import (
	"container/list"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID      int
	account string
	server  string
	status  string
}

type Page struct {
	Number int
	Active bool
}

var queue *list.List

const PageSize = 20

func initQueue() {
	queue = list.New()
	for i := 0; i < 140; i++ {
		task := Task{
			ID:      i + 1,
			account: "imap.gmail.com",
			server:  "jomo",
			status:  "In progress",
		}
		queue.PushFront(task)
	}
}

func addOneTask() {
	task := Task{
		ID:      queue.Len() + 1,
		account: "imap.gmail.com",
		server:  "jomo",
		status:  "In progress",
	}
	queue.PushFront(task)
}

// Updates pagination buttons
func handlePagination(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))
	pages := []Page{}
	startPage := index - 2
	endPage := index + 2

	if startPage < 1 {
		startPage = 1
	}

	if endPage > queue.Len()/PageSize {
		endPage = queue.Len() / PageSize
	}

	for i := startPage; i <= endPage; i++ {
		pages = append(pages, Page{
			Number: i,
			Active: i == index,
		})
	}

	ctx.HTML(200, "pagination.html", pages)
}

// Updates the page number on tbody
func handleQueuePoll(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")

	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	page := getPageByIndex(index)

	var rows string
	rows += fmt.Sprintf(`<thead>
                    <tr>
                        <th>Index</th>
                        <th>Server</th>
                        <th>Account</th>
                        <th>Status</th>
                    </tr>
                </thead>
                <tbody id="table-body" hx-get="/api/queue?page=%d" hx-swap="innerHTML"
                    hx-trigger="every 4s">`, index)

	for _, task := range page {
		rows += `<tr>
					<th>` + strconv.Itoa(task.ID) + `</th>
					<td>` + task.server + `</td>
					<td>` + task.account + `</td>
					<td>
						<span class="badge badge-outline-warning">` + task.status + `</span>
					</td>
				</tr>`
	}

	rows += `</tbody>`

	ctx.String(200, rows)
}

// Send only one page
func handleQueue(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	if queue.Len() == 0 {
		ctx.String(200, "")
		return
	}

	var rows string

	page := getPageByIndex(index)
	for _, task := range page {
		rows += `<tr>
					<th>` + strconv.Itoa(task.ID) + `</th>
					<td>` + task.server + `</td>
					<td>` + task.account + `</td>
					<td>
						<span class="badge badge-outline-warning">` + task.status + `</span>
					</td>
				</tr>`
	}

	ctx.String(200, rows)
}

func getPageByIndex(index int) []Task {
	var tasks []Task
	start := (index - 1) * PageSize
	end := start + PageSize

	for i, e := 0, queue.Front(); i < end && e != nil; i, e = i+1, e.Next() {
		if i >= start {
			tasks = append(tasks, e.Value.(Task))
		}
	}

	return tasks
}
