package message

// Key is a message key, which tells what the message purpose is.
type Key string

// Error keys
const (
	ErrorOperationFailed = Key("failure")
)

// Request keys
const (
	RequestActiveTask     = Key("get.active.task")
	RequestAppVersions    = Key("get.versions")
	RequestOpenDatabase   = Key("open.database")
	RequestAddTask        = Key("add.task")
	RequestAddGroup       = Key("add.group")
	RequestGetHistory     = Key("get.history")
	RequestGetTask        = Key("get.task")
	RequestGroups         = Key("get.groups")
	RequestGroupTasks     = Key("get.group.tasks")
	RequestSetHistory     = Key("set.history")
	RequestStartTask      = Key("start.task")
	RequestStopTask       = Key("stop.task")
	RequestTaskSlices     = Key("get.task.slices")
	RequestUpdateGroup    = Key("update.group")
	RequestUpdateTask     = Key("update.task")
	RequestWindowClose    = Key("window.close")
	RequestWindowMinimize = Key("window.minimize")
)

// Event keys
const (
	EventBackendStatusChanged = Key("backend.status")
)
