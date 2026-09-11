type Fetcher = (
  input: RequestInfo | URL,
  init?: RequestInit,
) => Promise<Response>;

type AuthClientDependencies = {
  fetcher: Fetcher;
  redirect: (url: string) => void;
};

type SessionResponse = {
  authenticated: boolean;
};

type LogoutResponse = {
  logout_url: string;
};

export type AuthClient = {
  checkSession: () => Promise<boolean>;
  startSignIn: () => Promise<void>;
  signOut: () => Promise<void>;
};

function isSessionResponse(
  value: unknown,
): value is SessionResponse {
  return (
    typeof value === 'object' &&
    value !== null &&
    'authenticated' in value &&
    typeof value.authenticated === 'boolean'
  );
}

function isLogoutResponse(
  value: unknown,
): value is LogoutResponse {
  return (
    typeof value === 'object' &&
    value !== null &&
    'logout_url' in value &&
    typeof value.logout_url === 'string' &&
    value.logout_url.trim() !== ''
  );
}

export function createAuthClient({
  fetcher,
  redirect,
}: AuthClientDependencies): AuthClient {
  const checkSession =
    async (): Promise<boolean> => {
      const response = await fetcher(
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

      if (response.status === 401) {
        return false;
      }

      if (!response.ok) {
        throw new Error(`Failed to check session: ${response.status}`);
      }

      let body: unknown;

      try {
        body = await response.json();
      } catch {
        throw new Error('Invalid session response');
      }

      if (!isSessionResponse(body)) {
        throw new Error('Invalid session response');
      }

      return body.authenticated;
    };

  const startSignIn = (): Promise<void> => {
    redirect('/api/auth/login');

    return Promise.resolve();
  };

  const signOut = async (): Promise<void> => {
    const response = await fetcher(
      '/api/auth/logout',
      {
        method: 'POST',
        credentials: 'include',
        cache: 'no-store',
      },
    );

    if (!response.ok) {
      throw new Error(`Failed to sign out: ${response.status}`);
    }

    let body: unknown;

    try {
      body = await response.json();
    } catch {
      throw new Error(
        'Invalid logout response',
      );
    }

    if (!isLogoutResponse(body)) {
      throw new Error(
        'Invalid logout response',
      );
    }

    redirect(body.logout_url);
  };

  return {
    checkSession,
    startSignIn,
    signOut,
  };
}

const authClient = createAuthClient({
  fetcher: (
    input: RequestInfo | URL,
    init?: RequestInit,
  ) => fetch(input, init),

  redirect: (url: string) => {
    window.location.assign(url);
  },
});

export const checkSession = authClient.checkSession;
export const startSignIn = authClient.startSignIn;
export const signOut = authClient.signOut;
