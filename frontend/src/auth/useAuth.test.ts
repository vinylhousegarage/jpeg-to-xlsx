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
  getCurrentUser: vi.fn(),
  signInWithRedirect: vi.fn(),
  amplifySignOut: vi.fn(),
  hubListen: vi.fn(),
}));

vi.mock('aws-amplify/auth', () => ({
  getCurrentUser: mocks.getCurrentUser,
  signInWithRedirect: mocks.signInWithRedirect,
  signOut: mocks.amplifySignOut,
}));

vi.mock('aws-amplify/utils', () => ({
  Hub: {
    listen: mocks.hubListen,
  },
}));

type AuthHubEvent = {
  payload: {
    event: string;
    data?: unknown;
  };
};

type AuthHubCallback = (
  event: AuthHubEvent,
) => void;

describe('useAuth', () => {
  let authHubCallback:
    | AuthHubCallback
    | undefined;
  let stopListening: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    vi.clearAllMocks();

    authHubCallback = undefined;
    stopListening = vi.fn();

    mocks.hubListen.mockImplementation(
      (
        _channel: string,
        callback: AuthHubCallback,
      ) => {
        authHubCallback = callback;

        return stopListening;
      },
    );
  });

  test('sets authenticated when a current user exists', async () => {
    mocks.getCurrentUser.mockResolvedValue({
      userId: 'test-user-id',
      username: 'test-user',
    });

    const { result } = renderHook(() => useAuth());

    expect(result.current.status).toBe('checking');

    await waitFor(() => {
      expect(result.current.status).toBe(
        'authenticated',
      );
    });

    expect(result.current.error).toBeNull();
    expect(mocks.getCurrentUser).toHaveBeenCalledOnce();
    expect(mocks.hubListen).toHaveBeenCalledWith(
      'auth',
      expect.any(Function),
    );
  });

  test('sets unauthenticated when a current user does not exist', async () => {
    const unauthenticatedError = new Error(
      'User is not authenticated',
    );

    unauthenticatedError.name =
      'UserUnAuthenticatedException';

    mocks.getCurrentUser.mockRejectedValue(
      unauthenticatedError,
    );

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.status).toBe(
        'unauthenticated',
      );
    });

    expect(result.current.error).toBeNull();
  });

  test('starts the managed login redirect', async () => {
    const unauthenticatedError = new Error(
      'User is not authenticated',
    );

    unauthenticatedError.name =
      'UserUnAuthenticatedException';

    mocks.getCurrentUser.mockRejectedValue(
      unauthenticatedError,
    );
    mocks.signInWithRedirect.mockResolvedValue(
      undefined,
    );

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.status).toBe(
        'unauthenticated',
      );
    });

    await act(async () => {
      await result.current.signIn();
    });

    expect(
      mocks.signInWithRedirect,
    ).toHaveBeenCalledOnce();

    expect(result.current.status).toBe(
      'redirecting',
    );
    expect(result.current.error).toBeNull();
  });

  test('sets error when the managed login redirect fails', async () => {
    const unauthenticatedError = new Error(
      'User is not authenticated',
    );
    const redirectError = new Error(
      'Redirect failed',
    );

    unauthenticatedError.name =
      'UserUnAuthenticatedException';

    mocks.getCurrentUser.mockRejectedValue(
      unauthenticatedError,
    );
    mocks.signInWithRedirect.mockRejectedValue(
      redirectError,
    );

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.status).toBe(
        'unauthenticated',
      );
    });

    await act(async () => {
      await result.current.signIn();
    });

    expect(result.current.status).toBe('error');
    expect(result.current.error).toBe(
      redirectError,
    );
  });

  test('checks the session after a signed-in event', async () => {
    const unauthenticatedError = new Error(
      'User is not authenticated',
    );

    unauthenticatedError.name =
      'UserUnAuthenticatedException';

    mocks.getCurrentUser
      .mockRejectedValueOnce(
        unauthenticatedError,
      )
      .mockResolvedValueOnce({
        userId: 'test-user-id',
        username: 'test-user',
      });

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.status).toBe(
        'unauthenticated',
      );
    });

    if (!authHubCallback) {
      throw new Error(
        'Auth Hub callback was not registered',
      );
    }

    act(() => {
      authHubCallback?.({
        payload: {
          event: 'signedIn',
        },
      });
    });

    await waitFor(() => {
      expect(result.current.status).toBe(
        'authenticated',
      );
    });

    expect(mocks.getCurrentUser).toHaveBeenCalledTimes(
      2,
    );
  });

  test('sets error after a redirect failure event', async () => {
    mocks.getCurrentUser.mockResolvedValue({
      userId: 'test-user-id',
      username: 'test-user',
    });

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.status).toBe(
        'authenticated',
      );
    });

    if (!authHubCallback) {
      throw new Error(
        'Auth Hub callback was not registered',
      );
    }

    const redirectError = new Error(
      'OAuth callback failed',
    );

    act(() => {
      authHubCallback?.({
        payload: {
          event:
            'signInWithRedirect_failure',
          data: redirectError,
        },
      });
    });

    expect(result.current.status).toBe('error');
    expect(result.current.error).toBe(
      redirectError,
    );
  });

  test('signs out the current user', async () => {
    mocks.getCurrentUser.mockResolvedValue({
      userId: 'test-user-id',
      username: 'test-user',
    });
    mocks.amplifySignOut.mockResolvedValue(
      undefined,
    );

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.status).toBe(
        'authenticated',
      );
    });

    await act(async () => {
      await result.current.signOut();
    });

    expect(
      mocks.amplifySignOut,
    ).toHaveBeenCalledOnce();

    expect(result.current.status).toBe(
      'unauthenticated',
    );
    expect(result.current.error).toBeNull();
  });

  test('stops listening to auth events when unmounted', async () => {
    mocks.getCurrentUser.mockResolvedValue({
      userId: 'test-user-id',
      username: 'test-user',
    });

    const { result, unmount } = renderHook(
      () => useAuth(),
    );

    await waitFor(() => {
      expect(result.current.status).toBe(
        'authenticated',
      );
    });

    unmount();

    expect(stopListening).toHaveBeenCalledOnce();
  });
});
