package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/msepp/stopwatch/stopwatchdb"
)

const dateFmt = "2006-01-02"

var start time.Time
var end time.Time
var groupID int
var taskID int

func main() {
	var dbPath string
	var startStr string
	var endStr string
	var dumpType string
	var err error

	flag.StringVar(&dbPath, "db", "data.dat", "path to database")
	flag.StringVar(&startStr, "start", "", "start date (YYYY-MM-DD for reports, RFC 3339 for slices). Defaults to start of current day for reports.")
	flag.StringVar(&endStr, "end", "", "end date (YYYY-MM-DD for reports, RFC 3339 for slices). Defaults to now for reports.")
	flag.IntVar(&groupID, "groupID", 0, "Group ID to dump/modify")
	flag.IntVar(&taskID, "taskID", 0, "Task ID to dump/modify")
	flag.StringVar(&dumpType, "type", "report", "operation type. 'slices' returns recorded slices, 'report' gives a nice report, 'setslice' allows setting a slice and 'rmslice' removes slice.")
	flag.Parse()

	if groupID <= 0 {
		log.Fatalf("groupID needs to be a positive non-zero integer")
	}

	now := time.Now()

	switch dumpType {
	case "setslice", "rmslice":
		if start, err = time.Parse(time.RFC3339, startStr); err != nil {
			log.Fatalf("Invalid start datetime: %s", err)
		}

		if dumpType != "rmslice" {
			if end, err = time.Parse(time.RFC3339, endStr); err != nil {
				log.Fatalf("Invalid end datetime: %s", err)
			}
		}

	default:
		if startStr == "" {
			start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		} else {
			if start, err = time.Parse(dateFmt, startStr); err != nil {
				log.Fatalf("Invalid start date: %s", err)
			}
		}

		if endStr == "" {
			end = time.Now()

		} else {
			if end, err = time.Parse(dateFmt, endStr); err != nil {
				log.Fatalf("Invalid end date: %s", err)
			}
		}
	}

	if st, err := os.Stat(dbPath); err != nil || st.IsDir() {
		log.Fatalf("Invalid database path")
	}

	// Open database
	db := stopwatchdb.New()
	if err = db.Open(dbPath); err != nil {
		log.Fatalf("Opening database failed: %s", err)
	}

	var result interface{}
	switch dumpType {
	case "report":
		result, err = db.GetUsage(groupID, start, end)
		if err != nil {
			log.Fatalf("dumping failed: %s", err)
		}

	case "slices":
		result, err = db.GetSlices(groupID, start, end)
		if err != nil {
			log.Fatalf("dumping failed: %s", err)
		}

	case "setslice":
		result, err = db.SetSlice(groupID, taskID, start, end)
		if err != nil {
			log.Fatalf("set slice failed: %s", err)
		}

	case "rmslice":
		result, err = db.RemoveSlice(groupID, taskID, start)
		if err != nil {
			log.Fatalf("removing slice failed: %s", err)
		}

	default:
		log.Fatalf("Invalid dump type")
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(result)
}
