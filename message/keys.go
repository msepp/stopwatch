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
	RequestAddGroup       = Key("add.group")
	RequestActiveTask     = Key("get.active.task")
	RequestGroups         = Key("get.groups")
	RequestGroupTasks     = Key("get.tasks")
	RequestStartTask      = Key("start.task")
	RequestStopTask       = Key("stop.task")
	RequestTaskSlices     = Key("get.task.slices")
)

// Event keys
const (
	EventBackendStatusChanged = Key("backend.status")
)
