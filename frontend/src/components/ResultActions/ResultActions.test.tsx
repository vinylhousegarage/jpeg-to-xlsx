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

import { ResultActions } from './ResultActions';

describe('ResultActions', () => {
  const onContinueFileSelected = vi.fn();
  const onLogout = vi.fn();

  beforeEach(() => {
    onContinueFileSelected.mockReset();
    onLogout.mockReset();

    onContinueFileSelected.mockResolvedValue(
      undefined,
    );
  });

  it(
    'renders the continue and logout buttons',
    () => {
      render(
        <ResultActions
          onContinueFileSelected={
            onContinueFileSelected
          }
          onLogout={onLogout}
        />,
      );

      expect(
        screen.getByRole('button', {
          name: 'つづけて撮影',
        }),
      ).toBeInTheDocument();

      expect(
        screen.getByRole('button', {
          name: 'ログアウト',
        }),
      ).toBeInTheDocument();
    },
  );

  it(
    'passes the selected file to onContinueFileSelected',
    async () => {
      const { container } = render(
        <ResultActions
          onContinueFileSelected={
            onContinueFileSelected
          }
          onLogout={onLogout}
        />,
      );

      const input = container.querySelector<HTMLInputElement>('input[type="file"]');

      if (!input) {
        throw new Error('File input was not found');
      }

      const file = new File(
        ['image data'],
        'photo.jpg',
        {
          type: 'image/jpeg',
        },
      );

      fireEvent.change(input, {
        target: {
          files: [file],
        },
      });

      await waitFor(() => {
        expect(onContinueFileSelected).toHaveBeenCalledWith(file);
      });

      expect(onContinueFileSelected).toHaveBeenCalledOnce();
      expect(onLogout).not.toHaveBeenCalled();
    },
  );

  it(
    'calls onLogout when the logout button is clicked',
    () => {
      render(
        <ResultActions
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

      expect(onLogout).toHaveBeenCalledOnce();
      expect(onContinueFileSelected).not.toHaveBeenCalled();
    },
  );
});
