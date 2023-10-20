package internal

func CancelTask(task *Task) {
	if task.Status != "In Progress" {
		updateTaskStatus(task, "Cancelled")
		log.Debug(len(taskChan))
		_ = <-taskChan
		log.Debug(len(taskChan))
		return
	}

}
