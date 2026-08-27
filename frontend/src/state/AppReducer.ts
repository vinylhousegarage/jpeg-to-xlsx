import { AppState, AppAction } from '../types';

export const appReducer = (
  state: AppState,
  action: AppAction,
): AppState => {
  switch (action.type) {
    case 'SET_SLACK_LINKED':
      return {
        ...state,
        isSlackLinked: action.isSlackLinked,
      };

    case 'SET_PREVIEW':
      return {
        ...state,
        phase: {
          type: 'preview',
          file: action.file,
          shotNumber: action.shotNumber,
        },
      };

    case 'RETAKE':
    case 'CONTINUE':
    case 'EXIT':
      return {
        ...state,
        phase: {
          type: 'input',
        },
      };

    case 'SEND':
      if (state.phase.type !== 'preview') {
        return state;
      }

      return {
        ...state,
        phase: {
          type: 'upload',
          shotNumber: state.phase.shotNumber,
        },
      };

    case 'START_UPLOAD':
      return {
        ...state,
        phase: {
          type: 'upload',
          shotNumber:
            'shotNumber' in state.phase
              ? state.phase.shotNumber
              : '',
        },
      };

    case 'UPLOAD_COMPLETE':
      return {
        ...state,
        phase: {
          type: 'result',
          shotNumber:
            'shotNumber' in state.phase
              ? state.phase.shotNumber
              : '',
          status: action.status,
          error: action.error,
        },
      };

    default:
      return state;
  }
};
