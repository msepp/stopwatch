import { Injectable } from '@angular/core';
import { Action } from '@ngrx/store';
import { AppVersion } from '../../model/app-version';

export const SET_VERSION = '[version] SET';
export class Set implements Action {
  readonly type = SET_VERSION;
  constructor(public payload: AppVersion) {}
}

export type All = Set;
