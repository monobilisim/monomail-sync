package internal

import "time"

func getFirstInProgressTaskIndex() int {
	// Do it from the back of the queue
	for i, e := queue.Back(), 0; i != nil; i, e = i.Prev(), e+1 {
		if i.Value.(Task).Status == "In Progress" {
			return queue.Len() - e
		}
	}
	return -1
}

// Infinite loop waiting for in progress tasks to appear in the queue and then calling the sync function on them and setting their status to "Done" when they are finished.
func syncLoop() {
	for {
		// Wait for a task to appear in the queue
		for queue.Len() == 0 {
			time.Sleep(1 * time.Second)
		}

		// Get the index of the first task in the queue that is in progress
		index := getFirstInProgressTaskIndex()
		if index == -1 {
			continue
		}

		// Get the task from the queue
		task := getTaskByIndex(index)

		// Set the task status to "In Progress"
		task.Status = "In Progress"
		setTaskByIndex(index, task)

		// Sync the task
		syncIMAP(task)

		// Set the task status to "Done"
		task.Status = "Done"
		setTaskByIndex(index, task)
	}
}

func getTaskByIndex(index int) Task {
	for i, e := 0, queue.Front(); i < index && e != nil; i, e = i+1, e.Next() {
		if i == index {
			return e.Value.(Task)
		}
	}
	return Task{}
}

func setTaskByIndex(index int, task Task) {
	for i, e := 0, queue.Front(); i < index && e != nil; i, e = i+1, e.Next() {
		if i == index {
			e.Value = task
		}
	}
}
