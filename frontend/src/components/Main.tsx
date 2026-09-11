import { useEffect } from 'react';
import { useAuthContext } from '../auth/AuthContext';
import { Spinner } from '../common/Spinner';
import { useImageProcessor } from '../hooks/useImageProcessor';
import { usePresignUpload } from '../hooks/usePresignUpload';
import { useAppState } from '../state/useContext';
import { createShotNumber } from '../utils';
import { CanceledDisplay } from './CanceledDisplay';
import { InputPhase } from './InputPhase';
import { PreviewPhase } from './PreviewPhase';
import { ResultPhase } from './ResultPhase';

export const Main = () => {
  const { state, dispatch } = useAppState();
  const { signOut } = useAuthContext();
  const { send } = usePresignUpload(dispatch);

  const handleCapture = (blob: Blob) => {
    const shotNumber =
      state.phase.type === 'preview'
        ? state.phase.shotNumber
        : createShotNumber();

    dispatch({
      type: 'SET_PREVIEW',
      file: blob,
      shotNumber,
    });
  };

  const {
    isCompressing,
    processImage,
  } = useImageProcessor(
    handleCapture,
    () => {},
  );

  useEffect(() => {
    const url = new URL(window.location.href);

    if (url.searchParams.get('slack') !== 'connected') {
      return;
    }

    dispatch({
      type: 'SET_SLACK_LINKED',
      isSlackLinked: true,
    });

    url.searchParams.delete('slack');

    window.history.replaceState(
      {},
      '',
      `${url.pathname}${url.search}${url.hash}`,
    );
  }, [dispatch]);

  const handleConnectSlack = () => {
    window.location.assign('/api/oauth/slack/login');
  };

  const handleCancel = () => {
    dispatch({
      type: 'CANCEL',
    });
  };

  const handleLogout = () => {
    void signOut();
  };

  switch (state.phase.type) {
    case 'input':
      return (
        <InputPhase
          isSlackLinked={state.isSlackLinked}
          isCompressing={isCompressing}
          onConnectSlack={handleConnectSlack}
          onFileSelected={processImage}
        />
      );

    case 'preview': {
      const { file, shotNumber } = state.phase;

      return (
        <PreviewPhase
          blob={file}
          isSending={false}
          isCompressing={isCompressing}
          onRetakeFileSelected={processImage}
          onSend={() =>
            send(file, shotNumber)
          }
          onCancel={handleCancel}
        />
      );
    }

    case 'canceled':
      return (
        <CanceledDisplay
          onContinueFileSelected={processImage}
          onLogout={handleLogout}
        />
      );

    case 'upload':
      return <Spinner />;

    case 'result':
      return (
        <ResultPhase
          state={state.phase}
          onFileSelected={processImage}
        />
      );

    default: {
      const _exhaustiveCheck: never = state.phase;

      return _exhaustiveCheck;
    }
  }
};
