import * as VersionActions from '../actions/version.actions';

export type Action = VersionActions.All;

import {AppVersion} from '../../model/app-version';

export function VersionReducer(state: AppVersion = null, action: Action): AppVersion {
  switch (action.type) {
    case VersionActions.SET_VERSION: {
      return action.payload;
    }

    default:
      return state;
  };
}
