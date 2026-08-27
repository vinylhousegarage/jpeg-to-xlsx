import { standardButtonStyle } from '../../../styles/button';

type Props = {
  error?: Error;
  onContinue: () => void;
  onExit: () => void;
};

export const ErrorDisplay: React.FC<Props> = ({
  error,
  onContinue,
  onExit,
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
        <button
          type="button"
          onClick={onContinue}
          style={standardButtonStyle}
        >
          撮り直し
        </button>

        <button
          type="button"
          onClick={onExit}
          style={standardButtonStyle}
        >
          終了
        </button>
      </div>
    </div>
  );
};
