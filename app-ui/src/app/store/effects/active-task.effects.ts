import { Injectable } from '@angular/core';
import { Store, Action } from '@ngrx/store';
import { Actions, Effect } from '@ngrx/effects';
import { AppState } from '../../model/app-state';

import * as ActiveTaskActions from '../actions/active-task.actions';

@Injectable()
export class ActiveTaskEffects {
  constructor() {}

}
