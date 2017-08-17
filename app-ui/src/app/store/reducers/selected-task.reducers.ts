import * as SelectedTaskActions from '../actions/selected-task.actions';
import {Task} from '../../model/task';

export type Action = SelectedTaskActions.All;

export function SelectedTaskReducer(state: Task = null, action: Action): Task {
  switch (action.type) {
    case SelectedTaskActions.SET_TASK: {
      return action.payload;
    }

    default:
        return state;
  };
}
