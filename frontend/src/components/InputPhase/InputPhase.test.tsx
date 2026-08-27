import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { InputPhase } from './InputPhase';

vi.mock('./SlackOAuth', () => ({
  SlackOAuth: ({
    onConnectSlack,
  }: {
    onConnectSlack: () => void;
  }) => (
    <button
      type="button"
      onClick={onConnectSlack}
    >
      DM通知を設定
    </button>
  ),
}));

vi.mock('./CameraInput', () => ({
  CameraInput: ({
    onFileSelected,
    disabled,
  }: {
    onFileSelected: (
      file: File,
    ) => Promise<void>;
    disabled?: boolean;
  }) => (
    <input
      aria-label="カメラを起動"
      type="file"
      disabled={disabled}
      onChange={(event) => {
        const file =
          event.currentTarget.files?.[0];

        if (file) {
          void onFileSelected(file);
        }
      }}
    />
  ),
}));

describe('InputPhase', () => {
  it(
    'renders SlackOAuth when Slack is not linked',
    () => {
      render(
        <InputPhase
          isSlackLinked={false}
          isCompressing={false}
          onConnectSlack={vi.fn()}
          onFileSelected={vi.fn()}
        />,
      );

      expect(
        screen.getByRole('button', {
          name: 'DM通知を設定',
        }),
      ).toBeInTheDocument();

      expect(
        screen.queryByLabelText(
          'カメラを起動',
        ),
      ).not.toBeInTheDocument();
    },
  );

  it(
    'calls onConnectSlack when the Slack button is clicked',
    () => {
      const onConnectSlack = vi.fn();

      render(
        <InputPhase
          isSlackLinked={false}
          isCompressing={false}
          onConnectSlack={onConnectSlack}
          onFileSelected={vi.fn()}
        />,
      );

      fireEvent.click(
        screen.getByRole('button', {
          name: 'DM通知を設定',
        }),
      );

      expect(
        onConnectSlack,
      ).toHaveBeenCalledTimes(1);
    },
  );

  it(
    'renders CameraInput when Slack is linked',
    () => {
      render(
        <InputPhase
          isSlackLinked
          isCompressing={false}
          onConnectSlack={vi.fn()}
          onFileSelected={vi.fn()}
        />,
      );

      expect(
        screen.getByRole('heading', {
          name: '画像を撮影',
        }),
      ).toBeInTheDocument();

      expect(
        screen.getByLabelText(
          'カメラを起動',
        ),
      ).toBeEnabled();

      expect(
        screen.queryByRole('button', {
          name: 'DM通知を設定',
        }),
      ).not.toBeInTheDocument();
    },
  );

  it(
    'shows processing state and disables CameraInput',
    () => {
      render(
        <InputPhase
          isSlackLinked
          isCompressing
          onConnectSlack={vi.fn()}
          onFileSelected={vi.fn()}
        />,
      );

      expect(
        screen.getByRole('heading', {
          name: '画像を処理しています',
        }),
      ).toBeInTheDocument();

      expect(
        screen.getByLabelText(
          'カメラを起動',
        ),
      ).toBeDisabled();
    },
  );

  it(
    'passes the selected file to onFileSelected',
    () => {
      const onFileSelected =
        vi.fn().mockResolvedValue(undefined);

      const file = new File(
        ['image'],
        'capture.jpg',
        {
          type: 'image/jpeg',
        },
      );

      render(
        <InputPhase
          isSlackLinked
          isCompressing={false}
          onConnectSlack={vi.fn()}
          onFileSelected={onFileSelected}
        />,
      );

      fireEvent.change(
        screen.getByLabelText(
          'カメラを起動',
        ),
        {
          target: {
            files: [file],
          },
        },
      );

      expect(
        onFileSelected,
      ).toHaveBeenCalledTimes(1);

      expect(
        onFileSelected,
      ).toHaveBeenCalledWith(file);
    },
  );
});
