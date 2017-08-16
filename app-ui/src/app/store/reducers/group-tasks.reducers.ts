import * as GroupTasksActions from '../actions/group-tasks.actions';

export type Action = GroupTasksActions.All;

import {Group} from '../../model/group';

export function GroupTasksReducer(state: Group[] = [], action: Action): Group[] {
  switch (action.type) {
    case GroupTasksActions.SET_TASKS: {
      return action.payload;
    }

    case GroupTasksActions.ADD_TASK: {
      state.push(action.payload);
      return state;
    }

    case GroupTasksActions.UPDATE_TASK: {
      const idx = state.findIndex(v => v.id === action.payload.id);
      if (idx === -1) {
        return state;
      }

      state[idx] = action.payload;
      return state;
    }

    default:
      return state;
  };
}
