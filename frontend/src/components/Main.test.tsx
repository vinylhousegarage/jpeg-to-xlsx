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
  AppState,
} from '../types';
import { Main } from './Main';

const {
  mockDispatch,
  mockProcessImage,
  mockSend,
  mockSignOut,
  mockUseAppState,
} = vi.hoisted(() => ({
  mockDispatch: vi.fn(),
  mockProcessImage: vi.fn(),
  mockSend: vi.fn(),
  mockSignOut: vi.fn(),
  mockUseAppState: vi.fn(),
}));

vi.mock('../auth/AuthContext', () => ({
  useAuthContext: () => ({
    signOut: mockSignOut,
  }),
}));

vi.mock('../state/useContext', () => ({
  useAppState: () => mockUseAppState(),
}));

vi.mock('../hooks/useImageProcessor', () => ({
  useImageProcessor: () => ({
    isCompressing: false,
    processImage: mockProcessImage,
  }),
}));

vi.mock('../hooks/usePresignUpload', () => ({
  usePresignUpload: () => ({
    send: mockSend,
  }),
}));

vi.mock('./InputPhase', () => ({
  InputPhase: () => (
    <div data-testid="input-phase">
      InputPhase
    </div>
  ),
}));

vi.mock('./PreviewPhase', () => ({
  PreviewPhase: ({
    onCancel,
  }: {
    onCancel: () => void;
  }) => (
    <div data-testid="preview-phase">
      <button
        type="button"
        onClick={onCancel}
      >
        中止
      </button>
    </div>
  ),
}));

vi.mock('./ResultPhase', () => ({
  ResultPhase: () => (
    <div data-testid="result-phase">
      ResultPhase
    </div>
  ),
}));

vi.mock('./CanceledDisplay', () => ({
  CanceledDisplay: ({
    onContinueFileSelected,
    onLogout,
  }: {
    onContinueFileSelected: (
      file: File,
    ) => Promise<void>;
    onLogout: () => void;
  }) => (
    <div data-testid="canceled-display">
      <p>処理を中止しました</p>

      <button
        type="button"
        onClick={() => {
          void onContinueFileSelected(
            new File(
              ['test-image'],
              'continue.jpg',
              {
                type: 'image/jpeg',
              },
            ),
          );
        }}
      >
        つづけて撮影
      </button>

      <button
        type="button"
        onClick={onLogout}
      >
        ログアウト
      </button>
    </div>
  ),
}));

vi.mock('../common/Spinner', () => ({
  Spinner: () => (
    <div data-testid="spinner">
      Spinner
    </div>
  ),
}));

describe('Main', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockSignOut.mockResolvedValue(undefined);

    window.history.replaceState({}, '', '/');
  });

  it('dispatches CANCEL when the cancel button is clicked in the preview phase', () => {
    const state: AppState = {
      isSlackLinked: true,
      phase: {
        type: 'preview',
        file: new Blob(
          ['test-image'],
          {
            type: 'image/jpeg',
          },
        ),
        shotNumber: 'SHOT-001',
      },
    };

    mockUseAppState.mockReturnValue({
      state,
      dispatch: mockDispatch,
    });

    render(<Main />);

    expect(screen.getByTestId('preview-phase')).toBeInTheDocument();

    fireEvent.click(
      screen.getByRole('button', {
        name: '中止',
      }),
    );

    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'CANCEL',
    });
  });

  it('renders CanceledDisplay in the canceled phase', () => {
    const state: AppState = {
      isSlackLinked: true,
      phase: {
        type: 'canceled',
      },
    };

    mockUseAppState.mockReturnValue({
      state,
      dispatch: mockDispatch,
    });

    render(<Main />);

    expect(screen.getByTestId('canceled-display')).toHaveTextContent('処理を中止しました');
    expect(screen.queryByTestId('preview-phase')).not.toBeInTheDocument();
  });

  it('passes processImage to CanceledDisplay', () => {
    const state: AppState = {
      isSlackLinked: true,
      phase: {
        type: 'canceled',
      },
    };

    mockUseAppState.mockReturnValue({
      state,
      dispatch: mockDispatch,
    });

    render(<Main />);

    fireEvent.click(
      screen.getByRole('button', {
        name: 'つづけて撮影',
      }),
    );

    expect(mockProcessImage).toHaveBeenCalledTimes(1);
    expect(mockProcessImage).toHaveBeenCalledWith(expect.any(File));
  });

  it('calls signOut from CanceledDisplay', () => {
    const state: AppState = {
      isSlackLinked: true,
      phase: {
        type: 'canceled',
      },
    };

    mockUseAppState.mockReturnValue({
      state,
      dispatch: mockDispatch,
    });

    render(<Main />);

    fireEvent.click(
      screen.getByRole('button', {
        name: 'ログアウト',
      }),
    );

    expect(mockSignOut).toHaveBeenCalledTimes(1);
  });
});
