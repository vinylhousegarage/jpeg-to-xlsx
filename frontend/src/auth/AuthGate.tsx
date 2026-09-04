import {
  type PropsWithChildren,
  useEffect,
} from 'react';

import { Spinner } from '../common';
import { standardButtonStyle } from '../styles/button';
import { useAuthContext } from './AuthContext';

let signInRedirectStarted = false;

export const AuthGate = ({
  children,
}: PropsWithChildren) => {
  const {
    status,
    error,
    signIn,
  } = useAuthContext();

  useEffect(() => {
    if (
      status === 'authenticated' ||
      status === 'error'
    ) {
      signInRedirectStarted = false;

      return;
    }

    if (
      status !== 'unauthenticated' ||
      signInRedirectStarted
    ) {
      return;
    }

    signInRedirectStarted = true;

    void signIn();
  }, [signIn, status]);

  if (status === 'authenticated') {
    return <>{children}</>;
  }

  if (status === 'error') {
    const handleRetry = () => {
      signInRedirectStarted = true;

      void signIn();
    };

    return (
      <div
        style={{
          maxWidth: '375px',
          margin: '0 auto',
          textAlign: 'center',
        }}
      >
        <h2>認証エラー</h2>

        <p>
          認証に失敗しました。もう一度お試しください。
        </p>

        {error && (
          <p
            role="alert"
            style={{
              color: '#b00020',
            }}
          >
            {error.message}
          </p>
        )}

        <button
          type="button"
          onClick={handleRetry}
          style={standardButtonStyle}
        >
          再試行
        </button>
      </div>
    );
  }

  return <Spinner />;
};
