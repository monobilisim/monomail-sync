package internal

type Task struct {
	ID      int
	Account string
	Server  string
	Status  string
}

type Pagination struct {
	Number int
	Active bool
}

type PageData struct {
	Index int
	Tasks []Task
}
