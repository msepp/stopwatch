import {Injectable} from '@angular/core';
import {Action, Store} from '@ngrx/store';
import {Actions, Effect} from '@ngrx/effects';
import {Task} from '../../model/task';

export const SET_ACTIVE_TASK = '[active task] SET';
export class Set implements Action {
  readonly type = SET_ACTIVE_TASK;
  constructor(public payload: Task) {}
}

export const CLEAR_ACTIVE_TASK = '[active task] CLEAR';
export class Clear implements Action {
  readonly type = CLEAR_ACTIVE_TASK;
  constructor() {}
}

export const UPDATE_ACTIVE_TASK = '[active task] UPDATE';
export class Update implements Action {
  readonly type = UPDATE_ACTIVE_TASK;
  constructor(public payload: Task) {}
}


export type All = Clear | Set | Update;
