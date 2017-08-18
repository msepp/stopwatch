import * as TaskHistoryActions from '../actions/task-history.actions';

export type Action = TaskHistoryActions.All;

import {Task} from '../../model/task';

export function TaskHistoryReducer(state: Task[] = [], action: Action): Task {
  switch (action.type) {
    case TaskHistoryActions.PUSH: {
      // Send task as first in history.
      state.unshift(action.payload);

      // Return 10 first.
      return state.slice(0, 10);
    }

    case TaskHistoryActions.REMOVE: {
      // See if in history.
      const idx = state.findIndex((t: Task) =>
        (t.id === action.payload.id && t.groupid === action.payload.groupid)
      );

      if (idx > -1) {
        state.splice(idx, 1);
      }

      return state;
    }

    case TaskHistoryActions.UPDATE: {
      // See if in history.
      const idx = state.findIndex((t: Task) =>
        (t.id === action.payload.id && t.groupid === action.payload.groupid)
      );

      if (idx > -1) {
        state[idx] = action.payload;
      }

      return state;
    }

    default:
      return state;
  };
}
