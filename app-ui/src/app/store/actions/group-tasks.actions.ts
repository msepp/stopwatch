import {Injectable} from '@angular/core';
import {Action, Store} from '@ngrx/store';
import {Actions, Effect} from '@ngrx/effects';
import {Task} from '../../model/task';

export const SET_GROUP_TASKS = '[group tasks] SET';
export class Set implements Action {
  readonly type = SET_GROUP_TASKS;
  constructor(public payload: Task[]) {}
}

export const ADD_GROUP_TASK = '[group tasks] ADD';
export class Add implements Action {
  readonly type = ADD_GROUP_TASK;
  constructor(public payload: Task) {}
}

export const UPDATE_GROUP_TASK = '[group tasks] UPDATE';
export class Update implements Action {
  readonly type = UPDATE_GROUP_TASK;
  constructor(public payload: Task) {}
}

export type All = Set | Add | Update;
