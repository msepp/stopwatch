import * as ProjectsActions from '../actions/projects.actions';

export type Action = ProjectsActions.All;

import {Project} from '../../model/project';

export function ProjectsReducer(state: Project[] = [], action: Action): Project[] {
  switch (action.type) {
    case ProjectsActions.SET_PROJECTS: {
      return action.payload;
    }

    case ProjectsActions.ADD_PROJECT: {
      state.push(action.payload);
      return state;
    }

    default:
      return state;
  };
}
