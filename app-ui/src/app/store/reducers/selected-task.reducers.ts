import {Action} from '@ngrx/store';

import {Task} from '../../model/task';
import {SelectedTaskActions} from '../actions/selected-task.actions';

export function SelectedTaskReducer(state: Task = null, action: Action) {
  switch (action.type) {
    default:
        return state;
  };
}
