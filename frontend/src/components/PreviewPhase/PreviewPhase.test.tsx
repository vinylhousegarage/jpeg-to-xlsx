import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
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
    isSending,
    isCompressing,
  }: {
    onRetakeFileSelected: (
      file: File,
    ) => Promise<void>;
    onSubmit: () => void;
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
          画像を確定し送信
        </button>
      </div>
    );
  },
}));

describe('PreviewPhase', () => {
  it(
    'renders the heading, preview area, and preview actions',
    () => {
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
        />,
      );

      expect(
        screen.getByRole('heading', {
          name: '画像を確認',
        }),
      ).toBeInTheDocument();

      expect(
        screen.getByTestId('preview-area'),
      ).toHaveTextContent('image/jpeg');

      expect(
        screen.getByTestId('preview-actions'),
      ).toBeInTheDocument();
    },
  );

  it(
    'passes a retake file when the retake button is clicked',
    () => {
      const onRetakeFileSelected =
        vi.fn().mockResolvedValue(undefined);
      const onSend = vi.fn();

      render(
        <PreviewPhase
          blob={new Blob(['test-image'])}
          onRetakeFileSelected={
            onRetakeFileSelected
          }
          onSend={onSend}
        />,
      );

      fireEvent.click(
        screen.getByRole('button', {
          name: '撮り直し',
        }),
      );

      expect(
        onRetakeFileSelected,
      ).toHaveBeenCalledTimes(1);

      expect(
        onRetakeFileSelected,
      ).toHaveBeenCalledWith(
        expect.any(File),
      );

      expect(onSend).not.toHaveBeenCalled();
    },
  );

  it(
    'calls onSend when the submit button is clicked',
    () => {
      const onRetakeFileSelected =
        vi.fn().mockResolvedValue(undefined);
      const onSend = vi.fn();

      render(
        <PreviewPhase
          blob={new Blob(['test-image'])}
          onRetakeFileSelected={
            onRetakeFileSelected
          }
          onSend={onSend}
        />,
      );

      fireEvent.click(
        screen.getByRole('button', {
          name: '画像を確定し送信',
        }),
      );

      expect(onSend).toHaveBeenCalledTimes(1);

      expect(
        onRetakeFileSelected,
      ).not.toHaveBeenCalled();
    },
  );

  it(
    'passes isSending to PreviewActions',
    () => {
      render(
        <PreviewPhase
          blob={new Blob(['test-image'])}
          onRetakeFileSelected={vi.fn()}
          onSend={vi.fn()}
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
          name: '画像を確定し送信',
        }),
      ).toBeDisabled();
    },
  );

  it(
    'shows the processing state and disables actions while compressing',
    () => {
      render(
        <PreviewPhase
          blob={new Blob(['test-image'])}
          onRetakeFileSelected={vi.fn()}
          onSend={vi.fn()}
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
          name: '画像を確定し送信',
        }),
      ).toBeDisabled();
    },
  );
});
