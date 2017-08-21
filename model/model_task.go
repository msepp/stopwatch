package model

import (
	"time"
)

// Task describes a single task
type Task struct {
	ID       int          `json:"id"`
	GroupID  int          `json:"groupid"`
	Name     string       `json:"name"`
	CostCode string       `json:"costcode"`
	Used     TaskDuration `json:"duration"`
	Running  *time.Time   `json:"running,omitempty"`
}

// NewTask initializes an task
func NewTask(groupid int, name, costcode string) *Task {
	return &Task{
		ID:       0,
		Name:     name,
		CostCode: costcode,
		GroupID:  groupid,
	}
}
