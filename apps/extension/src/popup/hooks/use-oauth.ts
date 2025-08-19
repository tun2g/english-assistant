import { useState, useEffect } from 'react';
import { oauthService } from '../../services/oauth-service';

export function useOAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        const result = await oauthService.checkAuthStatus();
        setIsAuthenticated(result.authenticated);
      } catch (error) {
        console.error('Failed to check OAuth status:', error);
        setIsAuthenticated(false);
      }
    };

    checkAuthStatus();
  }, []);

  const connect = async (): Promise<void> => {
    setIsLoading(true);
    try {
      await oauthService.initiateOAuth();
      
      // Poll for authentication completion
      const pollInterval = setInterval(async () => {
        const result = await oauthService.checkAuthStatus();
        if (result.authenticated) {
          clearInterval(pollInterval);
          setIsAuthenticated(true);
          setIsLoading(false);
        }
      }, 2000);

      // Stop polling after 60 seconds
      setTimeout(() => {
        clearInterval(pollInterval);
        setIsLoading(false);
      }, 60000);
    } catch (error) {
      console.error('OAuth connection failed:', error);
      setIsLoading(false);
      throw error;
    }
  };

  const disconnect = async (): Promise<void> => {
    setIsLoading(true);
    try {
      const result = await oauthService.revokeAuth();
      if (result.success) {
        setIsAuthenticated(false);
      }
    } catch (error) {
      console.error('OAuth disconnect failed:', error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    isAuthenticated,
    isLoading,
    connect,
    disconnect
  };
}