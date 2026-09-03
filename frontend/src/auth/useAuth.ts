import {
  useCallback,
  useEffect,
  useState,
} from 'react';

import {
  checkSession,
  signOut as signOutSession,
  startSignIn,
} from './authClient';

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

export function useAuth(): UseAuthResult {
  const [status, setStatus] =
    useState<AuthStatus>('checking');
  const [error, setError] =
    useState<Error | null>(null);

  const refreshSession = useCallback(async () => {
    try {
      const authenticated = await checkSession();

      setError(null);
      setStatus(
        authenticated
          ? 'authenticated'
          : 'unauthenticated',
      );
    } catch (sessionError) {
      setError(toError(sessionError));
      setStatus('error');
    }
  }, []);

  const signIn = useCallback(async () => {
    setError(null);
    setStatus('redirecting');

    try {
      await startSignIn();
    } catch (signInError) {
      setError(toError(signInError));
      setStatus('error');
    }
  }, []);

  const signOut = useCallback(async () => {
    setError(null);
    setStatus('checking');

    try {
      await signOutSession();

      setStatus('unauthenticated');
    } catch (signOutError) {
      setError(toError(signOutError));
      setStatus('error');
    }
  }, []);

  useEffect(() => {
    void refreshSession();
  }, [refreshSession]);

  return {
    status,
    error,
    signIn,
    signOut,
  };
}
