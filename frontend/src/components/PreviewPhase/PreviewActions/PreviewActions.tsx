import { CameraButton } from '../../CameraButton';
import { standardButtonStyle } from '../../../styles/button';

type Props = {
  onRetakeFileSelected: (
    file: File,
  ) => Promise<void>;
  onSubmit: () => void;
  onCancel: () => void;
  isSending?: boolean;
  isCompressing?: boolean;
};

export const PreviewActions = ({
  onRetakeFileSelected,
  onSubmit,
  onCancel,
  isSending = false,
  isCompressing = false,
}: Props) => {
  const disabled = isSending || isCompressing;

  return (
    <div
      className="button-group"
      style={{
        display: 'flex',
        justifyContent: 'center',
        gap: '10px',
        width: '100%',
      }}
    >
      <CameraButton
        onFileSelected={onRetakeFileSelected}
        disabled={disabled}
      >
        撮り直し
      </CameraButton>

      <button
        type="button"
        onClick={onSubmit}
        disabled={disabled}
        style={standardButtonStyle}
      >
        送信
      </button>

      <button
        type="button"
        onClick={onCancel}
        disabled={disabled}
        style={standardButtonStyle}
      >
        中止
      </button>
    </div>
  );
};
