package main

import (
	"errors"
	"fmt"
	"log"
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

	case message.RequestProjects:
		return HandleGetProjects(msg)

	case message.RequestProjectTasks:
		return HandleGetProjectTasks(msg)

	case message.RequestAddProject:
		return HandleAddProject(msg)

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

// HandleGetProjects returns current list of projects
func HandleGetProjects(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	// Retrieve projects
	plist, err := gState.db.ReadProjects()
	if err != nil {
		return nil, fmt.Errorf("Unable to read projects: %s", err)
	}

	// Get active task
	at, err := gState.db.GetActiveTask()
	if err != nil {
		return nil, fmt.Errorf("Unable to read active task: %s", err)
	}

	return map[string]interface{}{
		"projects": plist, "activeTask": at,
	}, nil
}

// HandleGetProjectTasks returns list of tasks for a project
func HandleGetProjectTasks(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var projectID int
	if mdata, ok := msg.DataMap(); ok {
		projectID, _ = mdata.GetInt("projectid")
	}

	// Retrieve tasks
	tlist, err := gState.db.ReadTasks(projectID)
	if err != nil {
		return nil, fmt.Errorf("Unable to read tasks: %s", err)
	}

	return tlist, nil
}

// HandleAddProject adds the project detailed in msg data
func HandleAddProject(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var name string
	if mdata, ok := msg.DataMap(); ok {
		name, _ = mdata.GetString("name")
	}

	if name == "" {
		return nil, fmt.Errorf("name must be given and not empty")
	}

	p, err := gState.db.AddProject(name)
	if err != nil {
		return nil, fmt.Errorf("failed to add project: %s", err)
	}

	return p, nil
}

// HandleAddTask adds a task for a project using details from msg data
func HandleAddTask(msg *message.Message) (interface{}, error) {
	if gState.app.SetOpRunning(true) == false {
		return nil, errors.New("backend busy")
	}
	defer gState.app.SetOpRunning(false)

	if gState.db == nil {
		return nil, fmt.Errorf("no database")
	}

	var name string
	var costcode string
	var projectID int

	if mdata, ok := msg.DataMap(); ok {
		name, _ = mdata.GetString("name")
		projectID, _ = mdata.GetInt("projectid")
		costcode, _ = mdata.GetString("costcode")
	}

	if name == "" {
		return nil, fmt.Errorf("name must be given and not empty")
	}
	if projectID <= 0 {
		return nil, fmt.Errorf("projectid must be a non-zero positive integer")
	}

	t, err := gState.db.AddTask(projectID, name, costcode)
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

	var projectID int
	var taskID int

	if mdata, ok := msg.DataMap(); ok {
		projectID, _ = mdata.GetInt("projectid")
		taskID, _ = mdata.GetInt("id")
	}

	log.Printf("pid: %d, tid: %d", projectID, taskID)

	// Locate task
	task, err := gState.db.GetTask(projectID, taskID)
	if err != nil {
		return nil, fmt.Errorf("failure starting task: %s", err)
	}

	activeTask, err := gState.db.GetActiveTask()
	if err != nil {
		return nil, fmt.Errorf("Unable test active task: %s", err)
	}

	// Stop active task if not the same as current task
	if activeTask != nil {
		if activeTask.ID == task.ID && activeTask.ProjectID == task.ProjectID {
			if activeTask.Running != nil {
				return activeTask, nil
			}
		} else {
			if _, err = gState.db.StopTask(activeTask.ProjectID, activeTask.ID); err != nil {
				return nil, fmt.Errorf("unable to stop current task: %s", err)
			}

			if err = gState.db.SetActiveTask(task.ProjectID, task.ID); err != nil {
				return nil, fmt.Errorf("unable to set active task: %s", err)
			}
		}
	}

	// Mark running
	return gState.db.StartTask(task.ProjectID, task.ID)
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

	var projectID int
	var taskID int

	if mdata, ok := msg.DataMap(); ok {
		projectID, _ = mdata.GetInt("projectid")
		taskID, _ = mdata.GetInt("id")
	}

	// Locate task
	task, err := gState.db.GetTask(projectID, taskID)
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
		if activeTask.ID == task.ID && activeTask.ProjectID == task.ProjectID {
			if err := gState.db.SetActiveTask(0, 0); err != nil {
				return nil, fmt.Errorf("unable to clear active task: %s")
			}
		}
	}

	// Mark running
	return gState.db.StopTask(task.ProjectID, task.ID)
}
