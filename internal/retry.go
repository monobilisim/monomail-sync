package internal

func RetryTask(task *Task) {
	if err := updateTaskStatus(task, "Pending"); err != nil {
		log.Errorf("Failed to add task: %s", err)
	}
	go func() {
		taskChan <- *task
	}()
}
