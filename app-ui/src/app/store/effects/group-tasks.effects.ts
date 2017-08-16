import { Injectable } from '@angular/core';
import { Store, Action } from '@ngrx/store';
import { Actions, Effect } from '@ngrx/effects';
import { AppState } from '../../model/app-state';

import * as GroupTasksActions from '../actions/group-tasks.actions';

@Injectable()
export class GroupTasksEffects {
  constructor() {}

}
