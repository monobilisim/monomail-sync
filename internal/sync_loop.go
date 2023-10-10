package internal

import "time"

func processPendingTasks() {
	for {
		// Get the index of the first pending task
		task := getFirstPendingTask()

		// If there are no pending tasks, wait for a new task to be added
		if task == nil {
			task := <-taskChan
			simulateSyncIMAP(&task)
			continue
		}

		simulateSyncIMAP(task)
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

func simulateSyncIMAP(details *Task) error {
	details.Status = "In Progress"
	time.Sleep(5000 * time.Millisecond)
	details.Status = "Done"
	return nil
}
