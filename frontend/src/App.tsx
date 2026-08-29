import { AuthGate } from './auth/AuthGate';
import { Main } from './components/Main';
import { AppProvider } from './state/AppProvider';

export const App = () => {
  return (
    <AuthGate>
      <AppProvider>
        <Main />
      </AppProvider>
    </AuthGate>
  );
};
