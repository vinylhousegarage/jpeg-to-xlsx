import React from 'react';
import ReactDOM from 'react-dom/client';
import 'aws-amplify/auth/enable-oauth-listener';
import { App } from './App';
import { configureAmplify } from './config/amplify';
import './index.css';

configureAmplify();

ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement,
).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
