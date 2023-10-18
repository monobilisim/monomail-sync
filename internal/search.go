package internal

import (
	"strconv"
	"sync"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

func GetSearchData(searchQuery string) PageData {
	results := searchInQueue(searchQuery)

	data := PageData{
		Index: 1,
		Tasks: results,
	}
	return data
}

func searchInQueue(searchQuery string) []*Task {
	var results []*Task
	chunkSize := 150                                       // number of tasks to process in each chunk
	numChunks := (queue.Len() + chunkSize - 1) / chunkSize // round up division
	chunkResults := make([][]*Task, numChunks)
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
					task := e.Value.(*Task)
					if fuzzy.Match(searchQuery, strconv.Itoa(task.ID)) ||
						fuzzy.Match(searchQuery, task.SourceAccount) ||
						fuzzy.Match(searchQuery, task.SourceServer) ||
						fuzzy.Match(searchQuery, task.DestinationAccount) ||
						fuzzy.Match(searchQuery, task.DestinationServer) ||
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
		if len(results) > 150 {
			break
		}
	}

	return results
}

func searchExactCredentials(sourceDetails, destinationDetails Credentials) []*Task {
	var results []*Task
	for i, e := 0, queue.Front(); i < queue.Len() && e != nil; i, e = i+1, e.Next() {
		task := e.Value.(*Task)
		if task.SourceAccount == sourceDetails.Account &&
			task.SourceServer == sourceDetails.Server &&
			task.SourcePassword == sourceDetails.Password &&
			task.DestinationAccount == destinationDetails.Account &&
			task.DestinationServer == destinationDetails.Server &&
			task.DestinationPassword == destinationDetails.Password {
			results = append(results, task)
		}
	}
	return results
}
