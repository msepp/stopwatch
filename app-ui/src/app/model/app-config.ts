import { Group, Task } from './';

export class AppConfig {
  constructor(
    public activeTask?: Task,
    public groups?: Group[]
  ) {}
}
