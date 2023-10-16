package internal

import (
	"os"
	"time"
)

func processPendingTasks() {
	for {
		// Get the index of the first pending task
		task := getFirstPendingTask()

		// If there are no pending tasks, wait for a new task to be added
		if task == nil {
			<-taskChan
			continue
		}

		//syncIMAP(task)
		simulateTask(task)
		time.Sleep(100 * time.Millisecond)
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

func simulateTask(task *Task) {

	currentTime := time.Now().Format("2006.01.02_15:04:05")

	logname := task.SourceAccount + "_" + task.DestinationAccount + "_" + currentTime + ".log"

	updateTaskLogFile(task, logname)

	// Write to log file
	logFile, err := os.Create("LOG_imapsync/" + logname)
	if err != nil {
		log.Error(err)
	}
	defer logFile.Close()

	logFile.WriteString("This is a test log file\n")

	updateTaskStatus(task, "In Progress")
	time.Sleep(10000 * time.Millisecond)
	updateTaskStatus(task, "Done")
}
