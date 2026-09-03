// @vitest-environment jsdom

import {
  act,
  renderHook,
  waitFor,
} from '@testing-library/react';
import {
  beforeEach,
  describe,
  expect,
  test,
  vi,
} from 'vitest';
import { useAuth } from './useAuth';

const mocks = vi.hoisted(() => ({
  checkSession: vi.fn(),
  startSignIn: vi.fn(),
  signOut: vi.fn(),
}));

vi.mock('./authClient', () => ({
  checkSession: mocks.checkSession,
  startSignIn: mocks.startSignIn,
  signOut: mocks.signOut,
}));

describe('useAuth', () => {
  beforeEach(() => {
    vi.clearAllMocks();

    mocks.checkSession.mockResolvedValue(false);
    mocks.startSignIn.mockResolvedValue(undefined);
    mocks.signOut.mockResolvedValue(undefined);
  });

  test('sets authenticated when the BFF session is valid', async () => {
    mocks.checkSession.mockResolvedValue(true);

    const { result } = renderHook(
      () => useAuth(),
    );

    expect(result.current.status).toBe('checking');

    await waitFor(() => {
      expect(result.current.status).toBe('authenticated');
    });

    expect(result.current.error).toBeNull();

    expect(mocks.checkSession).toHaveBeenCalledOnce();
  });

  test('sets unauthenticated when the BFF session is invalid', async () => {
    mocks.checkSession.mockResolvedValue(false);

    const { result } = renderHook(
      () => useAuth(),
    );

    await waitFor(() => {
      expect(result.current.status).toBe('unauthenticated');
    });

    expect(result.current.error).toBeNull();

    expect(mocks.checkSession).toHaveBeenCalledOnce();
  });

  test('sets error when the session check fails', async () => {
    const sessionError = new Error('Session check failed');

    mocks.checkSession.mockRejectedValue(sessionError);

    const { result } = renderHook(
      () => useAuth(),
    );

    await waitFor(() => {
      expect(result.current.status).toBe('error');
    });

    expect(result.current.error).toBe(sessionError);
  });

  test('converts a non-Error session failure', async () => {
    mocks.checkSession.mockRejectedValue('Session check failed');

    const { result } = renderHook(
      () => useAuth(),
    );

    await waitFor(() => {
      expect(result.current.status).toBe('error');
    });

    expect(result.current.error).toEqual(new Error('Authentication failed'));
  });

  test('starts the BFF login redirect', async () => {
    const { result } = renderHook(
      () => useAuth(),
    );

    await waitFor(() => {
      expect(result.current.status).toBe('unauthenticated');
    });

    await act(async () => {
      await result.current.signIn();
    });

    expect(mocks.startSignIn).toHaveBeenCalledOnce();

    expect(result.current.status).toBe('redirecting');

    expect(result.current.error).toBeNull();
  });

  test('sets error when the BFF login redirect fails', async () => {
    const redirectError = new Error('Login redirect failed');

    mocks.startSignIn.mockRejectedValue(redirectError);

    const { result } = renderHook(
      () => useAuth(),
    );

    await waitFor(() => {
      expect(result.current.status).toBe('unauthenticated');
    });

    await act(async () => {
      await result.current.signIn();
    });

    expect(result.current.status).toBe('error');

    expect(result.current.error).toBe(redirectError);
  });

  test('signs out the BFF session', async () => {
    mocks.checkSession.mockResolvedValue(true);

    const { result } = renderHook(
      () => useAuth(),
    );

    await waitFor(() => {
      expect(result.current.status).toBe('authenticated');
    });

    await act(async () => {
      await result.current.signOut();
    });

    expect(mocks.signOut).toHaveBeenCalledOnce();

    expect(result.current.status).toBe('unauthenticated');

    expect(result.current.error).toBeNull();
  });

  test('sets error when signing out fails', async () => {
    const signOutError = new Error('Sign out failed');

    mocks.checkSession.mockResolvedValue(true);
    mocks.signOut.mockRejectedValue(signOutError);

    const { result } = renderHook(
      () => useAuth(),
    );

    await waitFor(() => {
      expect(result.current.status).toBe('authenticated');
    });

    await act(async () => {
      await result.current.signOut();
    });

    expect(result.current.status).toBe('error');

    expect(result.current.error).toBe(signOutError);
  });
});
