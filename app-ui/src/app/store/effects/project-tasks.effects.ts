import { Injectable } from '@angular/core';
import { Store, Action } from '@ngrx/store';
import { Actions, Effect } from '@ngrx/effects';
import { AppState } from '../../model/app-state';

import * as ProjectTasksActions from '../actions/project-tasks.actions';

@Injectable()
export class ProjectTasksEffects {
  constructor() {}

}
