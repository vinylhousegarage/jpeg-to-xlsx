import {
  fireEvent,
  render,
  screen,
  waitFor,
} from '@testing-library/react';
import {
  describe,
  expect,
  it,
  vi,
} from 'vitest';

import { ErrorDisplay } from './ErrorDisplay';

describe('ErrorDisplay', () => {
  it(
    'renders the error heading and message',
    () => {
      render(
        <ErrorDisplay
          error={
            new Error('送信処理に失敗しました')
          }
          onRetakeFileSelected={
            vi.fn().mockResolvedValue(undefined)
          }
          onLogout={vi.fn()}
        />,
      );

      expect(
        screen.getByRole('heading', {
          name: '送信失敗',
        }),
      ).toBeInTheDocument();

      expect(
        screen.getByText(
          '送信処理に失敗しました',
        ),
      ).toBeInTheDocument();
    },
  );

  it(
    'does not render an error message when error is undefined',
    () => {
      render(
        <ErrorDisplay
          onRetakeFileSelected={
            vi.fn().mockResolvedValue(undefined)
          }
          onLogout={vi.fn()}
        />,
      );

      expect(screen.queryByText('送信処理に失敗しました')).not.toBeInTheDocument();
    },
  );

  it(
    'centers the display and button group',
    () => {
      const { container } = render(
        <ErrorDisplay
          onRetakeFileSelected={
            vi.fn().mockResolvedValue(undefined)
          }
          onLogout={vi.fn()}
        />,
      );

      const display = container.querySelector('.error-display');

      expect(display).toHaveStyle({
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      });

      const retakeButton =
        screen.getByRole('button', {
          name: '撮り直し',
        });

      const buttonGroup = retakeButton.parentElement;

      expect(buttonGroup).not.toBeNull();

      expect(buttonGroup).toHaveStyle({
        display: 'flex',
        justifyContent: 'center',
        gap: '10px',
        width: '100%',
      });
    },
  );

  it(
    'passes the selected file to onRetakeFileSelected',
    async () => {
      const onRetakeFileSelected = vi
        .fn()
        .mockResolvedValue(undefined);

      const onLogout = vi.fn();

      const { container } = render(
        <ErrorDisplay
          onRetakeFileSelected={
            onRetakeFileSelected
          }
          onLogout={onLogout}
        />,
      );

      const input = container.querySelector<HTMLInputElement>('input[type="file"]' );

      expect(input).not.toBeNull();

      const file = new File(
        ['image data'],
        'photo.jpg',
        {
          type: 'image/jpeg',
        },
      );

      fireEvent.change(input!, {
        target: {
          files: [file],
        },
      });

      await waitFor(() => {
        expect(
          onRetakeFileSelected,
        ).toHaveBeenCalledWith(file);
      });

      expect(onRetakeFileSelected).toHaveBeenCalledTimes(1);
      expect(onLogout).not.toHaveBeenCalled();
    },
  );

  it(
    'calls onLogout when the logout button is clicked',
    () => {
      const onRetakeFileSelected = vi.fn().mockResolvedValue(undefined);
      const onLogout = vi.fn();

      render(
        <ErrorDisplay
          onRetakeFileSelected={
            onRetakeFileSelected
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
      expect(onRetakeFileSelected).not.toHaveBeenCalled();
    },
  );
});
