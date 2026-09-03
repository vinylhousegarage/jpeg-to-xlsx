import {
  createContext,
  type PropsWithChildren,
  useContext,
} from 'react';

import {
  type AuthStatus,
  useAuth,
} from './useAuth';

type AuthContextValue = {
  status: AuthStatus;
  error: Error | null;
  signIn: () => Promise<void>;
  signOut: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export const AuthProvider = ({
  children,
}: PropsWithChildren) => {
  const auth = useAuth();

  return (
    <AuthContext.Provider value={auth}>
      {children}
    </AuthContext.Provider>
  );
};

export function useAuthContext(): AuthContextValue {
  const context = useContext(AuthContext);

  if (context === null) {
    throw new Error(
      'useAuthContext must be used within AuthProvider',
    );
  }

  return context;
}
