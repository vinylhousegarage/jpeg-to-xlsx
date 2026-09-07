import {
  fireEvent,
  render,
  screen,
  waitFor,
} from '@testing-library/react';
import {
  beforeEach,
  describe,
  expect,
  it,
  vi,
} from 'vitest';

import type {
  ResultPhase as ResultPhaseState,
} from '../../types';
import { ResultPhase } from './ResultPhase';

const mocks = vi.hoisted(() => ({
  useAuthContext: vi.fn(),
}));

vi.mock('../../auth/AuthContext', () => ({
  useAuthContext: mocks.useAuthContext,
}));

const successState: ResultPhaseState = {
  type: 'result',
  shotNumber: 'SHOT-001',
  status: 'success',
};

const createErrorState =
  (): ResultPhaseState => ({
    type: 'result',
    shotNumber: 'SHOT-001',
    status: 'error',
    error: new Error(
      '画像の送信に失敗しました',
    ),
  });

const createOnFileSelected = () =>
  vi.fn().mockResolvedValue(undefined);

const createImageFile = (): File =>
  new File(
    ['image data'],
    'photo.jpg',
    {
      type: 'image/jpeg',
    },
  );

const getFileInput = (
  container: HTMLElement,
): HTMLInputElement => {
  const input =
    container.querySelector<HTMLInputElement>('input[type="file"]');

  if (!input) {
    throw new Error('file input was not found');
  }

  return input;
};

const renderResultPhase = (
  state: ResultPhaseState,
  onFileSelected =
    createOnFileSelected(),
) => {
  const renderResult = render(
    <ResultPhase
      state={state}
      onFileSelected={onFileSelected}
    />,
  );

  return {
    ...renderResult,
    onFileSelected,
  };
};

const selectImageFile = (
  container: HTMLElement,
): File => {
  const input = getFileInput(container);
  const file = createImageFile();

  fireEvent.change(input, {
    target: {
      files: [file],
    },
  });

  return file;
};

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

  it(
    'renders SuccessDisplay when status is success',
    () => {
      renderResultPhase(successState);

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
    },
  );

  it(
    'renders ErrorDisplay with the error message when status is error',
    () => {
      renderResultPhase(createErrorState());

      expect(
        screen.getByRole('heading', {
          name: '送信失敗',
        }),
      ).toBeInTheDocument();

      expect(screen.getByText('画像の送信に失敗しました')).toBeInTheDocument();

      expect(
        screen.queryByRole('heading', {
          name: '送信完了',
        }),
      ).not.toBeInTheDocument();
    },
  );

  it(
    'passes the selected file from SuccessDisplay to onFileSelected',
    async () => {
      const {
        container,
        onFileSelected,
      } = renderResultPhase(successState);

      const file = selectImageFile(container);

      await waitFor(() => {
        expect(onFileSelected).toHaveBeenCalledWith(file);
      });

      expect(onFileSelected).toHaveBeenCalledTimes(1);
      expect(signOut).not.toHaveBeenCalled();
    },
  );

  it(
    'passes the selected file from ErrorDisplay to onFileSelected',
    async () => {
      const {
        container,
        onFileSelected,
      } = renderResultPhase(createErrorState());

      const file = selectImageFile(container);

      await waitFor(() => {
        expect(onFileSelected).toHaveBeenCalledWith(file);
      });

      expect(onFileSelected).toHaveBeenCalledTimes(1);
      expect(signOut).not.toHaveBeenCalled();
    },
  );

  it(
    'signs out from the success result',
    () => {
      const { onFileSelected } = renderResultPhase(successState);

      fireEvent.click(
        screen.getByRole('button', {
          name: 'ログアウト',
        }),
      );

      expect(signOut).toHaveBeenCalledTimes(1);
      expect(onFileSelected).not.toHaveBeenCalled();
    },
  );

  it(
    'signs out from the error result',
    () => {
      const { onFileSelected } = renderResultPhase(createErrorState());

      fireEvent.click(
        screen.getByRole('button', {
          name: 'ログアウト',
        }),
      );

      expect(signOut).toHaveBeenCalledTimes(1);
      expect(onFileSelected).not.toHaveBeenCalled();
    },
  );
});
