package main

// Config is the date returned from RequestConfig
type Config struct {
	Groups     []Group `json:"groups"`
	ActiveTask *Task   `json:"activeTask,omitempty"`
}
