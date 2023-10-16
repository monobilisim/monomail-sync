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
	LogFile             string
}

type PageData struct {
	Index int
	Tasks []*Task
}
