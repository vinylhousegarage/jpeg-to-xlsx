import { useContext } from 'react';
import { AppContext, AppContextType } from './AppContext';

export const useAppState = (): AppContextType => {
  const context = useContext(AppContext);

  if (context === undefined) {
    throw new Error('useAppState must be used within an AppProvider');
  }

  return context;
};
