import * as GroupsActions from '../actions/groups.actions';

export type Action = GroupsActions.All;

import {Group} from '../../model/group';

export function GroupsReducer(state: Group[] = [], action: Action): Group[] {
  switch (action.type) {
    case GroupsActions.SET_GROUPS: {
      return action.payload;
    }

    case GroupsActions.ADD_GROUP: {
      state.push(action.payload);
      return state;
    }

    default:
      return state;
  };
}
