import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { ResultPhase } from './ResultPhase';
import type {
  AppAction,
  ResultPhase as ResultPhaseState,
} from '../../types';

describe('ResultPhase', () => {
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
      screen.getByText('画像の送信に失敗しました'),
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

    expect(dispatch).toHaveBeenCalledTimes(1);
    expect(dispatch).toHaveBeenCalledWith({
      type: 'CONTINUE',
    } satisfies AppAction);
  });

  it('dispatches CONTINUE when the retake button is clicked', () => {
    const dispatch = vi.fn();
    const state: ResultPhaseState = {
      type: 'result',
      shotNumber: 'SHOT-001',
      status: 'error',
      error: new Error('画像の送信に失敗しました'),
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

    expect(dispatch).toHaveBeenCalledTimes(1);
    expect(dispatch).toHaveBeenCalledWith({
      type: 'CONTINUE',
    } satisfies AppAction);
  });

  it.each([
    {
      status: 'success' as const,
      buttonName: '終了',
    },
    {
      status: 'error' as const,
      buttonName: '終了',
    },
  ])(
    'dispatches EXIT from the $status result',
    ({ status, buttonName }) => {
      const dispatch = vi.fn();
      const state: ResultPhaseState = {
        type: 'result',
        shotNumber: 'SHOT-001',
        status,
        error:
          status === 'error'
            ? new Error('画像の送信に失敗しました')
            : undefined,
      };

      render(
        <ResultPhase
          state={state}
          dispatch={dispatch}
        />,
      );

      fireEvent.click(
        screen.getByRole('button', {
          name: buttonName,
        }),
      );

      expect(dispatch).toHaveBeenCalledTimes(1);
      expect(dispatch).toHaveBeenCalledWith({
        type: 'EXIT',
      } satisfies AppAction);
    },
  );
});
