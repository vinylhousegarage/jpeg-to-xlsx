import {
  fireEvent,
  render,
  screen,
  waitFor,
} from '@testing-library/react';
import {
  afterEach,
  describe,
  expect,
  it,
  vi,
} from 'vitest';

import { PreviewActions } from './PreviewActions';

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

describe('PreviewActions retake', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('opens the file input when the retake button is clicked', () => {
    const inputClickSpy = vi
      .spyOn(
        HTMLInputElement.prototype,
        'click',
      )
      .mockImplementation(() => {});

    render(
      <PreviewActions
        onRetakeFileSelected={vi.fn()}
        onSubmit={vi.fn()}
        onCancel={vi.fn()}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: '撮り直し',
      }),
    );

    expect(inputClickSpy).toHaveBeenCalledTimes(1);
  });

  it('passes the selected file to onRetakeFileSelected', async () => {
    const onRetakeFileSelected = vi.fn().mockResolvedValue(undefined);

    const { container } = render(
      <PreviewActions
        onRetakeFileSelected={
          onRetakeFileSelected
        }
        onSubmit={vi.fn()}
        onCancel={vi.fn()}
      />,
    );

    const input = getFileInput(container);

    const file = new File(
      ['retake-image'],
      'retake.jpg',
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
      expect(onRetakeFileSelected).toHaveBeenCalledTimes(1);
    });

    expect(onRetakeFileSelected).toHaveBeenCalledWith(file);
  });

  it('does not call onRetakeFileSelected when file selection is cancelled', () => {
    const onRetakeFileSelected = vi.fn();

    const { container } = render(
      <PreviewActions
        onRetakeFileSelected={
          onRetakeFileSelected
        }
        onSubmit={vi.fn()}
        onCancel={vi.fn()}
      />,
    );

    const input = getFileInput(container);

    fireEvent.change(input, {
      target: {
        files: [],
      },
    });

    expect(onRetakeFileSelected).not.toHaveBeenCalled();
  });

  it('configures the file input for the environment camera', () => {
    const { container } = render(
      <PreviewActions
        onRetakeFileSelected={vi.fn()}
        onSubmit={vi.fn()}
        onCancel={vi.fn()}
      />,
    );

    const input = getFileInput(container);

    expect(input).toHaveAttribute('accept', 'image/*');
    expect(input).toHaveAttribute('capture', 'environment');
  });
});
