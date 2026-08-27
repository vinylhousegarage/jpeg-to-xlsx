import { createContext, Dispatch } from 'react';
import { AppAction, AppState } from '../types';

export type AppContextType = {
  state: AppState;
  dispatch: Dispatch<AppAction>;
};

export const AppContext = createContext<AppContextType | undefined>(
  undefined,
);
