import * as ActiveTaskActions from '../actions/active-task.actions';

export type Action = ActiveTaskActions.All;

import {Task} from '../../model/task';

export function ActiveTaskReducer(state: Task = null, action: Action): Task {
  switch (action.type) {
    case ActiveTaskActions.SET_ACTIVE_TASK: {
      console.log(action.type, action.payload);
      return action.payload;
    }

    case ActiveTaskActions.STOP_ACTIVE_TASK: {
      if (state !== null && state.id === action.payload.id && state.projectid === action.payload.projectid) {
        return null;
      }

      return state;
    }

    default:
      return state;
  };
}
