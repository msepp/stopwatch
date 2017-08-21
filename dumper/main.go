package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/msepp/stopwatch/model"
	"github.com/msepp/stopwatch/stopwatchdb"
)

const dateFmt = "2006-01-02"

var start time.Time
var end time.Time
var db *bolt.DB
var groupID int

func main() {
	var dbPath string
	var startStr string
	var endStr string
	var err error

	flag.StringVar(&dbPath, "db", "data.dat", "path to database")
	flag.StringVar(&startStr, "start", "", "start date YYYY-MM-DD. Defaults to start of current day.")
	flag.StringVar(&endStr, "end", "", "end date YYYY-MMM-DD. Defaults to now")
	flag.IntVar(&groupID, "groupID", 0, "Group ID to dump")
	flag.Parse()

	if groupID <= 0 {
		log.Fatalf("groupID needs to be a positive non-zero integer")
	}

	now := time.Now()

	if startStr == "" {
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	} else {
		if start, err = time.Parse(dateFmt, startStr); err != nil {
			log.Fatalf("Invalid start date: %s", err)
		}
	}

	if endStr == "" {
		end = time.Now()
		end = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		end = end.AddDate(0, 0, 1)

	} else {
		if end, err = time.Parse(dateFmt, endStr); err != nil {
			log.Fatalf("Invalid end date: %s", err)
		}
	}

	if end.Before(start) || end.Equal(start) {
		log.Fatalf("End is a date after start or same as start.")
	}

	if st, err := os.Stat(dbPath); err != nil || st.IsDir() {
		log.Fatalf("Invalid database path")
	}

	start = start.UTC()
	end = end.UTC()

	// Open database
	if db, err = bolt.Open(dbPath, 0600, &bolt.Options{
		ReadOnly: true,
	}); err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	if err = dump(); err != nil {
		log.Fatalf("dumping failed: %s", err)
	}
}

func dump() error {
	output := map[string]map[string]model.TaskDuration{}

	tasks := []*model.Task{}

	// open group tasks
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(stopwatchdb.BucketTasks))
		bg := b.Bucket(stopwatchdb.Itob(groupID))
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
		return err
	}

	// Go through each task, day by day.
	for _, task := range tasks {
		if err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(stopwatchdb.BucketSlices))
			bs := b.Bucket(bytes.Join([][]byte{
				stopwatchdb.Itob(task.GroupID),
				stopwatchdb.Itob(task.ID)},
				[]byte("-"),
			))
			if bs == nil {
				log.Printf("Unable to find slices for task %d:%s", task.ID, task.Name)
				return nil
			}

			min := []byte(start.Format(time.RFC3339))
			max := []byte(end.Format(time.RFC3339))

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

				if _, ok := output[task.CostCode]; !ok {
					output[task.CostCode] = map[string]model.TaskDuration{}
				}

				if _, ok := output[task.CostCode][date]; !ok {
					output[task.CostCode][date] = model.TaskDuration{}
				}

				od := output[task.CostCode][date]
				od.Add(dur)
				output[task.CostCode][date] = od
			}

			// Iterate slices until we hit end.
			return nil
		}); err != nil {
			return err
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}
