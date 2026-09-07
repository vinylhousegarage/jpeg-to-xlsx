import { CameraButton } from '../../CameraButton';
import { standardButtonStyle } from '../../../styles/button';

type Props = {
  onContinueFileSelected: (
    file: File,
  ) => Promise<void>;
  onLogout: () => void;
};

export const SuccessDisplay: React.FC<Props> = ({
  onContinueFileSelected,
  onLogout,
}) => {
  return (
    <div
      className="success-display"
      style={{
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      }}
    >
      <h2>送信完了</h2>

      <div
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
    </div>
  );
};
