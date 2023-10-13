package internal

import "time"

func processPendingTasks() {
	for {
		// Get the index of the first pending task
		task := getFirstPendingTask()

		// If there are no pending tasks, wait for a new task to be added
		if task == nil {
			<-taskChan
			continue
		}

		syncIMAP(task)
		time.Sleep(2000 * time.Millisecond)
	}
}

// return pointer to first pending task from back of queue
func getFirstPendingTask() *Task {
	for e := queue.Back(); e != nil; e = e.Prev() {
		if task, ok := e.Value.(*Task); ok && task.Status == "Pending" {
			return task
		}
	}
	return nil
}
