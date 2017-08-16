import * as ProjectTasksActions from '../actions/project-tasks.actions';

export type Action = ProjectTasksActions.All;

import {Project} from '../../model/project';

export function ProjectTasksReducer(state: Project[] = [], action: Action): Project[] {
  switch (action.type) {
    case ProjectTasksActions.SET_PROJECT_TASKS: {
      return action.payload;
    }

    case ProjectTasksActions.ADD_PROJECT_TASK: {
      state.push(action.payload);
      return state;
    }

    case ProjectTasksActions.UPDATE_PROJECT_TASK: {
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
