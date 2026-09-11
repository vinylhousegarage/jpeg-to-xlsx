// @vitest-environment jsdom

import {
  beforeEach,
  describe,
  expect,
  test,
  vi,
} from 'vitest';

import { createAuthClient } from './authClient';

const logoutURL =
  'https://test.auth.ap-northeast-1.' +
  'amazoncognito.com/logout?' +
  'client_id=test-client-id&' +
  'logout_uri=https%3A%2F%2Fexample.com%2F';

describe('authClient', () => {
  const fetcher = vi.fn();
  const redirect = vi.fn();

  beforeEach(() => {
    fetcher.mockReset();
    redirect.mockReset();
  });

  test(
    'returns true when the BFF session is authenticated',
    async () => {
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
    },
  );

  test('returns false when the BFF session is unauthenticated', async () => {
    fetcher.mockResolvedValue(new Response(JSON.stringify({
      authenticated: false,
    }),
    {
      status: 401,
      headers: {
        'Content-Type': 'application/json',
      },
    }));

    const client = createAuthClient({
      fetcher,
      redirect,
    });

    await expect(client.checkSession()).resolves.toBe(false);
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

    await expect(client.checkSession()).rejects.toThrow('Failed to check session: 500');
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

    await expect(client.checkSession()).rejects.toThrow('Invalid session response');
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

  test('signs out through the BFF endpoint and redirects to Cognito logout', async () => {
    fetcher.mockResolvedValue(
      new Response(
        JSON.stringify({
          logout_url: logoutURL,
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

    expect(redirect).toHaveBeenCalledOnce();
    expect(redirect).toHaveBeenCalledWith(logoutURL);
  });

  test('throws without redirecting when the logout request fails', async () => {
    fetcher.mockResolvedValue(
      new Response(null, {
        status: 500,
      }),
    );

    const client = createAuthClient({
      fetcher,
      redirect,
    });

    await expect(client.signOut()).rejects.toThrow('Failed to sign out: 500');

    expect(redirect).not.toHaveBeenCalled();
  });

  test('throws without redirecting when the logout response is malformed', async () => {
    fetcher.mockResolvedValue(
      new Response(
        '{"logout_url":',
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

    await expect(client.signOut()).rejects.toThrow('Invalid logout response');

    expect(redirect).not.toHaveBeenCalled();
  });

  test.each([
    {
      name: 'missing logout URL',
      body: {},
    },
    {
      name: 'non-string logout URL',
      body: {
        logout_url: 123,
      },
    },
    {
      name: 'empty logout URL',
      body: {
        logout_url: '',
      },
    },
    {
      name: 'whitespace logout URL',
      body: {
        logout_url: '   ',
      },
    },
  ])(
    'throws without redirecting for $name',
    async ({ body }) => {
      fetcher.mockResolvedValue(
        new Response(
          JSON.stringify(body),
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

      await expect(client.signOut()).rejects.toThrow('Invalid logout response');

      expect(redirect).not.toHaveBeenCalled();
    },
  );
});
