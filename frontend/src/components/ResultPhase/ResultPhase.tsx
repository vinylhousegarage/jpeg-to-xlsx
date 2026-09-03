import { Dispatch } from 'react';

import { useAuthContext } from '../../auth/AuthContext';
import {
  AppAction,
  ResultPhase as ResultPhaseState,
} from '../../types';
import { ErrorDisplay } from './ErrorDisplay';
import { SuccessDisplay } from './SuccessDisplay';

type Props = {
  state: ResultPhaseState;
  dispatch: Dispatch<AppAction>;
};

export const ResultPhase = ({
  state,
  dispatch,
}: Props) => {
  const { signOut } = useAuthContext();

  const handleContinue = () => {
    dispatch({ type: 'CONTINUE' });
  };

  const handleLogout = () => {
    void signOut();
  };

  const handleExit = () => {
    dispatch({ type: 'EXIT' });
  };

  if (state.status === 'success') {
    return (
      <SuccessDisplay
        onContinue={handleContinue}
        onLogout={handleLogout}
      />
    );
  }

  return (
    <ErrorDisplay
      error={state.error}
      onContinue={handleContinue}
      onExit={handleExit}
    />
  );
};
