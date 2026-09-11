import {
  fireEvent,
  render,
  screen,
} from '@testing-library/react';
import {
  describe,
  expect,
  it,
  vi,
} from 'vitest';

import { CanceledDisplay } from './CanceledDisplay';

vi.mock('../ResultActions', () => ({
  ResultActions: ({
    onContinueFileSelected,
    onLogout,
  }: {
    onContinueFileSelected: (
      file: File,
    ) => Promise<void>;
    onLogout: () => void;
  }) => (
    <div data-testid="result-actions">
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

describe('CanceledDisplay', () => {
  it('displays the canceled message and result actions', () => {
    render(
      <CanceledDisplay
        onContinueFileSelected={vi.fn()}
        onLogout={vi.fn()}
      />,
    );

    expect(
      screen.getByRole('heading', {
        name: '処理を中止しました',
      }),
    ).toBeInTheDocument();

    expect(screen.getByTestId('result-actions')).toBeInTheDocument();
  });

  it('passes a selected file to onContinueFileSelected', () => {
    const onContinueFileSelected = vi.fn().mockResolvedValue(undefined);
    const onLogout = vi.fn();

    render(
      <CanceledDisplay
        onContinueFileSelected={
          onContinueFileSelected
        }
        onLogout={onLogout}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: 'つづけて撮影',
      }),
    );

    expect(onContinueFileSelected).toHaveBeenCalledTimes(1);
    expect(onContinueFileSelected).toHaveBeenCalledWith(expect.any(File));
    expect(onLogout).not.toHaveBeenCalled();
  });

  it('calls onLogout when the logout button is clicked', () => {
    const onContinueFileSelected = vi.fn();
    const onLogout = vi.fn();

    render(
      <CanceledDisplay
        onContinueFileSelected={
          onContinueFileSelected
        }
        onLogout={onLogout}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: 'ログアウト',
      }),
    );

    expect(onLogout).toHaveBeenCalledTimes(1);
    expect(onContinueFileSelected).not.toHaveBeenCalled();
  });
});
