package internal

import (
	"container/list"
)

type Credentials struct {
	Server   string
	Account  string
	Password string
}

var queue *list.List
var taskChan chan Task

const PageSize = 20

func GetQueueData(index int) PageData {
	if queue.Len() == 0 {
		return PageData{}
	}

	tasks := getPageByIndex(index)

	data := PageData{
		Index: index,
		Tasks: tasks,
	}

	return data
}

func getPageByIndex(index int) []*Task {
	var tasks []*Task
	start := (index - 1) * PageSize
	end := start + PageSize

	for i, e := 0, queue.Front(); i < end && e != nil; i, e = i+1, e.Next() {
		if i >= start {
			tasks = append(tasks, e.Value.(*Task))
		}
	}
	return tasks
}

func InitQueue() {
	SetupLogger()
	queue = list.New()
	taskChan = make(chan Task)
	for i := 0; i < 10; i++ {
		addOneTask()
	}
	go processPendingTasks()
}

func addOneTask() {
	task := &Task{
		ID:                 queue.Len() + 1,
		SourceAccount:      "jomo",
		SourceServer:       "imap.gmail.com",
		DestinationAccount: "emin",
		DestinationServer:  "imap.yandex.com",
		Status:             "Pending",
	}
	queue.PushFront(task)
	go func() {
		taskChan <- *task
	}()
}

func AddTask(sourceDetails, destinationDetails Credentials) {
	task := &Task{
		ID:                  queue.Len() + 1,
		SourceAccount:       sourceDetails.Account,
		SourceServer:        sourceDetails.Server,
		SourcePassword:      sourceDetails.Password,
		DestinationAccount:  destinationDetails.Account,
		DestinationServer:   destinationDetails.Server,
		DestinationPassword: destinationDetails.Password,
		Status:              "Pending",
	}

	queue.PushFront(task)
	go func() {
		taskChan <- *task
	}()
}
