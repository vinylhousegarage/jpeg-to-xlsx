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
      const onRetakeFileSelected = vi.fn();
      const onSubmit = vi.fn();
      const onCancel = vi.fn();

      render(
        <PreviewActions
          onRetakeFileSelected={onRetakeFileSelected}
          onSubmit={onSubmit}
          onCancel={onCancel}
        />,
      );

      fireEvent.click(
        screen.getByRole('button', {
          name: '送信',
        }),
      );

      expect(onSubmit).toHaveBeenCalledTimes(1);
      expect(onRetakeFileSelected).not.toHaveBeenCalled();
      expect(onCancel).not.toHaveBeenCalled();
    },
  );

  it(
    'calls onCancel when the cancel button is clicked',
    () => {
      const onRetakeFileSelected = vi.fn();
      const onSubmit = vi.fn();
      const onCancel = vi.fn();

      render(
        <PreviewActions
          onRetakeFileSelected={onRetakeFileSelected}
          onSubmit={onSubmit}
          onCancel={onCancel}
        />,
      );

      fireEvent.click(
        screen.getByRole('button', {
          name: '中止',
        }),
      );

      expect(onCancel).toHaveBeenCalledTimes(1);
      expect(onSubmit).not.toHaveBeenCalled();
      expect(onRetakeFileSelected).not.toHaveBeenCalled();
    },
  );

  it(
    'calls onRetakeFileSelected when an image is selected',
    () => {
      const onRetakeFileSelected = vi.fn();
      const onSubmit = vi.fn();
      const onCancel = vi.fn();

      const { container } = render(
        <PreviewActions
          onRetakeFileSelected={onRetakeFileSelected}
          onSubmit={onSubmit}
          onCancel={onCancel}
        />,
      );

      const fileInput = container.querySelector<HTMLInputElement>('input[type="file"]');

      expect(fileInput).not.toBeNull();

      const file = new File(
        ['image-data'],
        'retake.jpg',
        {
          type: 'image/jpeg',
        },
      );

      fireEvent.change(fileInput!, {
        target: {
          files: [file],
        },
      });

      expect(onRetakeFileSelected).toHaveBeenCalledTimes(1);
      expect(onRetakeFileSelected).toHaveBeenCalledWith(file);
      expect(onSubmit).not.toHaveBeenCalled();
      expect(onCancel).not.toHaveBeenCalled();
    },
  );

  it(
    'does not call onRetakeFileSelected when no file is selected',
    () => {
      const onRetakeFileSelected = vi.fn();

      const { container } = render(
        <PreviewActions
          onRetakeFileSelected={onRetakeFileSelected}
          onSubmit={vi.fn()}
          onCancel={vi.fn()}
        />,
      );

      const fileInput = container.querySelector<HTMLInputElement>('input[type="file"]');

      expect(fileInput).not.toBeNull();

      fireEvent.change(fileInput!, {
        target: {
          files: [],
        },
      });

      expect(onRetakeFileSelected).not.toHaveBeenCalled();
    },
  );

  it(
    'disables all buttons while sending',
    () => {
      const onRetakeFileSelected = vi.fn();
      const onSubmit = vi.fn();
      const onCancel = vi.fn();

      render(
        <PreviewActions
          onRetakeFileSelected={onRetakeFileSelected}
          onSubmit={onSubmit}
          onCancel={onCancel}
          isSending
        />,
      );

      const retakeButton = screen.getByRole(
        'button',
        {
          name: '撮り直し',
        },
      );

      const submitButton = screen.getByRole(
        'button',
        {
          name: '送信',
        },
      );

      const cancelButton = screen.getByRole(
        'button',
        {
          name: '中止',
        },
      );

      expect(retakeButton).toBeDisabled();
      expect(submitButton).toBeDisabled();
      expect(cancelButton).toBeDisabled();

      fireEvent.click(retakeButton);
      fireEvent.click(submitButton);
      fireEvent.click(cancelButton);

      expect(onRetakeFileSelected).not.toHaveBeenCalled();
      expect(onSubmit).not.toHaveBeenCalled();
      expect(onCancel).not.toHaveBeenCalled();
    },
  );

  it(
    'disables all buttons while compressing',
    () => {
      const onRetakeFileSelected = vi.fn();
      const onSubmit = vi.fn();
      const onCancel = vi.fn();

      render(
        <PreviewActions
          onRetakeFileSelected={onRetakeFileSelected}
          onSubmit={onSubmit}
          onCancel={onCancel}
          isCompressing
        />,
      );

      const retakeButton = screen.getByRole(
        'button',
        {
          name: '撮り直し',
        },
      );

      const submitButton = screen.getByRole(
        'button',
        {
          name: '送信',
        },
      );

      const cancelButton = screen.getByRole(
        'button',
        {
          name: '中止',
        },
      );

      expect(retakeButton).toBeDisabled();
      expect(submitButton).toBeDisabled();
      expect(cancelButton).toBeDisabled();

      fireEvent.click(retakeButton);
      fireEvent.click(submitButton);
      fireEvent.click(cancelButton);

      expect(onRetakeFileSelected).not.toHaveBeenCalled();
      expect(onSubmit).not.toHaveBeenCalled();
      expect(onCancel).not.toHaveBeenCalled();
    },
  );

  it(
    'enables all buttons by default',
    () => {
      render(
        <PreviewActions
          onRetakeFileSelected={vi.fn()}
          onSubmit={vi.fn()}
          onCancel={vi.fn()}
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

      expect(
        screen.getByRole('button', {
          name: '中止',
        }),
      ).toBeEnabled();
    },
  );
});
