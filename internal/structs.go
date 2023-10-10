package internal

type Task struct {
	ID                  int
	SourceAccount       string
	SourceServer        string
	SourcePassword      string
	DestinationAccount  string
	DestinationServer   string
	DestinationPassword string
	Status              string
}

type PageData struct {
	Index int
	Tasks []Task
}
