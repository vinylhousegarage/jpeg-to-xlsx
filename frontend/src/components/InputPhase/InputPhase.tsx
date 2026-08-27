import { SlackOAuth } from './SlackOAuth';
import { CameraInput } from './CameraInput';

type Props = {
  isSlackLinked: boolean;
  isCompressing: boolean;
  onConnectSlack: () => void;
  onFileSelected: (
    file: File,
  ) => Promise<void>;
};

export const InputPhase = ({
  isSlackLinked,
  isCompressing,
  onConnectSlack,
  onFileSelected,
}: Props) => {
  if (!isSlackLinked) {
    return (
      <SlackOAuth
        onConnectSlack={onConnectSlack}
      />
    );
  }

  return (
    <div
      style={{
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      }}
    >
      <h2>
        {isCompressing
          ? '画像を処理しています'
          : '画像を撮影'}
      </h2>

      <CameraInput
        onFileSelected={onFileSelected}
        disabled={isCompressing}
      />
    </div>
  );
};
