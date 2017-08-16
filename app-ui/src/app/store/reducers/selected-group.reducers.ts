import * as SelectedGroupActions from '../actions/selected-group.actions';

export type Action = SelectedGroupActions.All;

import {Group} from '../../model/group';

export function SelectedGroupReducer(state: Group = null, action: Action): Group {
  switch (action.type) {
    case SelectedGroupActions.SET_GROUP: {
      return action.payload;
    }

    default:
      return state;
  };
}
