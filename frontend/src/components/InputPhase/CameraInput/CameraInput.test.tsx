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
import { CameraInput } from './CameraInput';

describe('CameraInput', () => {
  it('renders the camera input label', () => {
    render(
      <CameraInput
        onFileSelected={vi.fn()}
      />,
    );

    expect(
      screen.getByLabelText('カメラを起動'),
    ).toBeInTheDocument();
  });

  it('calls onFileSelected with the selected file', async () => {
    const onFileSelected = vi.fn().mockResolvedValue(undefined);

    render(
      <CameraInput
        onFileSelected={onFileSelected}
      />,
    );

    const input = screen.getByLabelText(
      'カメラを起動',
    ) as HTMLInputElement;

    const file = new File(
      ['test-image'],
      'test.jpg',
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
    });

    expect(onFileSelected).toHaveBeenCalledWith(file);
  });

  it('does not call onFileSelected when no file is selected', () => {
    const onFileSelected = vi.fn().mockResolvedValue(undefined);

    render(
      <CameraInput
        onFileSelected={onFileSelected}
      />,
    );

    const input = screen.getByLabelText(
      'カメラを起動',
    );

    fireEvent.change(input, {
      target: {
        files: [],
      },
    });

    expect(onFileSelected).not.toHaveBeenCalled();
  });

  it('disables the file input when disabled is true', () => {
    render(
      <CameraInput
        onFileSelected={vi.fn()}
        disabled
      />,
    );

    expect(
      screen.getByLabelText('カメラを起動'),
    ).toBeDisabled();
  });

  it('enables the file input by default', () => {
    render(
      <CameraInput
        onFileSelected={vi.fn()}
      />,
    );

    expect(
      screen.getByLabelText('カメラを起動'),
    ).toBeEnabled();
  });

  it('resets the input value after processing the selected file', async () => {
    const onFileSelected = vi.fn().mockResolvedValue(undefined);

    render(
      <CameraInput
        onFileSelected={onFileSelected}
      />,
    );

    const input = screen.getByLabelText(
      'カメラを起動',
    ) as HTMLInputElement;

    const file = new File(
      ['test-image'],
      'test.jpg',
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

    expect(input.value).toBe('');
  });
});
