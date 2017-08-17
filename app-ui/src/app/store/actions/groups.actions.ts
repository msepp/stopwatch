import {Injectable} from '@angular/core';
import {Action} from '@ngrx/store';
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

export const UPDATE_GROUP = '[groups] UPDATE';
export class Update implements Action {
  readonly type = UPDATE_GROUP;
  constructor(public payload: Group) {}
}

export const LOAD_GROUPS = '[groups] LOAD';
export class Load implements Action {
  readonly type = LOAD_GROUPS;
  constructor(public payload?: any) {}
}

export type All = Set | Add | Load | Update;
