// @vitest-environment jsdom

import {
  beforeEach,
  describe,
  expect,
  test,
  vi,
} from 'vitest';

import { createAuthClient } from './authClient';

describe('authClient', () => {
  const fetcher = vi.fn();
  const redirect = vi.fn();

  beforeEach(() => {
    fetcher.mockReset();
    redirect.mockReset();
  });

  test('returns true when the BFF session is authenticated', async () => {
    fetcher.mockResolvedValue(
      new Response(
        JSON.stringify({
          authenticated: true,
        }),
        {
          status: 200,
          headers: {
            'Content-Type': 'application/json',
          },
        },
      ),
    );

    const client = createAuthClient({
      fetcher,
      redirect,
    });

    await expect(client.checkSession()).resolves.toBe(true);

    expect(fetcher).toHaveBeenCalledOnce();

    expect(fetcher).toHaveBeenCalledWith(
      '/api/auth/session',
      {
        method: 'GET',
        credentials: 'include',
        cache: 'no-store',
        headers: {
          Accept: 'application/json',
        },
      },
    );
  });

  test('returns false when the BFF session is unauthenticated', async () => {
    fetcher.mockResolvedValue(
      new Response(
        JSON.stringify({
          authenticated: false,
        }),
        {
          status: 401,
          headers: {
            'Content-Type': 'application/json',
          },
        },
      ),
    );

    const client = createAuthClient({
      fetcher,
      redirect,
    });

    await expect(
      client.checkSession(),
    ).resolves.toBe(false);
  });

  test('throws when the session request fails', async () => {
    fetcher.mockResolvedValue(
      new Response(null, {
        status: 500,
      }),
    );

    const client = createAuthClient({
      fetcher,
      redirect,
    });

    await expect(client.checkSession()).rejects.toThrow(
      'Failed to check session: 500',
    );
  });

  test('throws when the session response is invalid', async () => {
    fetcher.mockResolvedValue(
      new Response(
        JSON.stringify({
          authenticated: 'true',
        }),
        {
          status: 200,
          headers: {
            'Content-Type': 'application/json',
          },
        },
      ),
    );

    const client = createAuthClient({
      fetcher,
      redirect,
    });

    await expect(client.checkSession()).rejects.toThrow(
      'Invalid session response',
    );
  });

  test('redirects to the BFF login endpoint', async () => {
    const client = createAuthClient({
      fetcher,
      redirect,
    });

    await client.startSignIn();

    expect(redirect).toHaveBeenCalledOnce();
    expect(redirect).toHaveBeenCalledWith('/api/auth/login');
    expect(fetcher).not.toHaveBeenCalled();
  });

  test('signs out through the BFF endpoint', async () => {
    fetcher.mockResolvedValue(
      new Response(null, {
        status: 204,
      }),
    );

    const client = createAuthClient({
      fetcher,
      redirect,
    });

    await client.signOut();

    expect(fetcher).toHaveBeenCalledOnce();

    expect(fetcher).toHaveBeenCalledWith(
      '/api/auth/logout',
      {
        method: 'POST',
        credentials: 'include',
        cache: 'no-store',
      },
    );
  });

  test('throws when the logout request fails', async () => {
    fetcher.mockResolvedValue(
      new Response(null, {
        status: 500,
      }),
    );

    const client = createAuthClient({
      fetcher,
      redirect,
    });

    await expect(client.signOut()).rejects.toThrow(
      'Failed to sign out: 500',
    );
  });
});
