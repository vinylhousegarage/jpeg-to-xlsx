import { useAuthContext } from '../../auth/AuthContext';
import {
  ResultPhase as ResultPhaseState,
} from '../../types';
import { ErrorDisplay } from './ErrorDisplay';
import { SuccessDisplay } from './SuccessDisplay';

type Props = {
  state: ResultPhaseState;
  onFileSelected: (
    file: File,
  ) => Promise<void>;
};

export const ResultPhase = ({
  state,
  onFileSelected,
}: Props) => {
  const { signOut } = useAuthContext();

  const handleLogout = () => {
    void signOut();
  };

  if (state.status === 'success') {
    return (
      <SuccessDisplay
        onContinueFileSelected={
          onFileSelected
        }
        onLogout={handleLogout}
      />
    );
  }

  return (
    <ErrorDisplay
      error={state.error}
      onRetakeFileSelected={
        onFileSelected
      }
      onLogout={handleLogout}
    />
  );
};
