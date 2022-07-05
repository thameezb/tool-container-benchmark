package benchmark

import "time"

const charSet = "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"

const (
	CreateProject string = "CreateProject"
	UpdateProject string = "UpdateProject"
	DeleteProject string = "DeleteProject"
	RunBenchmark  string = "RunBenchmark"
)

type Msg struct {
	GitlabID int `json:"gitlabId"`
}

type Time struct {
	projectCreationStartTime time.Time
	projectCreationEndTime   time.Time
	eventStartTime           time.Time
	eventEndTime             time.Time
}
