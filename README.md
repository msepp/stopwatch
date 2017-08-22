# Stopwatch

A simple software for tracking time spent on different tasks.

Uses [go-astilectron](https://github.com/asticode/go-astilectron) for Golang/Electron binding.

## Features
 * Add tasks, grouped into named groups.
 * Start & Stop timer per task.
 * Edit groups & task details.
 * Get a table of time used per cost code for a group.

## TODO
 * Editing of recorded time to fix mishaps (eg. forgot to stop task).
  * Already in the CLI.

## Screenshot

![Group view](https://raw.githubusercontent.com/msepp/stopwatch/master/screenshot.png "Group view screenshot")

## Requirements

 * You should have [Go](https://golang.org) installed and set up.
 * Uses Makefiles for build automation, so you should have `make` installed.
 * `wget` is required for automatically downloading vendored packages
 * [asar](https://github.com/electron/asar) is used for bundling UI resources
   * `npm install -g asar`
 * [angular-cli](https://github.com/angular/angular-cli) for UI scaffolding.
   * See [angular-cli](https://github.com/angular/angular-cli) for installation.
 * [go-bindata](https://github.com/lestrrat/go-bindata) is used for packing binary files into executable.
   * `go get -u github.com/lestrrat/go-bindata/...`
 * [go-homedir](https://github.com/mitchellh/go-homedir)
   * `go get -u github.com/mitchellh/go-homedir`
 * [go-astilectron](https://github.com/asticode/go-astilectron)
   * `go get -u github.com/asticode/go-astilectron`
 * Windows builds tested using git bash.

## Building

Make sure you have the requirements installed and run the following commands:

```sh
go get -u github.com/msepp/stopwatch/...
cd $GOPATH/src/github.com/msepp/stopwatch
make
```
This will build the sample for your OS/Arch, if supported.

Building non-host targets happens with `make stopwatch-GOOS-ARCH[.exe]`, where GOOS, ARCH should be replaced with target values.
