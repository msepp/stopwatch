
// Request keys
export const REQUEST_WINDOW_CLOSE    = 'window.close';
export const REQUEST_WINDOW_MINIMIZE = 'window.minimize';
export const REQUEST_APP_VERSIONS    = 'get.versions';
export const REQUEST_ACTIVE_TASK     = 'get.active.task';
export const REQUEST_GROUP_TASKS     = 'get.group.tasks';
export const REQUEST_GROUPS          = 'get.groups';
export const REQUEST_GET_USAGE       = 'get.usage';
export const REQUEST_ADD_GROUP       = 'add.group';
export const REQUEST_UPDATE_GROUP    = 'update.group';
export const REQUEST_ADD_TASK        = 'add.task';
export const REQUEST_GET_TASK        = 'get.task';
export const REQUEST_START_TASK      = 'start.task';
export const REQUEST_STOP_TASK       = 'stop.task';
export const REQUEST_UPDATE_TASK     = 'update.task';
export const REQUEST_OPEN_DATABASE   = 'open.database';
export const REQUEST_GET_TASK_HISTORY = 'get.history';
export const REQUEST_SET_TASK_HISTORY = 'set.history';

// Event keys
export const EVENT_BACKEND_STATUS = 'backend.status';

// Message types
export const MESSAGE_TYPE_ALERT    = 'alert';
export const MESSAGE_TYPE_ERROR    = 'error';
export const MESSAGE_TYPE_EVENT    = 'event';
export const MESSAGE_TYPE_RESPONSE = 'response';

// Backend status keys
export const BACKEND_STATUS_IDLE = 'idle';

// BackendStatus describes backend status message
export class BackendStatus {
  constructor(
    public status?: string,
  ) {}
}

// Message describes an IPC message sent between render and backend process.
// Can be either a response, error (for a request) or an event.
export class Message {
  constructor(
    public id: string,
    public key: string,
    public type: string,
    public data: any
  ) {}
}
