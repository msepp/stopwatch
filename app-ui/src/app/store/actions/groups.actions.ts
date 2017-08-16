import {Injectable} from '@angular/core';
import {Action, Store} from '@ngrx/store';
import {Actions, Effect} from '@ngrx/effects';
import {Group} from '../../model/group';

export const SET_GROUPS = '[groups] SET';
export class Set implements Action {
  readonly type = SET_GROUPS;
  constructor(public payload: Group[]) {}
}

export const ADD_GROUP = '[groups] ADD';
export class Add implements Action {
  readonly type = ADD_GROUP;
  constructor(public payload: Group) {}
}

export type All = Set | Add;
