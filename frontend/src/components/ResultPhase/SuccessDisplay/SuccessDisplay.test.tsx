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

import { SuccessDisplay } from './SuccessDisplay';

describe('SuccessDisplay', () => {
  it(
    'renders the success heading',
    () => {
      render(
        <SuccessDisplay
          onContinueFileSelected={
            vi.fn().mockResolvedValue(undefined)
          }
          onLogout={vi.fn()}
        />,
      );

      expect(
        screen.getByRole('heading', {
          name: '送信完了',
        }),
      ).toBeInTheDocument();
    },
  );

  it(
    'centers the display and button group',
    () => {
      const { container } = render(
        <SuccessDisplay
          onContinueFileSelected={
            vi.fn().mockResolvedValue(undefined)
          }
          onLogout={vi.fn()}
        />,
      );

      const display = container.querySelector('.success-display');

      expect(display).toHaveStyle({
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      });

      const continueButton = screen.getByRole('button', {
        name: 'つづけて撮影',
      });

      const buttonGroup = continueButton.parentElement;

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
    'passes the selected file to onContinueFileSelected',
    async () => {
      const onContinueFileSelected = vi
        .fn()
        .mockResolvedValue(undefined);

      const onLogout = vi.fn();

      const { container } = render(
        <SuccessDisplay
          onContinueFileSelected={
            onContinueFileSelected
          }
          onLogout={onLogout}
        />,
      );

      const input = container.querySelector<HTMLInputElement>(
        'input[type="file"]',
      );

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
          onContinueFileSelected,
        ).toHaveBeenCalledWith(file);
      });

      expect(onContinueFileSelected).toHaveBeenCalledTimes(1);
      expect(onLogout).not.toHaveBeenCalled();
    },
  );

  it(
    'calls onLogout when the logout button is clicked',
    () => {
      const onContinueFileSelected = vi
        .fn()
        .mockResolvedValue(undefined);

      const onLogout = vi.fn();

      render(
        <SuccessDisplay
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
    },
  );
});
