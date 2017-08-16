import * as BackendConnActions from '../actions/backend-conn.actions';

export type Action = BackendConnActions.All;

export function BackendConnReducer(state: boolean = false, action: Action): boolean {
  switch (action.type) {
    case BackendConnActions.SET_BACKEND_CONN: {
      return action.payload;
    }

    default:
      return state;
  };
}
