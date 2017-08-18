import { AppVersion, Group, Task } from './';

export class AppState {
  constructor(
    public backendConn: boolean,
    public activeTask?: Task,
    public groupTasks?: Task[],
    public taskHistory?: Task[],
    public groups?: Group[],
    public selectedGroup?: Group,
    public selectedTask?: Task,
    public version?: AppVersion
  ) {}
}
