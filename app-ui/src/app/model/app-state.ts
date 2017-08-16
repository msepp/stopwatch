import { AppVersion, Group, Task } from './';

export class AppState {
  constructor(
    public activeTask?: Task,
    public groupTasks?: Task[],
    public groups?: Group[],
    public selectedGroup?: Group,
    public selectedTask?: Task,
    public version?: AppVersion
  ) {}
}
