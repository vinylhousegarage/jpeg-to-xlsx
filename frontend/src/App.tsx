import { AppProvider } from './state/AppProvider';
import { Main } from './components/Main'; 

export const App = () => {
  return (
    <AppProvider>
      <Main />
    </AppProvider>
  );
};
