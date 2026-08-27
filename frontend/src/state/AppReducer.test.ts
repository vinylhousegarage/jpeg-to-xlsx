import { describe, expect, it } from 'vitest';
import { appReducer } from './AppReducer';
import { AppState, AppAction } from '../types';

describe('appReducer', () => {
  const initialState: AppState = {
    isSlackLinked: false,
    phase: {
      type: 'input',
    },
  };

  it('should return the initial state when action is unknown', () => {
    const action = {
      type: 'INVALID_TYPE',
    } as unknown as AppAction;

    expect(appReducer(initialState, action)).toEqual(initialState);
  });

  it('should handle SET_SLACK_LINKED', () => {
    const action: AppAction = {
      type: 'SET_SLACK_LINKED',
      isSlackLinked: true,
    };

    const state = appReducer(initialState, action);

    expect(state).toEqual({
      isSlackLinked: true,
      phase: {
        type: 'input',
      },
    });
  });

  it('should handle SET_PREVIEW', () => {
    const file = new File([''], 'test.png');

    const action: AppAction = {
      type: 'SET_PREVIEW',
      file,
      shotNumber: 'SHOT-001',
    };

    const state = appReducer(initialState, action);

    expect(state).toEqual({
      isSlackLinked: false,
      phase: {
        type: 'preview',
        file,
        shotNumber: 'SHOT-001',
      },
    });
  });

  it('should preserve Slack linked state when handling SET_PREVIEW', () => {
    const linkedState: AppState = {
      isSlackLinked: true,
      phase: {
        type: 'input',
      },
    };

    const file = new File([''], 'test.png');

    const action: AppAction = {
      type: 'SET_PREVIEW',
      file,
      shotNumber: 'SHOT-001',
    };

    const state = appReducer(linkedState, action);

    expect(state.isSlackLinked).toBe(true);
    expect(state.phase.type).toBe('preview');
  });

  it('should handle UPLOAD_COMPLETE', () => {
    const uploadState: AppState = {
      isSlackLinked: true,
      phase: {
        type: 'upload',
        shotNumber: 'SHOT-001',
      },
    };

    const action: AppAction = {
      type: 'UPLOAD_COMPLETE',
      status: 'success',
      error: undefined,
    };

    const state = appReducer(uploadState, action);

    expect(state).toEqual({
      isSlackLinked: true,
      phase: {
        type: 'result',
        shotNumber: 'SHOT-001',
        status: 'success',
        error: undefined,
      },
    });
  });

  it('should handle EXIT and preserve Slack linked state', () => {
    const previewState: AppState = {
      isSlackLinked: true,
      phase: {
        type: 'preview',
        shotNumber: 'SHOT-001',
        file: new File([''], 'test.png'),
      },
    };

    const action: AppAction = {
      type: 'EXIT',
    };

    const state = appReducer(previewState, action);

    expect(state).toEqual({
      isSlackLinked: true,
      phase: {
        type: 'input',
      },
    });
  });
});
