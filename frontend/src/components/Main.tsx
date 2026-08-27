import { useEffect } from 'react';
import { useImageProcessor } from '../hooks/useImageProcessor';
import { usePresignUpload } from '../hooks/usePresignUpload';
import { useAppState } from '../state/useContext';
import { createShotNumber } from '../utils/createShotNumber';
import { Spinner } from '../common/Spinner';
import { InputPhase } from './InputPhase';
import { PreviewPhase } from './PreviewPhase';
import { ResultPhase } from './ResultPhase';

export const Main = () => {
  const { state, dispatch } = useAppState();
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

    if (
      url.searchParams.get('slack') !==
      'connected'
    ) {
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
    const apiBaseURL = (
      import.meta.env.VITE_API_BASE_URL ||
      window.location.origin
    ).replace(/\/$/, '');

    window.location.assign(
      `${apiBaseURL}/api/oauth/slack/login`,
    );
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
            send(
              file,
              shotNumber,
            )
          }
        />
      );
    }

    case 'upload':
      return <Spinner />;

    case 'result':
      return (
        <ResultPhase
          state={state.phase}
          dispatch={dispatch}
        />
      );

    default: {
      const _exhaustiveCheck: never =
        state.phase;

      return _exhaustiveCheck;
    }
  }
};
