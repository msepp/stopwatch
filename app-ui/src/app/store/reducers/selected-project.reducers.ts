import * as SelectedProjectActions from '../actions/selected-project.actions';

export type Action = SelectedProjectActions.All;

import {Project} from '../../model/project';

export function SelectedProjectReducer(state: Project = null, action: Action): Project {
  switch (action.type) {
    case SelectedProjectActions.SET_PROJECT: {
      return action.payload;
    }

    default:
      return state;
  };
}
