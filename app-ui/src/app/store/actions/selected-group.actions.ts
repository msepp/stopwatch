import { Injectable } from '@angular/core';
import { Action } from '@ngrx/store';
import { Group } from '../../model/group';

export const SET_GROUP = '[group] SET';
export class Set implements Action {
  readonly type = SET_GROUP;
  constructor(public payload: Group) {}
}

export type All = Set;
