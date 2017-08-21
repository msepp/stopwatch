import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { Actions, Effect } from '@ngrx/effects';

import * as GroupTasksActions from '../actions/group-tasks.actions';
import * as ActiveTaskActions from '../actions/active-task.actions';
import * as TaskHistoryActions from '../actions/task-history.actions';

@Injectable()
export class GroupTasksEffects {
  // Listen for the 'UPDATE_TASK' action
  @Effect() update$ = this.actions$.ofType(GroupTasksActions.UPDATE_TASK)
    // Map the payload into JSON to use as the request body
    .flatMap((a: GroupTasksActions.Update) => [
      new GroupTasksActions.Updated(a.payload),
      new ActiveTaskActions.Update(a.payload),
      new TaskHistoryActions.Update(a.payload)
    ]);

  constructor(
    private actions$: Actions
  ) {}
}
