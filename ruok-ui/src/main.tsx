import React from 'react';
import ReactDOM from 'react-dom/client';
import { HashRouter as Router } from 'react-router-dom';
import { CssVarsProvider } from '@mui/joy/styles';
import CssBaseline from '@mui/joy/CssBaseline';
import App from './App.tsx';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <CssVarsProvider
      defaultMode="dark" // the selector to apply CSS theme variables stylesheet.
      colorSchemeSelector="#dark-mode-by-default"
      //
      // the local storage key to use
      modeStorageKey="dark-mode-by-default"
      //
      // set as root provider
      //disableNestedContext
    >
      <CssBaseline />
      <Router>
        <App />
      </Router>
    </CssVarsProvider>
  </React.StrictMode>,
);
