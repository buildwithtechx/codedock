import React from 'react';
import { createRoot } from 'react-dom/client';
import { registerServiceWorker } from './registerServiceWorker';
import { AppRouter } from './router';
import './styles.css';

registerServiceWorker();

createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AppRouter />
  </React.StrictMode>
);
