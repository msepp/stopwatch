import { Injectable } from '@angular/core';
import { Store, Action } from '@ngrx/store';
import { Actions, Effect } from '@ngrx/effects';
import { AppState } from '../../model/app-state';

import * as GroupsAction from '../actions/groups.actions';

@Injectable()
export class GroupsEffects {
  constructor() {}

}
