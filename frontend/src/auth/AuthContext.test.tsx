// @vitest-environment jsdom

import {
  fireEvent,
  render,
  screen,
} from '@testing-library/react';
import {
  beforeEach,
  describe,
  expect,
  test,
  vi,
} from 'vitest';

import {
  AuthProvider,
  useAuthContext,
} from './AuthContext';

const mocks = vi.hoisted(() => ({
  useAuth: vi.fn(),
}));

vi.mock('./useAuth', () => ({
  useAuth: mocks.useAuth,
}));

const TestConsumer = () => {
  const {
    status,
    error,
    signIn,
    signOut,
  } = useAuthContext();

  return (
    <div>
      <p>status: {status}</p>

      {error && (
        <p role="alert">{error.message}</p>
      )}

      <button
        type="button"
        onClick={() => {
          void signIn();
        }}
      >
        sign in
      </button>

      <button
        type="button"
        onClick={() => {
          void signOut();
        }}
      >
        sign out
      </button>
    </div>
  );
};

describe('AuthContext', () => {
  const signIn = vi.fn();
  const signOut = vi.fn();

  beforeEach(() => {
    mocks.useAuth.mockReset();
    signIn.mockReset();
    signOut.mockReset();

    signIn.mockResolvedValue(undefined);
    signOut.mockResolvedValue(undefined);

    mocks.useAuth.mockReturnValue({
      status: 'authenticated',
      error: null,
      signIn,
      signOut,
    });
  });

  test('provides the authentication state to children', () => {
    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    expect(
      screen.getByText(
        'status: authenticated',
      ),
    ).toBeDefined();

    expect(
      screen.queryByRole('alert'),
    ).toBeNull();

    expect(mocks.useAuth).toHaveBeenCalledOnce();
  });

  test('provides the authentication error to children', () => {
    mocks.useAuth.mockReturnValue({
      status: 'error',
      error: new Error(
        'Session check failed',
      ),
      signIn,
      signOut,
    });

    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    expect(
      screen.getByText('status: error'),
    ).toBeDefined();

    expect(
      screen.getByRole('alert').textContent,
    ).toBe('Session check failed');
  });

  test('provides authentication actions to children', () => {
    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: 'sign in',
      }),
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: 'sign out',
      }),
    );

    expect(signIn).toHaveBeenCalledOnce();
    expect(signOut).toHaveBeenCalledOnce();
  });

  test('throws when used outside AuthProvider', () => {
    expect(() => {
      render(<TestConsumer />);
    }).toThrow(
      'useAuthContext must be used within AuthProvider',
    );
  });
});
