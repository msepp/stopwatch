import {Injectable} from '@angular/core';
import {Action, Store} from '@ngrx/store';
import {Actions, Effect} from '@ngrx/effects';
import {Task} from '../../model/task';

export const SET_ACTIVE_TASK = '[active task] SET';
export class Set implements Action {
  readonly type = SET_ACTIVE_TASK;
  constructor(public payload: Task) {}
}

export const STOP_ACTIVE_TASK = '[active task] STOP';
export class Stop implements Action {
  readonly type = STOP_ACTIVE_TASK;
  constructor(public payload: Task) {}
}


export type All = Stop | Set;
