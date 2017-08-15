package message

// Key is a message key, which tells what the message purpose is.
type Key string

// Error keys
const (
	ErrorOperationFailed = Key("failure")
)

// Request keys
const (
	RequestWindowClose    = Key("window.close")
	RequestWindowMinimize = Key("window.minimize")
	RequestAppVersions    = Key("get.versions")
	RequestOpenDatabase   = Key("open.database")
	RequestAddTask        = Key("add.task")
	RequestAddProject     = Key("add.project")
	RequestStartTask      = Key("start.task")
	RequestStopTask       = Key("stop.task")
	RequestProjects       = Key("get.projects")
	RequestProjectTasks   = Key("get.project.tasks")
	RequestTaskSlices     = Key("get.task.slices")
)

// Event keys
const (
	EventBackendStatusChanged = Key("backend.status")
)
