package internal

type Task struct {
	ID                 int
	SourceAccount      string
	SourceServer       string
	DestinationAccount string
	DestinationServer  string
	Status             string
}

type Pagination struct {
	Number int
	Active bool
}

type PageData struct {
	Index int
	Tasks []Task
}
