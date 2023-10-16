package internal

import "os"

func GetLogFromTask(task *Task) (string, error) {
	logPath := "LOG_imapsync/" + task.LogFile
	log, err := os.ReadFile(logPath)
	if err != nil {
		return "", err
	}
	return string(log), nil
}

func GetTaskFromID(id int) *Task {
	for e := queue.Back(); e != nil; e = e.Prev() {
		if task, ok := e.Value.(*Task); ok && task.ID == id {
			return task
		}
	}
	return nil
}
