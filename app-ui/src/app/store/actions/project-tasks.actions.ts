import {Injectable} from '@angular/core';
import {Action, Store} from '@ngrx/store';
import {Actions, Effect} from '@ngrx/effects';
import {Project} from '../../model/project';

export const SET_PROJECT_TASKS = '[project tasks] SET';
export class Set implements Action {
  readonly type = SET_PROJECT_TASKS;
  constructor(public payload: Project[]) {}
}

export const ADD_PROJECT_TASK = '[project tasks] ADD';
export class Add implements Action {
  readonly type = ADD_PROJECT_TASK;
  constructor(public payload: Project) {}
}

export const UPDATE_PROJECT_TASK = '[project tasks] UPDATE';
export class Update implements Action {
  readonly type = UPDATE_PROJECT_TASK;
  constructor(public payload: Project) {}
}

export type All = Set | Add | Update;
