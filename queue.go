package main

import (
	"container/list"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID      int
	Account string
	Server  string
	Status  string
}

type Pagination struct {
	Number int
	Active bool
}

type PageData struct {
	Index int
	Tasks []Task
}

var queue *list.List

const PageSize = 20

func handleQueuePoll(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
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

// Send only one page
func handleQueue(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))

	addOneTask()

	if queue.Len() == 0 {
		ctx.String(200, "")
		return
	}

	tasks := getPageByIndex(index)

	data := PageData{
		Index: index,
		Tasks: tasks,
	}

	ctx.HTML(200, "queue.html", data)
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
	for i := 0; i < 140; i++ {
		task := Task{
			ID:      i + 1,
			Account: "jomo",
			Server:  "imap.gmail.com",
			Status:  "In progress",
		}
		queue.PushFront(task)
	}
}

func addOneTask() {
	task := Task{
		ID:      queue.Len() + 1,
		Account: "jomo",
		Server:  "imap.gmail.com",
		Status:  "In progress",
	}
	queue.PushFront(task)
}

// Updates pagination buttons
func handlePagination(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	index, _ := strconv.Atoi(ctx.Request.FormValue("page"))
	pages := []Pagination{}
	startPage := index - 2
	endPage := index + 2

	if startPage < 1 {
		startPage = 1
	}

	if endPage > queue.Len()/PageSize {
		endPage = queue.Len() / PageSize
	}

	if (index <= 2 || index >= endPage-2 || index == endPage || index == endPage-1) && endPage-startPage+1 < 5 && endPage < queue.Len()/PageSize {
		endPage = startPage + 4
	}

	for i := startPage; i <= endPage; i++ {
		pages = append(pages, Pagination{
			Number: i,
			Active: i == index,
		})
	}

	ctx.HTML(200, "pagination.html", pages)
}
