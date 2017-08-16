// +build linux windows
// +build 386 amd64

package main

import (
	"log"

	"github.com/msepp/stopwatch/bootstrap"
)

// Various handles that are used globally
var gState = struct {
	app *bootstrap.App
	db  *StopwatchDB
}{}

func main() {

	// Init new application
	gState.app = bootstrap.New(Asset, RestoreAsset, HandleGUIMessage)

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
