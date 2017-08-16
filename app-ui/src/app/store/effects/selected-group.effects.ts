import { Injectable } from '@angular/core';
import { Store, Action } from '@ngrx/store';
import { Actions, Effect } from '@ngrx/effects';
import { AppState } from '../../model/app-state';

import * as SelectedGroupActions from '../actions/selected-group.actions';

@Injectable()
export class SelectedGroupEffects {
  constructor() {}

}
