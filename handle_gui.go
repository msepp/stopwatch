package main

import (
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/msepp/stopwatch/bootstrap"
	"github.com/msepp/stopwatch/message"
	"github.com/msepp/stopwatch/stopwatchdb"
)

// HandleGUIMessage is called when we receive messages from the user interface.
func HandleGUIMessage(msg *message.Message) (interface{}, error) {
	switch msg.Key {

	case message.RequestActiveTask:
		return HandleGetActiveTask(msg)

	case message.RequestAddGroup:
		return HandleAddGroup(msg)

	case message.RequestAddTask:
		return HandleAddTask(msg)

	case message.RequestAppVersions:
		return HandleGetAppVersions(msg)

	case message.RequestGetHistory:
		return HandleGetHistory(msg)

	case message.RequestGroups:
		return HandleGetGroups(msg)

	case message.RequestGroupTasks:
		return HandleGetGroupTasks(msg)

	case message.RequestGetTask:
		return HandleGetTask(msg)

	case message.RequestGetUsage:
		return HandleGetUsage(msg)

	case message.RequestOpenDatabase:
		return HandleOpenDatabase(msg)

	case message.RequestSetHistory:
		return HandleSetHistory(msg)

	case message.RequestStartTask:
		return HandleStartTask(msg)

	case message.RequestStopTask:
		return HandleStopTask(msg)

	case message.RequestUpdateGroup:
		return HandleUpdateGroup(msg)

	case message.RequestUpdateTask:
		return HandleUpdateTask(msg)

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
	if gState.db == nil {
		gState.db = stopwatchdb.New()
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

// HandleGetHistory returns task history
func HandleGetHistory(msg *message.Message) (interface{}, error) {
	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	// Retrieve task usage
	history, err := gState.db.ReadHistory()
	if err != nil {
		return nil, fmt.Errorf("Unable to read history: %s", err)
	}

	return history, nil
}

// HandleSetHistory saves history information
func HandleSetHistory(msg *message.Message) (interface{}, error) {
	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var payload ReqPayloadSetHistory
	if err := msg.Into(&payload); err != nil {
		return nil, fmt.Errorf("payload invalid or missing: %s", err)
	}

	// Save history items
	err := gState.db.SaveHistory(payload.History)
	if err != nil {
		return nil, fmt.Errorf("Unable to save history: %s", err)
	}

	return nil, nil
}

// HandleGetGroups returns all known groups
func HandleGetGroups(msg *message.Message) (interface{}, error) {
	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	// Retrieve groups
	groups, err := gState.db.ReadGroups()
	if err != nil {
		return nil, fmt.Errorf("Unable to read groups: %s", err)
	}

	return groups, nil
}

// HandleGetTask returns details of a single task
func HandleGetTask(msg *message.Message) (interface{}, error) {
	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var payload ReqPayloadGetTask
	if err := msg.Into(&payload); err != nil || payload.TaskID <= 0 || payload.GroupID <= 0 {
		return nil, fmt.Errorf("payload invalid or missing: %s", err)
	}

	// Retrieve task
	task, err := gState.db.ReadTask(payload.GroupID, payload.TaskID)
	if err != nil {
		return nil, fmt.Errorf("Unable to read task: %s", err)
	}

	return task, nil
}

// HandleGetActiveTask returns current active task or nil if not set.
func HandleGetActiveTask(msg *message.Message) (interface{}, error) {
	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	// Get active task
	at, err := gState.db.GetActiveTask()
	if err != nil {
		return nil, fmt.Errorf("Unable to read active task: %s", err)
	}

	return at, nil
}

// HandleGetGroupTasks returns list of tasks for a group
func HandleGetGroupTasks(msg *message.Message) (interface{}, error) {
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

// HandleUpdateGroup updates group value in database
func HandleUpdateGroup(msg *message.Message) (interface{}, error) {
	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var payload ReqPayloadUpdateGroup
	if err := msg.Into(&payload); err != nil {
		return nil, fmt.Errorf("payload invalid: %s", err)
	}

	if payload.GroupID <= 0 {
		return nil, fmt.Errorf("group id must be non-zero positive integer")
	}

	// Locate task
	grp, err := gState.db.GetGroup(payload.GroupID)
	if err != nil {
		return nil, fmt.Errorf("group not found: %s", err)
	}

	// Save with new name
	grp.Name = payload.Name

	return grp, gState.db.SaveGroup(grp)
}

// HandleUpdateTask updates task value in database
func HandleUpdateTask(msg *message.Message) (interface{}, error) {
	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var payload ReqPayloadUpdateTask
	if err := msg.Into(&payload); err != nil {
		return nil, fmt.Errorf("payload invalid: %s", err)
	}

	if payload.GroupID <= 0 || payload.TaskID <= 0 {
		return nil, fmt.Errorf("group id and task id must be non-zero positive integers")
	}

	if payload.Name == "" {
		return nil, fmt.Errorf("name can't be empty")
	}

	// Locate task
	task, err := gState.db.GetTask(payload.GroupID, payload.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task not found", err)
	}

	// Save with new name
	task.Name = payload.Name
	task.CostCode = payload.CostCode

	return task, gState.db.SaveTask(task)
}

// HandleAddTask adds a task for a group using details from msg data
func HandleAddTask(msg *message.Message) (interface{}, error) {
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

// HandleGetUsage handle request for usage statistics
func HandleGetUsage(msg *message.Message) (interface{}, error) {
	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var payload ReqPayloadGetUsage
	if err := msg.Into(&payload); err != nil {
		return nil, fmt.Errorf("payload invalid: %s", err)
	}

	if payload.GroupID <= 0 {
		return nil, errors.New("invalid group id")
	}

	var start time.Time
	var end time.Time
	var err error

	if start, err = time.Parse("2006-01-02", payload.StartDate); err != nil {
		return nil, fmt.Errorf("start date invalid: %s", err)
	}
	if end, err = time.Parse("2006-01-02", payload.EndDate); err != nil {
		return nil, fmt.Errorf("end date invalid: %s", err)
	}

	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start and end must be defined and not be empty")
	}

	return gState.db.GetUsage(payload.GroupID, start, end)
}
