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
	StartedAt           int64
	EndedAt             int64
	LogFile             string
}

type PageData struct {
	Index int
	Tasks []*Task
}
