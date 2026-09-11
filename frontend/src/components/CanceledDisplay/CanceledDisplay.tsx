import { ResultActions } from '../ResultActions';

type Props = {
  onContinueFileSelected: (
    file: File,
  ) => Promise<void>;
  onLogout: () => void;
};

export const CanceledDisplay = ({
  onContinueFileSelected,
  onLogout,
}: Props) => {
  return (
    <div
      className="canceled-display"
      style={{
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      }}
    >
      <h2>処理を中止しました</h2>

      <ResultActions
        onContinueFileSelected={
          onContinueFileSelected
        }
        onLogout={onLogout}
      />
    </div>
  );
};
