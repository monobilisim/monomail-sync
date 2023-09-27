package internal

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func handleSearch(ctx *gin.Context) {
	searchQuery := ctx.PostForm("search-input")

	if searchQuery == "" {
		handleQueue(ctx)
		return
	}

	results := searchChunk(searchQuery)

	data := PageData{
		Index: 1,
		Tasks: results,
	}

	ctx.HTML(200, "queue.html", data)
}

func searchInQueue(searchQuery string) []Task {
	var results []Task
	chunkSize := 100                                       // number of tasks to process in each chunk
	numChunks := (queue.Len() + chunkSize - 1) / chunkSize // round up division
	chunkResults := make([][]Task, numChunks)
	var wg sync.WaitGroup
	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * chunkSize
			end := start + chunkSize
			if end > queue.Len() {
				end = queue.Len()
			}
			for j, e := 0, queue.Front(); j < end && e != nil; j, e = j+1, e.Next() {
				if j >= start {
					task := e.Value.(Task)
					if fuzzy.Match(searchQuery, strconv.Itoa(task.ID)) ||
						fuzzy.Match(searchQuery, task.Account) ||
						fuzzy.Match(searchQuery, task.Server) ||
						fuzzy.Match(searchQuery, task.Status) {
						chunkResults[i] = append(chunkResults[i], task)
					}
				}
			}
		}(i)

		if len(results) > 150 {
			break
		}
	}

	wg.Wait()
	for _, chunkResult := range chunkResults {
		results = append(results, chunkResult...)
	}
	return results
}

// Fuzzy match each field in the queue
func searchChunk(searchQuery string) []Task {
	var results []Task
	for e := queue.Front(); e != nil; e = e.Next() {
		task := e.Value.(Task)
		if fuzzy.Match(searchQuery, strconv.Itoa(task.ID)) ||
			fuzzy.Match(searchQuery, task.Account) ||
			fuzzy.Match(searchQuery, task.Server) ||
			fuzzy.Match(searchQuery, task.Status) {
			results = append(results, task)
		}

		if len(results) > 100 {
			break
		}
	}
	return results
}
