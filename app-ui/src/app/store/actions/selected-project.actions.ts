import { Injectable } from '@angular/core';
import { Action } from '@ngrx/store';
import { Project } from '../../model/project';

export const SET_PROJECT = '[project] SET';
export class Set implements Action {
  readonly type = SET_PROJECT;
  constructor(public payload: Project) {}
}

export type All = Set;
