import {
  getCurrentUser,
  signInWithRedirect,
  signOut as amplifySignOut,
} from 'aws-amplify/auth';
import { Hub } from 'aws-amplify/utils';
import {
  useCallback,
  useEffect,
  useState,
} from 'react';

export type AuthStatus =
  | 'checking'
  | 'unauthenticated'
  | 'redirecting'
  | 'authenticated'
  | 'error';

type UseAuthResult = {
  status: AuthStatus;
  error: Error | null;
  signIn: () => Promise<void>;
  signOut: () => Promise<void>;
};

function toError(value: unknown): Error {
  if (value instanceof Error) {
    return value;
  }

  return new Error('Authentication failed');
}

function isUnauthenticatedError(
  error: unknown,
): boolean {
  return (
    error instanceof Error &&
    error.name === 'UserUnAuthenticatedException'
  );
}

export function useAuth(): UseAuthResult {
  const [status, setStatus] =
    useState<AuthStatus>('checking');
  const [error, setError] =
    useState<Error | null>(null);

  const checkSession = useCallback(async () => {
    try {
      await getCurrentUser();

      setError(null);
      setStatus('authenticated');
    } catch (sessionError) {
      if (isUnauthenticatedError(sessionError)) {
        setError(null);
        setStatus('unauthenticated');

        return;
      }

      setError(toError(sessionError));
      setStatus('error');
    }
  }, []);

  const signIn = useCallback(async () => {
    setError(null);
    setStatus('redirecting');

    try {
      await signInWithRedirect();
    } catch (signInError) {
      setError(toError(signInError));
      setStatus('error');
    }
  }, []);

  const signOut = useCallback(async () => {
    setError(null);
    setStatus('checking');

    try {
      await amplifySignOut();
      setStatus('unauthenticated');
    } catch (signOutError) {
      setError(toError(signOutError));
      setStatus('error');
    }
  }, []);

  useEffect(() => {
    const stopListening = Hub.listen(
      'auth',
      ({ payload }) => {
        switch (payload.event) {
          case 'signedIn':
          case 'signInWithRedirect':
            void checkSession();
            break;

          case 'signedOut':
          case 'tokenRefresh_failure':
            setError(null);
            setStatus('unauthenticated');
            break;

          case 'signInWithRedirect_failure':
            setError(toError(payload.data));
            setStatus('error');
            break;

          default:
            break;
        }
      },
    );

    void checkSession();

    return stopListening;
  }, [checkSession]);

  return {
    status,
    error,
    signIn,
    signOut,
  };
}
