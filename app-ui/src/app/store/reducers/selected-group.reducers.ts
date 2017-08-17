import * as SelectedGroupActions from '../actions/selected-group.actions';
import {Group} from '../../model/group';

export type Action = SelectedGroupActions.All;

export function SelectedGroupReducer(state: Group = null, action: Action): Group {
  switch (action.type) {
    case SelectedGroupActions.SET_GROUP: {
      return action.payload;
    }

    default:
      return state;
  };
}
