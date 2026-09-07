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

import { PreviewActions } from './PreviewActions';

describe('PreviewActions', () => {
  it(
    'calls onSubmit when the submit button is clicked',
    () => {
      const onRetakeFileSelected =
        vi.fn();
      const onSubmit = vi.fn();

      render(
        <PreviewActions
          onRetakeFileSelected={
            onRetakeFileSelected
          }
          onSubmit={onSubmit}
        />,
      );

      fireEvent.click(
        screen.getByRole('button', {
          name: '送信',
        }),
      );

      expect(onSubmit).toHaveBeenCalledTimes(1);
      expect(onRetakeFileSelected).not.toHaveBeenCalled();
    },
  );

  it(
    'disables both buttons while sending',
    () => {
      const onRetakeFileSelected =
        vi.fn();
      const onSubmit = vi.fn();

      render(
        <PreviewActions
          onRetakeFileSelected={
            onRetakeFileSelected
          }
          onSubmit={onSubmit}
          isSending
        />,
      );

      const retakeButton =
        screen.getByRole('button', {
          name: '撮り直し',
        });

      const submitButton =
        screen.getByRole('button', {
          name: '送信',
        });

      expect(retakeButton).toBeDisabled();
      expect(submitButton).toBeDisabled();

      fireEvent.click(retakeButton);
      fireEvent.click(submitButton);

      expect(onRetakeFileSelected).not.toHaveBeenCalled();
      expect(onSubmit).not.toHaveBeenCalled();
    },
  );

  it(
    'disables both buttons while compressing',
    () => {
      const onRetakeFileSelected =
        vi.fn();
      const onSubmit = vi.fn();

      render(
        <PreviewActions
          onRetakeFileSelected={
            onRetakeFileSelected
          }
          onSubmit={onSubmit}
          isCompressing
        />,
      );

      const retakeButton = screen.getByRole('button', {
        name: '撮り直し',
      });

      const submitButton = screen.getByRole('button', {
        name: '送信',
      });

      expect(retakeButton).toBeDisabled();
      expect(submitButton).toBeDisabled();

      fireEvent.click(retakeButton);
      fireEvent.click(submitButton);

      expect(onRetakeFileSelected).not.toHaveBeenCalled();
      expect(onSubmit).not.toHaveBeenCalled();
    },
  );

  it(
    'enables both buttons by default',
    () => {
      render(
        <PreviewActions
          onRetakeFileSelected={vi.fn()}
          onSubmit={vi.fn()}
        />,
      );

      expect(
        screen.getByRole('button', {
          name: '撮り直し',
        }),
      ).toBeEnabled();

      expect(
        screen.getByRole('button', {
          name: '送信',
        }),
      ).toBeEnabled();
    },
  );
});
