import * as ActiveTaskActions from '../actions/active-task.actions';

export type Action = ActiveTaskActions.All;

import {Task} from '../../model/task';

export function ActiveTaskReducer(state: Task = null, action: Action): Task {
  switch (action.type) {
    case ActiveTaskActions.SET_ACTIVE_TASK: {
      return action.payload;
    }

    case ActiveTaskActions.CLEAR_ACTIVE_TASK: {
      return null;
    }

    default:
      return state;
  };
}
