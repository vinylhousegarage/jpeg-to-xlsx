import { useReducer, ReactNode } from 'react';
import { AppContext } from './AppContext';
import { appReducer } from './AppReducer';
import { AppState } from '../types';

const initialState: AppState = {
  isSlackLinked: false,
  phase: {
    type: 'input',
  },
};

export const AppProvider = ({
  children,
}: {
  children: ReactNode;
}) => {
  const [state, dispatch] = useReducer(
    appReducer,
    initialState,
  );

  return (
    <AppContext.Provider
      value={{ state, dispatch }}
    >
      {children}
    </AppContext.Provider>
  );
};
