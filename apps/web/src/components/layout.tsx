import { ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Button } from '@english/ui';
import { hasValidToken, removeAuthTokens } from '@english/shared';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const isAuthenticated = hasValidToken();

  async function handleLogout() {
    await removeAuthTokens();
    window.location.href = '/'; 
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border">
        <div className="container mx-auto px-4 py-4">
          <nav className="flex items-center justify-between">
            <Link to="/" className="text-xl font-bold text-primary">
              English Learning
            </Link>
            
            <div className="flex items-center gap-4">
              {isAuthenticated ? (
                <>
                  <Link to="/dashboard">
                    <Button 
                      variant={location.pathname === '/dashboard' ? 'default' : 'ghost'}
                    >
                      Dashboard
                    </Button>
                  </Link>
                  <Button variant="outline" onClick={handleLogout}>
                    Logout
                  </Button>
                </>
              ) : (
                <Link to="/login">
                  <Button>
                    Login
                  </Button>
                </Link>
              )}
            </div>
          </nav>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        {children}
      </main>
    </div>
  );
}