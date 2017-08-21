package stopwatchdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/msepp/stopwatch/model"
)

const dateFmt = "2006-01-02"

// Usage is a time usage report. Contains work for one group.
type Usage struct {
	// Daily contains daily usage per cost code.
	Daily map[string]map[string]model.TaskDuration
	// Contains total usage per cost code.
	Total map[string]model.TaskDuration
}

func (db *StopwatchDB) GetUsage(group int, start, end time.Time) (*Usage, error) {
	daily := map[string]map[string]model.TaskDuration{}
	total := map[string]model.TaskDuration{}
	dates := []string{}

	tasks := []*model.Task{}

	// Normalize dates to begin of start date and end of end date.
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day()+1, 0, 0, 0, 0, end.Location())
	end = end.Add(-1 * time.Second)

	// Handle in UTC
	start = start.UTC()
	end = end.UTC()

	// Check order
	if start.Equal(end) || start.After(end) {
		return nil, errors.New("start must be a time before end.")
	}

	// open group tasks
	if err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketTasks))
		bg := b.Bucket(Itob(group))
		if bg == nil {
			return errors.New("group not found")
		}

		return bg.ForEach(func(k []byte, v []byte) error {
			var t model.Task
			json.Unmarshal(v, &t)
			tasks = append(tasks, &t)
			return nil
		})
	}); err != nil {
		return nil, err
	}

	// Minimum and maximum timestamps.
	min := []byte(start.Format(time.RFC3339))
	max := []byte(end.Format(time.RFC3339))

	// Generate dates to result.
	for start.Before(end) {
		dates = append(dates, start.Format(dateFmt))
		start = start.AddDate(0, 0, 1)
	}

	// Go through each task, day by day.
	for _, task := range tasks {
		if err := db.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BucketSlices))
			bs := b.Bucket(bytes.Join([][]byte{
				Itob(task.GroupID),
				Itob(task.ID)},
				[]byte("-"),
			))
			if bs == nil {
				log.Printf("Unable to find slices for task %d:%s", task.ID, task.Name)
				return nil
			}

			// Seek to start date
			c := bs.Cursor()
			for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) < 0; k, v = c.Next() {
				starttime, _ := time.Parse(time.RFC3339, string(k))
				date := starttime.Format(dateFmt)

				if v == nil {
					continue
				}

				endtime, _ := time.Parse(time.RFC3339, string(v))
				if endtime.IsZero() {
					continue
				}

				dur := endtime.Sub(starttime)

				log.Printf("%s/%s: %s", task.CostCode, task.Name, endtime)

				if _, ok := daily[task.CostCode]; !ok {
					daily[task.CostCode] = map[string]model.TaskDuration{}
					for _, d := range dates {
						daily[task.CostCode][d] = model.TaskDuration{}
					}
				}

				if _, ok := total[task.CostCode]; !ok {
					total[task.CostCode] = model.TaskDuration{}
				}

				od := daily[task.CostCode][date]
				od.Add(dur)
				daily[task.CostCode][date] = od

				od = total[task.CostCode]
				od.Add(dur)
				total[task.CostCode] = od
			}

			// Iterate slices until we hit end.
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return &Usage{
		Daily: daily,
		Total: total,
	}, nil
}
