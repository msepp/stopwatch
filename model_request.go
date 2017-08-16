package main

// ReqPayloadAddGroup defines data fields available when adding a group
type ReqPayloadAddGroup struct {
	// Name of the new group. Required.
	Name string `json:"name"`
}

// ReqPayloadGetGroupTasks defines data fields available when reading groups tasks
type ReqPayloadGetGroupTasks struct {
	// GroupID of the group to fetch. Required.
	GroupID int `json:"groupID"`
}

// ReqPayloadAddTask defines data fields available when adding a task
type ReqPayloadAddTask struct {
	// CostCode for the task
	CostCode string `json:"costcode"`
	// GroupID the task should be added to. Required.
	GroupID int `json:"groupid"`
	// Name of the new task. Required.
	Name string `json:"name"`
}

// ReqPayloadSetTaskStatus defines data fields required for changing task status
type ReqPayloadSetTaskStatus struct {
	// GroupID of the target task. Required.
	GroupID int `json:"groupid"`
	// TaskID of the target task. Required.
	TaskID int `json:"taskid"`
}
