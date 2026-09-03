// @vitest-environment jsdom

import {
  StrictMode,
  type PropsWithChildren,
} from 'react';
import {
  fireEvent,
  render,
  screen,
  waitFor,
} from '@testing-library/react';
import {
  beforeEach,
  describe,
  expect,
  test,
  vi,
} from 'vitest';

import { AuthGate } from './AuthGate';
import type { AuthStatus } from './useAuth';

const mocks = vi.hoisted(() => ({
  useAuth: vi.fn(),
}));

vi.mock('./useAuth', () => ({
  useAuth: mocks.useAuth,
}));

type MockAuthOptions = {
  status: AuthStatus;
  error?: Error | null;
  signIn?: ReturnType<typeof vi.fn>;
  signOut?: ReturnType<typeof vi.fn>;
};

function createMockAuth({
  status,
  error = null,
  signIn = vi.fn().mockResolvedValue(undefined),
  signOut = vi.fn().mockResolvedValue(undefined),
}: MockAuthOptions) {
  return {
    status,
    error,
    signIn,
    signOut,
  };
}

const TestContent = ({ children }: PropsWithChildren) => {
  return <div>{children}</div>;
};

describe('AuthGate', () => {
  beforeEach(() => {
    vi.clearAllMocks();

    mocks.useAuth.mockReturnValue(
      createMockAuth({
        status: 'authenticated',
      }),
    );

    const { unmount } = render(
      <AuthGate>
        <TestContent>
          protected content
        </TestContent>
      </AuthGate>,
    );

    unmount();
    vi.clearAllMocks();
  });

  test('shows the spinner while checking the BFF session', () => {
    mocks.useAuth.mockReturnValue(
      createMockAuth({
        status: 'checking',
      }),
    );

    render(
      <AuthGate>
        <TestContent>
          protected content
        </TestContent>
      </AuthGate>,
    );

    expect(screen.getByText('処理中...')).toBeDefined();
    expect(screen.queryByText('protected content')).toBeNull();
  });

  test('shows the spinner while redirecting to login', () => {
    const signIn = vi.fn().mockResolvedValue(undefined);

    mocks.useAuth.mockReturnValue(
      createMockAuth({
        status: 'redirecting',
        signIn,
      }),
    );

    render(
      <AuthGate>
        <TestContent>
          protected content
        </TestContent>
      </AuthGate>,
    );

    expect(screen.getByText('処理中...')).toBeDefined();
    expect(signIn).not.toHaveBeenCalled();
    expect(screen.queryByText('protected content')).toBeNull();
  });

  test('starts BFF login when unauthenticated', async () => {
    const signIn = vi.fn().mockResolvedValue(undefined);

    mocks.useAuth.mockReturnValue(
      createMockAuth({
        status: 'unauthenticated',
        signIn,
      }),
    );

    render(
      <AuthGate>
        <TestContent>
          protected content
        </TestContent>
      </AuthGate>,
    );

    await waitFor(() => {
      expect(signIn).toHaveBeenCalledOnce();
    });

    expect(screen.getByText('処理中...')).toBeDefined();
    expect(screen.queryByText('protected content')).toBeNull();
  });

  test('does not start duplicate login redirects in StrictMode', async () => {
    const signIn = vi.fn().mockResolvedValue(undefined);

    mocks.useAuth.mockReturnValue(
      createMockAuth({
        status: 'unauthenticated',
        signIn,
      }),
    );

    render(
      <StrictMode>
        <AuthGate>
          <TestContent>
            protected content
          </TestContent>
        </AuthGate>
      </StrictMode>,
    );

    await waitFor(() => {
      expect(signIn).toHaveBeenCalledOnce();
    });
  });

  test('shows children when authenticated', () => {
    mocks.useAuth.mockReturnValue(
      createMockAuth({
        status: 'authenticated',
      }),
    );

    render(
      <AuthGate>
        <TestContent>
          protected content
        </TestContent>
      </AuthGate>,
    );

    expect(screen.getByText('protected content')).toBeDefined();
    expect(screen.queryByText('処理中...')).toBeNull();
  });

  test('shows an authentication error', () => {
    const authError = new Error('Session check failed');

    mocks.useAuth.mockReturnValue(
      createMockAuth({
        status: 'error',
        error: authError,
      }),
    );

    render(
      <AuthGate>
        <TestContent>
          protected content
        </TestContent>
      </AuthGate>,
    );

    expect(screen.getByText('認証エラー')).toBeDefined();
    expect(screen.getByText('Session check failed')).toBeDefined();
    expect(screen.queryByText('protected content')).toBeNull();
  });

  test('retries BFF login after an error', () => {
    const signIn = vi.fn().mockResolvedValue(undefined);

    mocks.useAuth.mockReturnValue(
      createMockAuth({
        status: 'error',
        error: new Error('Session check failed'),
        signIn,
      }),
    );

    render(
      <AuthGate>
        <TestContent>
          protected content
        </TestContent>
      </AuthGate>,
    );

    fireEvent.click(screen.getByRole('button', {
      name: '再試行',
    }));

    expect(signIn).toHaveBeenCalledOnce();
  });
});
