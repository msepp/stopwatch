import {Injectable} from '@angular/core';
import {Action} from '@ngrx/store';
import { Task } from '../../model/task';

export const SET_TASK = '[task] SET';
export class Set implements Action {
  readonly type = SET_TASK;
  constructor(public payload: Task) {}
}

export type All = Set;
