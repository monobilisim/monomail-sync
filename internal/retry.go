package internal

func RetryTask(task *Task) {
	newTask := *task
	newTask.Status = "Pending"
	newTask.ID = queue.Len() + 1
	queue.PushFront(&newTask)
	if err := AddTaskToDB(&newTask); err != nil {
		log.Errorf("Failed to persist retried task: %v", err)
		return
	}
	go func() {
		taskChan <- newTask
	}()
}
