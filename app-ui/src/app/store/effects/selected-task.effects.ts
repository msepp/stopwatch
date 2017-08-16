import { Injectable } from '@angular/core';
import { Store, Action } from '@ngrx/store';
import { Actions, Effect } from '@ngrx/effects';
import { AppState } from '../../model/app-state';

import { SelectedTaskActions } from '../actions/selected-task.actions';

@Injectable()
export class SelectedTaskEffects {
  constructor() {}

}
