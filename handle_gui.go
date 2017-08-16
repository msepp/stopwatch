package main

import (
	"errors"
	"fmt"
	"path"

	"github.com/msepp/stopwatch/bootstrap"
	"github.com/msepp/stopwatch/message"
)

// HandleGUIMessage is called when we receive messages from the user interface.
func HandleGUIMessage(msg *message.Message) (interface{}, error) {
	switch msg.Key {
	case message.RequestAppVersions:
		return HandleGetAppVersions(msg)

	case message.RequestOpenDatabase:
		return HandleOpenDatabase(msg)

	case message.RequestConfig:
		return HandleGetConfig(msg)

	case message.RequestGroupTasks:
		return HandleGetGroupTasks(msg)

	case message.RequestAddGroup:
		return HandleAddGroup(msg)

	case message.RequestAddTask:
		return HandleAddTask(msg)

	case message.RequestStartTask:
		return HandleStartTask(msg)

	case message.RequestStopTask:
		return HandleStopTask(msg)

	default:
		return nil, fmt.Errorf("Unrecognized GUI message key: %s", msg.Key)
	}
}

// HandleGetAppVersions returns versions used in the updater
func HandleGetAppVersions(msg *message.Message) (interface{}, error) {
	return map[string]string{
		"app":         bootstrap.Version(),
		"build":       bootstrap.Build(),
		"electron":    bootstrap.ElectronVersion(),
		"astilectron": bootstrap.AstilectronVersion(),
	}, nil
}

// HandleOpenDatabase attempts to open database. Returns if the operation
// succeeded.
func HandleOpenDatabase(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		gState.db = NewStopwatchDB()
	}

	rootdir, err := bootstrap.PersistentDataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to create database dir: %s", err)
	}

	if err := gState.db.Open(path.Join(rootdir, "data.dat")); err != nil {
		gState.db = nil
		return nil, fmt.Errorf("failed to open database: %s", err)
	}

	return nil, nil
}

// HandleGetConfig returns app config and last known state
func HandleGetConfig(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	// Retrieve groups
	groups, err := gState.db.ReadGroups()
	if err != nil {
		return nil, fmt.Errorf("Unable to read groups: %s", err)
	}

	// Get active task
	at, err := gState.db.GetActiveTask()
	if err != nil {
		return nil, fmt.Errorf("Unable to read active task: %s", err)
	}

	return &Config{
		Groups:     groups,
		ActiveTask: at,
	}, nil
}

// HandleGetGroupTasks returns list of tasks for a group
func HandleGetGroupTasks(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var payload ReqPayloadGetGroupTasks
	if err := msg.Into(&payload); err != nil || payload.GroupID <= 0 {
		return nil, fmt.Errorf("payload invalid or missing: %s", err)
	}

	// Retrieve tasks
	tlist, err := gState.db.ReadTasks(payload.GroupID)
	if err != nil {
		return nil, fmt.Errorf("Unable to read tasks: %s", err)
	}

	return tlist, nil
}

// HandleAddGroup adds the group detailed in msg data
func HandleAddGroup(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var payload ReqPayloadAddGroup
	if err := msg.Into(&payload); err != nil || payload.Name == "" {
		return nil, fmt.Errorf("payload invalid or missing: %s", err)
	}

	p, err := gState.db.AddGroup(payload.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to add group: %s", err)
	}

	return p, nil
}

// HandleAddTask adds a task for a group using details from msg data
func HandleAddTask(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var payload ReqPayloadAddTask
	if err := msg.Into(&payload); err != nil {
		return nil, fmt.Errorf("payload invalid: %s", err)
	}

	if payload.Name == "" {
		return nil, fmt.Errorf("name must be given and not empty")
	}
	if payload.GroupID <= 0 {
		return nil, fmt.Errorf("groupid must be a non-zero positive integer")
	}

	t, err := gState.db.AddTask(payload.GroupID, payload.Name, payload.CostCode)
	if err != nil {
		return nil, fmt.Errorf("failed to add task: %s", err)
	}

	return t, nil
}

// HandleStartTask starts a task
func HandleStartTask(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var payload ReqPayloadSetTaskStatus
	if err := msg.Into(&payload); err != nil {
		return nil, fmt.Errorf("payload invalid: %s", err)
	}

	if payload.GroupID <= 0 || payload.TaskID <= 0 {
		return nil, fmt.Errorf("group and task IDs must be non-zero positive integers")
	}

	// Locate task
	task, err := gState.db.GetTask(payload.GroupID, payload.TaskID)
	if err != nil {
		return nil, fmt.Errorf("failure starting task: %s", err)
	}

	activeTask, err := gState.db.GetActiveTask()
	if err != nil {
		return nil, fmt.Errorf("Unable test active task: %s", err)
	}

	// Stop active task if not the same as current task
	if activeTask != nil {
		if activeTask.ID == task.ID && activeTask.GroupID == task.GroupID {
			if activeTask.Running != nil {
				return activeTask, nil
			}
		} else {
			if _, err = gState.db.StopTask(activeTask.GroupID, activeTask.ID); err != nil {
				return nil, fmt.Errorf("unable to stop current task: %s", err)
			}
		}
	}

	if err = gState.db.SetActiveTask(task.GroupID, task.ID); err != nil {
		return nil, fmt.Errorf("unable to set active task: %s", err)
	}

	// Mark running
	return gState.db.StartTask(task.GroupID, task.ID)
}

// HandleStopTask stops a task
func HandleStopTask(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var payload ReqPayloadSetTaskStatus
	if err := msg.Into(&payload); err != nil {
		return nil, fmt.Errorf("payload invalid: %s", err)
	}

	if payload.GroupID <= 0 || payload.TaskID <= 0 {
		return nil, fmt.Errorf("group and task IDs must be non-zero positive integers")
	}

	// Locate task
	task, err := gState.db.GetTask(payload.GroupID, payload.TaskID)
	if err != nil {
		return nil, fmt.Errorf("failure stopping task: %s", err)
	}

	// Stop active task if it matches current task
	activeTask, err := gState.db.GetActiveTask()
	if err != nil {
		return nil, fmt.Errorf("Unable test active task: %s", err)
	}

	// Stop active task if it matches the given task
	if activeTask != nil {
		if activeTask.ID == task.ID && activeTask.GroupID == task.GroupID {
			if err := gState.db.SetActiveTask(0, 0); err != nil {
				return nil, fmt.Errorf("unable to clear active task: %s")
			}
		}
	}

	// Mark running
	return gState.db.StopTask(task.GroupID, task.ID)
}
