package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func handleSearch(ctx *gin.Context) {
	searchQuery := ctx.PostForm("search-input")

	if searchQuery == "" {
		handleQueue(ctx)
		return
	}

	results := searchInQueue(searchQuery)

	data := PageData{
		Index: 1,
		Tasks: results,
	}

	ctx.HTML(200, "queue.html", data)
}

// Fuzzy match each field in the queue
func searchInQueue(searchQuery string) []Task {
	var results []Task
	for e := queue.Front(); e != nil; e = e.Next() {
		task := e.Value.(Task)
		if fuzzy.Match(searchQuery, fmt.Sprint(task.ID)) ||
			fuzzy.Match(searchQuery, task.Account) ||
			fuzzy.Match(searchQuery, task.Server) ||
			fuzzy.Match(searchQuery, task.Status) {
			results = append(results, task)
		}
	}
	return results
}
