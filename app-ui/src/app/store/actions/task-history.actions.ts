import {Injectable} from '@angular/core';
import {Action, Store} from '@ngrx/store';
import {Actions, Effect} from '@ngrx/effects';
import {Task} from '../../model/task';

export const ADD = '[task history] ADD';
export class Add implements Action {
  readonly type = ADD;
  constructor(public payload: Task) {}
}

export const PUSH = '[task history] PUSH';
export class Push implements Action {
  readonly type = PUSH;
  constructor(public payload: Task) {}
}

export const REMOVE = '[task history] REMOVE';
export class Remove implements Action {
  readonly type = REMOVE;
  constructor(public payload: Task) {}
}

export const UPDATE = '[task history] UPDATE';
export class Update implements Action {
  readonly type = UPDATE;
  constructor(public payload: Task) {}
}

export const SET = '[task history] SET';
export class Set implements Action {
  readonly type = SET;
  constructor(public payload: Task[]) {}
}

export type All = Add|Remove|Push|Update|Set;
