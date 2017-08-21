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

	} else {
		if end, err = time.Parse(dateFmt, endStr); err != nil {
			log.Fatalf("Invalid end date: %s", err)
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

	usage, err := db.GetUsage(groupID, start, end)
	if err != nil {
		log.Fatalf("dumping failed: %s", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(usage)
}
