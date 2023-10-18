package internal

func RetryTask(task *Task) {
	newTask := *task
	newTask.Status = "Pending"
	newTask.ID = queue.Len() + 1
	queue.PushFront(&newTask)
	AddTaskToDB(task)
	go func() {
		taskChan <- newTask
	}()
}
