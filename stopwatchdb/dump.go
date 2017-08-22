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

// Usage is the time used for a date
type Usage struct {
	Date string
	Used model.TaskDuration
}

// CostUsage is time used per cost code over a period of time
type CostUsage struct {
	CostCode string
	Usage    []Usage
	Total    model.TaskDuration
}

// UsageReport is a time usage report. Contains work for one group.
type UsageReport struct {
	// Array with the dates in the report
	Dates []Usage
	// CostCodes contains the time used per cost code
	CostCodes []CostUsage
	// Combined is the combined total time used
	Combined model.TaskDuration
}

func (db *StopwatchDB) GetUsage(group int, start, end time.Time) (*UsageReport, error) {
	daily := map[string]map[string]model.TaskDuration{}
	total := map[string]model.TaskDuration{}
	dates := []Usage{}
	combined := model.TaskDuration{}

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
		dates = append(dates, Usage{Date: start.Format(dateFmt)})
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
						daily[task.CostCode][d.Date] = model.TaskDuration{}
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

				combined.Add(dur)
			}

			// Iterate slices until we hit end.
			return nil
		}); err != nil {
			return nil, err
		}
	}

	rep := UsageReport{
		Dates:     dates,
		CostCodes: []CostUsage{},
		Combined:  combined,
	}

	// Transform result for easier use in UI
	for cost, usage := range daily {
		// Omit tasks that have no time recorded
		if total[cost].Duration == 0 {
			continue
		}

		c := CostUsage{CostCode: cost, Total: total[cost], Usage: []Usage{}}

		for di, date := range rep.Dates {
			c.Usage = append(c.Usage, Usage{Date: date.Date, Used: usage[date.Date]})
			rep.Dates[di].Used.Add(usage[date.Date].Duration)
		}

		rep.CostCodes = append(rep.CostCodes, c)
	}

	return &rep, nil
}
