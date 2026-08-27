import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { SuccessDisplay } from './SuccessDisplay';

describe('SuccessDisplay', () => {
  it(
    'renders the success heading',
    () => {
      render(
        <SuccessDisplay
          onContinue={vi.fn()}
          onExit={vi.fn()}
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
          onContinue={vi.fn()}
          onExit={vi.fn()}
        />,
      );

      const display =
        container.querySelector(
          '.success-display',
        );

      expect(display).toHaveStyle({
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      });

      const continueButton =
        screen.getByRole('button', {
          name: 'つづけて撮影',
        });

      const buttonGroup =
        continueButton.parentElement;

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
    'calls onContinue when the continue button is clicked',
    () => {
      const onContinue = vi.fn();
      const onExit = vi.fn();

      render(
        <SuccessDisplay
          onContinue={onContinue}
          onExit={onExit}
        />,
      );

      fireEvent.click(
        screen.getByRole('button', {
          name: 'つづけて撮影',
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
        <SuccessDisplay
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
