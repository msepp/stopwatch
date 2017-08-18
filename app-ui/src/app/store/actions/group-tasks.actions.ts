import {Injectable} from '@angular/core';
import {Action} from '@ngrx/store';
import {Task} from '../../model/task';

export const SET_TASKS = '[group tasks] SET';
export class Set implements Action {
  readonly type = SET_TASKS;
  constructor(public payload: Task[]) {}
}

export const ADD_TASK = '[group tasks] ADD';
export class Add implements Action {
  readonly type = ADD_TASK;
  constructor(public payload: Task) {}
}

export const UPDATE_TASK = '[group tasks] UPDATE';
export class Update implements Action {
  readonly type = UPDATE_TASK;
  constructor(public payload: Task) {}
}

export const TASK_UPDATED = '[group tasks] UPDATED';
export class Updated implements Action {
  readonly type = TASK_UPDATED;
  constructor(public payload: Task) {}
}

export const START_TASK = '[group tasks] START';
export class Start implements Action {
  readonly type = START_TASK;
  constructor(public payload: Task) {}
}

export const STOP_TASK = '[group tasks] STOP';
export class Stop implements Action {
  readonly type = STOP_TASK;
  constructor(public payload: Task) {}
}

export type All = Set | Add | Update | Updated | Start | Stop;
