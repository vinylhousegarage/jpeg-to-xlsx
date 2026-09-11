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

import { PreviewPhase } from './PreviewPhase';

vi.mock('./PreviewArea', () => ({
  PreviewArea: ({
    blob,
  }: {
    blob: Blob;
  }) => (
    <div data-testid="preview-area">
      {blob.type}
    </div>
  ),
}));

vi.mock('./PreviewActions', () => ({
  PreviewActions: ({
    onRetakeFileSelected,
    onSubmit,
    onCancel,
    isSending,
    isCompressing,
  }: {
    onRetakeFileSelected: (
      file: File,
    ) => Promise<void>;
    onSubmit: () => void;
    onCancel: () => void;
    isSending?: boolean;
    isCompressing?: boolean;
  }) => {
    const disabled =
      isSending || isCompressing;

    return (
      <div data-testid="preview-actions">
        <button
          type="button"
          onClick={() => {
            void onRetakeFileSelected(
              new File(
                ['retake-image'],
                'retake.jpg',
                {
                  type: 'image/jpeg',
                },
              ),
            );
          }}
          disabled={disabled}
        >
          撮り直し
        </button>

        <button
          type="button"
          onClick={onSubmit}
          disabled={disabled}
        >
          送信
        </button>

        <button
          type="button"
          onClick={onCancel}
          disabled={disabled}
        >
          中止
        </button>
      </div>
    );
  },
}));

describe('PreviewPhase', () => {
  it('renders the heading, preview area, and preview actions', () => {
    const blob = new Blob(
      ['test-image'],
      {
        type: 'image/jpeg',
      },
    );

    render(
      <PreviewPhase
        blob={blob}
        onRetakeFileSelected={vi.fn()}
        onSend={vi.fn()}
        onCancel={vi.fn()}
      />,
    );

    expect(
      screen.getByRole('heading', {
        name: '画像を確認',
      }),
    ).toBeInTheDocument();

    expect(screen.getByTestId('preview-area')).toHaveTextContent('image/jpeg');
    expect(screen.getByTestId('preview-actions')).toBeInTheDocument();
  });

  it('passes a retake file when the retake button is clicked', () => {
    const onRetakeFileSelected = vi.fn().mockResolvedValue(undefined);
    const onSend = vi.fn();
    const onCancel = vi.fn();

    render(
      <PreviewPhase
        blob={new Blob(['test-image'])}
        onRetakeFileSelected={
          onRetakeFileSelected
        }
        onSend={onSend}
        onCancel={onCancel}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: '撮り直し',
      }),
    );

    expect(onRetakeFileSelected).toHaveBeenCalledTimes(1);
    expect(onRetakeFileSelected).toHaveBeenCalledWith(expect.any(File));
    expect(onSend).not.toHaveBeenCalled();
    expect(onCancel).not.toHaveBeenCalled();
  });

  it('calls onSend when the submit button is clicked', () => {
    const onRetakeFileSelected = vi.fn().mockResolvedValue(undefined);
    const onSend = vi.fn();
    const onCancel = vi.fn();

    render(
      <PreviewPhase
        blob={new Blob(['test-image'])}
        onRetakeFileSelected={
          onRetakeFileSelected
        }
        onSend={onSend}
        onCancel={onCancel}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: '送信',
      }),
    );

    expect(onSend).toHaveBeenCalledTimes(1);
    expect(onRetakeFileSelected).not.toHaveBeenCalled();
    expect(onCancel).not.toHaveBeenCalled();
  });

  it('calls onCancel when the cancel button is clicked', () => {
    const onRetakeFileSelected = vi.fn().mockResolvedValue(undefined);
    const onSend = vi.fn();
    const onCancel = vi.fn();

    render(
      <PreviewPhase
        blob={new Blob(['test-image'])}
        onRetakeFileSelected={
          onRetakeFileSelected
        }
        onSend={onSend}
        onCancel={onCancel}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: '中止',
      }),
    );

    expect(onCancel).toHaveBeenCalledTimes(1);
    expect(onRetakeFileSelected).not.toHaveBeenCalled();
    expect(onSend).not.toHaveBeenCalled();
    },
  );

  it('passes isSending to PreviewActions', () => {
      render(
        <PreviewPhase
          blob={new Blob(['test-image'])}
          onRetakeFileSelected={vi.fn()}
          onSend={vi.fn()}
          onCancel={vi.fn()}
          isSending
        />,
      );

      expect(
        screen.getByRole('button', {
          name: '撮り直し',
        }),
      ).toBeDisabled();

      expect(
        screen.getByRole('button', {
          name: '送信',
        }),
      ).toBeDisabled();

      expect(
        screen.getByRole('button', {
          name: '中止',
        }),
      ).toBeDisabled();
    },
  );

  it('shows the processing state and disables actions while compressing', () => {
    render(
      <PreviewPhase
        blob={new Blob(['test-image'])}
        onRetakeFileSelected={vi.fn()}
        onSend={vi.fn()}
        onCancel={vi.fn()}
        isCompressing
      />,
    );

    expect(
      screen.getByRole('heading', {
        name: '画像を処理しています',
      }),
    ).toBeInTheDocument();

    expect(
      screen.getByRole('button', {
        name: '撮り直し',
      }),
    ).toBeDisabled();

    expect(
      screen.getByRole('button', {
        name: '送信',
      }),
    ).toBeDisabled();

    expect(
      screen.getByRole('button', {
        name: '中止',
      }),
    ).toBeDisabled();
  });
});
