package stopwatchmodel

import (
	"encoding/json"
	"time"
)

// TaskDuration contains time spent on a task in nanoseconds
type TaskDuration struct {
	time.Duration
}

// MarshalJSON transforms the duration into JSON string
func (td TaskDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(td.String())
}

// UnmarshalJSON converts given bytes into a duration from JSON
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

// Add adds given duration to receiver
func (td *TaskDuration) Add(d time.Duration) time.Duration {
	td.Duration = td.Duration + d
	return td.Duration
}
