import { Dispatch } from 'react';
import { AppAction, ResultPhase as ResultPhaseState } from '../../types';
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
  const handleContinue = () => {
    dispatch({ type: 'CONTINUE' });
  };

  const handleExit = () => {
    dispatch({ type: 'EXIT' });
  };

  if (state.status === 'success') {
    return (
      <SuccessDisplay
        onContinue={handleContinue}
        onExit={handleExit}
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
