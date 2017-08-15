package main

import (
	"encoding/json"
	"time"
)

type TaskDuration struct {
	time.Duration
}

func (td TaskDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(td.String())
}

func (td *TaskDuration) UnmarshalJSON(in []byte) error {
	var s string
	if err := json.Unmarshal(in, &s); err != nil {
		return err
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	td.Duration = d
	return nil
}

func (td *TaskDuration) Add(d time.Duration) time.Duration {
	td.Duration = td.Duration + d
	return td.Duration
}

// Task describes a single task
type Task struct {
	ID        int          `json:"id"`
	ProjectID int          `json:"projectid"`
	Name      string       `json:"name"`
	CostCode  string       `json:"costcode"`
	Used      TaskDuration `json:"duration"`
	Running   *time.Time   `json:"running,omitempty"`
}

// NewTask initializes an task
func NewTask(project int, name, costcode string) *Task {
	return &Task{
		ID:        0,
		Name:      name,
		CostCode:  costcode,
		ProjectID: project,
	}
}
