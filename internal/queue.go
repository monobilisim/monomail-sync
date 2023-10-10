package internal

import (
	"container/list"
	"strconv"

	"github.com/gin-gonic/gin"
)

var queue *list.List

const PageSize = 20

// Send only one page
func handleQueue(ctx *gin.Context) {
	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	go addOneTask()

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

func initQueue() {
	queue = list.New()
	for i := 0; i < 2000; i++ {
		addOneTask()
	}
}

func addOneTask() {
	task := Task{
		ID:                 queue.Len() + 1,
		SourceAccount:      "jomo",
		SourceServer:       "imap.gmail.com",
		DestinationAccount: "emin",
		DestinationServer:  "imap.yandex.com",
		Status:             "In progress",
	}
	queue.PushFront(task)
}
