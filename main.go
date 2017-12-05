// +build linux windows
// +build 386 amd64

package main

import (
	"flag"
	"log"

	app "github.com/msepp/stopwatch/stopwatchapp"
	"github.com/msepp/stopwatch/stopwatchdb"
)

// Various handles that are used globally
var gState = struct {
	app          *app.App
	db           *stopwatchdb.StopwatchDB
	databasePath string
	workDir      string
}{}

func main() {
	flag.StringVar(&gState.workDir, "work", "", "Working directory. If not set, an attempt is made to guess users home.")
	flag.StringVar(&gState.databasePath, "db", "", "Database path. If none given, a database is created under working dir.")
	// Init
	flag.Parse()

	// Init new application
	gState.app = app.New(Asset, RestoreAsset, HandleGUIMessage)

	// Set app working dir if given
	if gState.workDir != "" {
		gState.app.SetWorkingDir(gState.workDir)
	}

	// Bootstrap to get things going
	if err := gState.app.Bootstrap(); err != nil {
		log.Fatalln(err)
	}

	// Wait for app to exit
	gState.app.Wait()

	// Close database
	if gState.db != nil {
		gState.db.Close()
	}
}
