import {Injectable} from '@angular/core';
import {Action} from '@ngrx/store';

export const SET_BACKEND_CONN = '[backend conn] SET';
export class Set implements Action {
  readonly type = SET_BACKEND_CONN;
  constructor(public payload: boolean) {}
}

export type All = Set;
