import React from 'react';
import ReactDOM from 'react-dom/client';
import { PopupApp } from './components/popup-app';

// Mount the React app
const container = document.getElementById('app');
if (container) {
  const root = ReactDOM.createRoot(container);
  root.render(<PopupApp />);
}