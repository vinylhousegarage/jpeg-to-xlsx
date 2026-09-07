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

import { CameraButton } from './CameraButton';

const getFileInput = (
  container: HTMLElement,
): HTMLInputElement => {
  const input =
    container.querySelector<HTMLInputElement>(
      'input[type="file"]',
    );

  if (!input) {
    throw new Error('file input was not found');
  }

  return input;
};

describe('CameraButton', () => {
  it(
    'renders the specified button label',
    () => {
      render(
        <CameraButton
          onFileSelected={
            vi.fn().mockResolvedValue(undefined)
          }
        >
          撮り直し
        </CameraButton>,
      );

      expect(
        screen.getByRole('button', {
          name: '撮り直し',
        }),
      ).toBeInTheDocument();
    },
  );

  it(
    'opens the file input when the button is clicked',
    () => {
      const { container } = render(
        <CameraButton
          onFileSelected={
            vi.fn().mockResolvedValue(undefined)
          }
        >
          撮り直し
        </CameraButton>,
      );

      const input = getFileInput(container);

      const clickSpy = vi
        .spyOn(input, 'click')
        .mockImplementation(
          () => undefined,
        );

      fireEvent.click(
        screen.getByRole('button', {
          name: '撮り直し',
        }),
      );

      expect(clickSpy).toHaveBeenCalledTimes(1);
    },
  );

  it(
    'passes the selected file to onFileSelected',
    async () => {
      const onFileSelected = vi
        .fn()
        .mockResolvedValue(undefined);

      const { container } = render(
        <CameraButton
          onFileSelected={onFileSelected}
        >
          撮り直し
        </CameraButton>,
      );

      const input = getFileInput(container);

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
        expect(onFileSelected).toHaveBeenCalledWith(file);
      });

      expect(onFileSelected).toHaveBeenCalledTimes(1);
    },
  );

  it(
    'does not call onFileSelected when no file is selected',
    () => {
      const onFileSelected = vi.fn().mockResolvedValue(undefined);

      const { container } = render(
        <CameraButton
          onFileSelected={onFileSelected}
        >
          撮り直し
        </CameraButton>,
      );

      const input = getFileInput(container);

      fireEvent.change(input, {
        target: {
          files: [],
        },
      });

      expect(onFileSelected).not.toHaveBeenCalled();
    },
  );

  it(
    'disables the button and file input when disabled',
    () => {
      const { container } = render(
        <CameraButton
          onFileSelected={
            vi.fn().mockResolvedValue(undefined)
          }
          disabled
        >
          撮り直し
        </CameraButton>,
      );

      expect(
        screen.getByRole('button', {
          name: '撮り直し',
        }),
      ).toBeDisabled();

      expect(getFileInput(container)).toBeDisabled();
    },
  );

  it(
    'resets the file input after processing',
    async () => {
      const onFileSelected = vi.fn().mockResolvedValue(undefined);

      const { container } = render(
        <CameraButton
          onFileSelected={onFileSelected}
        >
          撮り直し
        </CameraButton>,
      );

      const input = getFileInput(container);

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
        expect(onFileSelected).toHaveBeenCalledTimes(1);
        expect(input).toHaveValue('');
      });
    },
  );
});
