import { Injectable } from '@angular/core';
import { Store, Action } from '@ngrx/store';
import { Actions, Effect } from '@ngrx/effects';
import { AppState } from '../../model/app-state';

import * as ProjectsAction from '../actions/projects.actions';

@Injectable()
export class ProjectsEffects {
  constructor() {}

}
