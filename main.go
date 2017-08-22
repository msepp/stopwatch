// +build linux windows
// +build 386 amd64

package main

import (
	"log"

	app "github.com/msepp/stopwatch/stopwatchapp"
	"github.com/msepp/stopwatch/stopwatchdb"
)

// Various handles that are used globally
var gState = struct {
	app *app.App
	db  *stopwatchdb.StopwatchDB
}{}

func main() {

	// Init new application
	gState.app = app.New(Asset, RestoreAsset, HandleGUIMessage)

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
