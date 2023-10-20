package internal

// CancelTask cancels a task and removes it from channel
func CancelTask(task *Task) {
	if task.Status != "In Progress" {
		updateTaskStatus(task, "Cancelled")

		select {
		case <-taskChan:
			log.Info("Task cancelled: ", task.ID)
		default:
			log.Info("Task not found in channel: ", task.ID)
		}

	} else {
		cancel()
		updateTaskStatus(task, "Cancelled")
	}
}
