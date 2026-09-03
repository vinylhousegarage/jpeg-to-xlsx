import {
  fireEvent,
  render,
  screen,
} from '@testing-library/react';
import {
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from 'vitest';

import type {
  AppAction,
  ResultPhase as ResultPhaseState,
} from '../../types';
import { ResultPhase } from './ResultPhase';

const mocks = vi.hoisted(() => ({
  useAuthContext: vi.fn(),
}));

vi.mock('../../auth/AuthContext', () => ({
  useAuthContext: mocks.useAuthContext,
}));

describe('ResultPhase', () => {
  const signIn = vi.fn();
  const signOut = vi.fn();

  beforeEach(() => {
    mocks.useAuthContext.mockReset();
    signIn.mockReset();
    signOut.mockReset();

    signIn.mockResolvedValue(undefined);
    signOut.mockResolvedValue(undefined);

    mocks.useAuthContext.mockReturnValue({
      status: 'authenticated',
      error: null,
      signIn,
      signOut,
    });
  });

  it('renders SuccessDisplay when status is success', () => {
    const dispatch = vi.fn();
    const state: ResultPhaseState = {
      type: 'result',
      shotNumber: 'SHOT-001',
      status: 'success',
    };

    render(
      <ResultPhase
        state={state}
        dispatch={dispatch}
      />,
    );

    expect(
      screen.getByRole('heading', {
        name: '送信完了',
      }),
    ).toBeInTheDocument();

    expect(
      screen.queryByRole('heading', {
        name: '送信失敗',
      }),
    ).not.toBeInTheDocument();
  });

  it('renders ErrorDisplay with the error message when status is error', () => {
    const dispatch = vi.fn();
    const error = new Error('画像の送信に失敗しました');
    const state: ResultPhaseState = {
      type: 'result',
      shotNumber: 'SHOT-001',
      status: 'error',
      error,
    };

    render(
      <ResultPhase
        state={state}
        dispatch={dispatch}
      />,
    );

    expect(
      screen.getByRole('heading', {
        name: '送信失敗',
      }),
    ).toBeInTheDocument();

    expect(
      screen.getByText(
        '画像の送信に失敗しました',
      ),
    ).toBeInTheDocument();

    expect(
      screen.queryByRole('heading', {
        name: '送信完了',
      }),
    ).not.toBeInTheDocument();
  });

  it('dispatches CONTINUE when the continue button is clicked', () => {
    const dispatch = vi.fn();
    const state: ResultPhaseState = {
      type: 'result',
      shotNumber: 'SHOT-001',
      status: 'success',
    };

    render(
      <ResultPhase
        state={state}
        dispatch={dispatch}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: 'つづけて撮影',
      }),
    );

    expect(dispatch).toHaveBeenCalledOnce();

    expect(dispatch).toHaveBeenCalledWith({
      type: 'CONTINUE',
    } satisfies AppAction);

    expect(signOut).not.toHaveBeenCalled();
  });

  it('dispatches CONTINUE when the retake button is clicked', () => {
    const dispatch = vi.fn();
    const state: ResultPhaseState = {
      type: 'result',
      shotNumber: 'SHOT-001',
      status: 'error',
      error: new Error(
        '画像の送信に失敗しました',
      ),
    };

    render(
      <ResultPhase
        state={state}
        dispatch={dispatch}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: '撮り直し',
      }),
    );

    expect(dispatch).toHaveBeenCalledOnce();

    expect(dispatch).toHaveBeenCalledWith({
      type: 'CONTINUE',
    } satisfies AppAction);
  });

  it('signs out when the logout button is clicked', () => {
    const dispatch = vi.fn();
    const state: ResultPhaseState = {
      type: 'result',
      shotNumber: 'SHOT-001',
      status: 'success',
    };

    render(
      <ResultPhase
        state={state}
        dispatch={dispatch}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: 'ログアウト',
      }),
    );

    expect(signOut).toHaveBeenCalledOnce();
    expect(dispatch).not.toHaveBeenCalled();
  });

  it('dispatches EXIT from the error result', () => {
    const dispatch = vi.fn();
    const state: ResultPhaseState = {
      type: 'result',
      shotNumber: 'SHOT-001',
      status: 'error',
      error: new Error(
        '画像の送信に失敗しました',
      ),
    };

    render(
      <ResultPhase
        state={state}
        dispatch={dispatch}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: '終了',
      }),
    );

    expect(dispatch).toHaveBeenCalledOnce();

    expect(dispatch).toHaveBeenCalledWith({
      type: 'EXIT',
    } satisfies AppAction);

    expect(signOut).not.toHaveBeenCalled();
  });
});
