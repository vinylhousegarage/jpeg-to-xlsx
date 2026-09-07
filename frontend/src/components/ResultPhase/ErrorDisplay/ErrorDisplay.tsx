import { CameraButton } from '../../CameraButton';
import { standardButtonStyle } from '../../../styles/button';

type Props = {
  error?: Error;
  onRetakeFileSelected: (
    file: File,
  ) => Promise<void>;
  onLogout: () => void;
};

export const ErrorDisplay: React.FC<Props> = ({
  error,
  onRetakeFileSelected,
  onLogout,
}) => {
  return (
    <div
      className="error-display"
      style={{
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      }}
    >
      <h2>送信失敗</h2>

      {error && (
        <p style={{ color: 'red' }}>
          {error.message}
        </p>
      )}

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
            onRetakeFileSelected
          }
        >
          撮り直し
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
