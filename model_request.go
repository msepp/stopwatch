package main

// ReqPayloadAddGroup defines data fields available when adding a group
type ReqPayloadAddGroup struct {
	// Name of the new group. Required.
	Name string `json:"name" mapstructure:"name"`
}

// ReqPayloadUpdateGroup defines data fields required to update a group
type ReqPayloadUpdateGroup struct {
	// GroupID is the groups ID. Required.
	GroupID int `json:"id" mapstructure:"id"`
	// Name of the new group. Required.
	Name string `json:"name" mapstructure:"name"`
}

// ReqPayloadGetGroupTasks defines data fields available when reading groups tasks
type ReqPayloadGetGroupTasks struct {
	// GroupID of the group to fetch. Required.
	GroupID int `json:"id" mapstructure:"id"`
}

// ReqPayloadAddTask defines data fields available when adding a task
type ReqPayloadAddTask struct {
	// CostCode for the task
	CostCode string `json:"costcode" mapstructure:"costcode"`
	// GroupID the task should be added to. Required.
	GroupID int `json:"groupid" mapstructure:"groupid"`
	// Name of the new task. Required.
	Name string `json:"name" mapstructure:"name"`
}

// ReqPayloadSetTaskStatus defines data fields required for changing task status
type ReqPayloadSetTaskStatus struct {
	// GroupID of the target task. Required.
	GroupID int `json:"groupid" mapstructure:"groupid"`
	// TaskID of the target task. Required.
	TaskID int `json:"taskid" mapstructure:"id"`
}

// ReqPayloadUpdateTask defines data fields required for updating task data
type ReqPayloadUpdateTask struct {
	// GroupID of the target task. Required.
	GroupID int `json:"groupid" mapstructure:"groupid"`
	// TaskID of the target task. Required.
	TaskID int `json:"taskid" mapstructure:"id"`
	// Name of the new task. Required.
	Name string `json:"name" mapstructure:"name"`
	// CostCode for the task.
	CostCode string `json:"costcode" mapstructure:"costcode"`
}
