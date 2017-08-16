import { Injectable } from '@angular/core';
import { Store, Action } from '@ngrx/store';
import { Actions, Effect } from '@ngrx/effects';
import { AppState } from '../../model/app-state';

import * as SelectedProjectActions from '../actions/selected-project.actions';

@Injectable()
export class SelectedProjectEffects {
  constructor() {}

}
