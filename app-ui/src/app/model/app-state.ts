import { Project } from './project';
import { Task } from './task';
import { AppVersion } from './app-version';

export class AppState {
  constructor(
    public activeTask?: Task,
    public projectTasks?: Task[],
    public projects?: Project[],
    public selectedProject?: Project,
    public selectedTask?: Task,
    public version?: AppVersion
  ) {}
}
