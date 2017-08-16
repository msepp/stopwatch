import {Injectable} from '@angular/core';
import {Action, Store} from '@ngrx/store';
import {Actions, Effect} from '@ngrx/effects';
import {Project} from '../../model/project';

export const SET_PROJECTS = '[projects] SET';
export class Set implements Action {
  readonly type = SET_PROJECTS;
  constructor(public payload: Project[]) {}
}

export const ADD_PROJECT = '[projects] ADD';
export class Add implements Action {
  readonly type = ADD_PROJECT;
  constructor(public payload: Project) {}
}

export type All = Set | Add;
