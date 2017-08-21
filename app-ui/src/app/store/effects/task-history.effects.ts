import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { Actions, Effect } from '@ngrx/effects';

import * as TaskHistoryActions from '../actions/task-history.actions';

@Injectable()
export class TaskHistoryEffects {
  // Listen for the 'PUSH' action
  @Effect() push$ = this.actions$.ofType(TaskHistoryActions.PUSH)
    // Push does a REMOVE and ADD to keep only one copy in history at a time.
    .flatMap((a: TaskHistoryActions.Update) => [
      new TaskHistoryActions.Remove(a.payload),
      new TaskHistoryActions.Add(a.payload)
    ]);

  constructor(
    private actions$: Actions
  ) {}
}
