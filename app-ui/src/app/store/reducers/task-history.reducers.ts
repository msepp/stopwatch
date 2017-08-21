import * as TaskHistoryActions from '../actions/task-history.actions';

export type Action = TaskHistoryActions.All;

import {Task, taskMatch} from '../../model/task';

export function TaskHistoryReducer(state: Task[] = [], action: Action): Task {
  switch (action.type) {
    case TaskHistoryActions.SET: {
      return action.payload;
    }

    case TaskHistoryActions.ADD: {
      console.log('adding task to history:', state);

      // Send task as first in history.
      state.unshift(action.payload);

      // Return 10 first.
      return state.slice(0, 10);
    }

    case TaskHistoryActions.REMOVE: {
      // If not in history, return state.
      if (!state.find((t: Task) => taskMatch(t, action.payload))) {
        console.log('task not found in history:', action.payload);
        return state;
      }

      // Remove from history by creating new array with only non-matching values
      const history: Task[] = [];
      state.forEach((t: Task) => {
        if (taskMatch(t, action.payload) === false) {
          history.push(t);
        }
      });

      console.log('new history ', history);
      return history;
    }

    case TaskHistoryActions.UPDATE: {
      // See if in history.
      const idx = state.findIndex((t: Task) => taskMatch(t, action.payload));
      if (idx > -1) {
        state[idx] = action.payload;
      }

      return state.slice(0, 10);
    }

    default:
      return state;
  };
}
