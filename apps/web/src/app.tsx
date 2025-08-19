import { Routes, Route } from 'react-router-dom';
import { useEffect } from 'react';
import { initializeTokens } from '@english/shared';
import { HomePage } from './pages/home-page';
import { LoginPage } from './pages/login-page';
import { DashboardPage } from './pages/dashboard-page';
import { Layout } from './components/layout';

export function App() {
  useEffect(() => {
    // Initialize auth tokens on app start
    async function init() {
      await initializeTokens();
    }
    init();
  }, []);

  return (
    <Layout>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/dashboard" element={<DashboardPage />} />
      </Routes>
    </Layout>
  );
}