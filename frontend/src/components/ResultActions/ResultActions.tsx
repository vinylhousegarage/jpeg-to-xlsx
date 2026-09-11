import { CameraButton } from '../CameraButton';
import { standardButtonStyle } from '../../styles/button';

type Props = {
  onContinueFileSelected: (
    file: File,
  ) => Promise<void>;
  onLogout: () => void;
};

export const ResultActions = ({
  onContinueFileSelected,
  onLogout,
}: Props) => {
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
        onFileSelected={
          onContinueFileSelected
        }
      >
        つづけて撮影
      </CameraButton>

      <button
        type="button"
        onClick={onLogout}
        style={standardButtonStyle}
      >
        ログアウト
      </button>
    </div>
  );
};
