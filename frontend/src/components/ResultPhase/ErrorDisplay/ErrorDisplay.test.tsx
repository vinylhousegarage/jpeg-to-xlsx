import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { ErrorDisplay } from './ErrorDisplay';

describe('ErrorDisplay', () => {
  it(
    'renders the error heading and message',
    () => {
      render(
        <ErrorDisplay
          error={
            new Error(
              '送信処理に失敗しました',
            )
          }
          onContinue={vi.fn()}
          onExit={vi.fn()}
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
          onContinue={vi.fn()}
          onExit={vi.fn()}
        />,
      );

      expect(
        screen.queryByText(
          '送信処理に失敗しました',
        ),
      ).not.toBeInTheDocument();
    },
  );

  it(
    'centers the display and button group',
    () => {
      const { container } = render(
        <ErrorDisplay
          onContinue={vi.fn()}
          onExit={vi.fn()}
        />,
      );

      const display =
        container.querySelector(
          '.error-display',
        );

      expect(display).toHaveStyle({
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      });

      const retakeButton =
        screen.getByRole('button', {
          name: '撮り直し',
        });

      const buttonGroup =
        retakeButton.parentElement;

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
    'calls onContinue when the retake button is clicked',
    () => {
      const onContinue = vi.fn();
      const onExit = vi.fn();

      render(
        <ErrorDisplay
          onContinue={onContinue}
          onExit={onExit}
        />,
      );

      fireEvent.click(
        screen.getByRole('button', {
          name: '撮り直し',
        }),
      );

      expect(
        onContinue,
      ).toHaveBeenCalledTimes(1);

      expect(
        onExit,
      ).not.toHaveBeenCalled();
    },
  );

  it(
    'calls onExit when the exit button is clicked',
    () => {
      const onContinue = vi.fn();
      const onExit = vi.fn();

      render(
        <ErrorDisplay
          onContinue={onContinue}
          onExit={onExit}
        />,
      );

      fireEvent.click(
        screen.getByRole('button', {
          name: '終了',
        }),
      );

      expect(
        onExit,
      ).toHaveBeenCalledTimes(1);

      expect(
        onContinue,
      ).not.toHaveBeenCalled();
    },
  );
});
