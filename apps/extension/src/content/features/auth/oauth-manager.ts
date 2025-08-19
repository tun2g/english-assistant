import { oauthService } from '../../../services/oauth-service';
import { NotificationManager } from '../notifications/notification-manager';
import { OAUTH_CONFIG } from '../../../shared/constants/extension-constants';
import type { OAuthAuthenticationResult } from '../../../shared/types/extension-types';

export class OAuthManager {
  private isAuthenticated: boolean = false;
  private pollingInterval: NodeJS.Timeout | null = null;

  constructor() {
    this.checkInitialAuthStatus();
  }

  // Check initial OAuth authentication status
  private async checkInitialAuthStatus(): Promise<void> {
    try {
      const result = await oauthService.checkAuthStatus();
      this.isAuthenticated = result.authenticated;
      
      if (result.authenticated) {
        console.log('OAuthManager: OAuth authentication active');
      } else {
        console.log('OAuthManager: OAuth authentication required for real captions');
      }
    } catch (error) {
      console.error('Failed to check OAuth status:', error);
      this.isAuthenticated = false;
    }
  }

  // Get current authentication status
  get authenticated(): boolean {
    return this.isAuthenticated;
  }

  // Check OAuth status (async refresh)
  async checkAuthStatus(): Promise<OAuthAuthenticationResult> {
    try {
      const result = await oauthService.checkAuthStatus();
      this.isAuthenticated = result.authenticated;
      return result;
    } catch (error) {
      console.error('Failed to check OAuth status:', error);
      this.isAuthenticated = false;
      return { 
        authenticated: false, 
        error: error instanceof Error ? error.message : 'Unknown error'
      };
    }
  }

  // Initiate OAuth flow
  async connect(): Promise<void> {
    try {
      NotificationManager.showLoading('Opening YouTube authentication...');
      await oauthService.initiateOAuth();
      
      // Start polling for authentication completion
      this.startAuthenticationPolling();
    } catch (error) {
      console.error('OAuth connection failed:', error);
      NotificationManager.showError('Failed to start authentication: ' + (error as Error).message);
      throw error;
    }
  }

  // Revoke OAuth authorization
  async disconnect(): Promise<boolean> {
    try {
      const result = await oauthService.revokeAuth();
      if (result.success) {
        this.isAuthenticated = false;
        NotificationManager.showSuccess('YouTube account disconnected');
        return true;
      } else {
        NotificationManager.showError('Failed to disconnect: ' + (result.error || 'Unknown error'));
        return false;
      }
    } catch (error) {
      console.error('OAuth disconnect failed:', error);
      NotificationManager.showError('Failed to disconnect account');
      return false;
    }
  }

  // Start polling to detect when authentication completes
  private startAuthenticationPolling(): void {
    if (this.pollingInterval) {
      clearInterval(this.pollingInterval);
    }

    let pollCount = 0;
    
    this.pollingInterval = setInterval(async () => {
      pollCount++;
      
      try {
        const result = await this.checkAuthStatus();
        
        if (result.authenticated) {
          // Authentication successful!
          this.stopPolling();
          NotificationManager.showSuccess('YouTube authentication successful! 🎉');
        } else if (pollCount >= OAUTH_CONFIG.MAX_POLLS) {
          // Timeout
          this.stopPolling();
          NotificationManager.showWarning('Authentication timeout. Please try again.');
        }
      } catch (error) {
        console.error('Authentication polling error:', error);
        if (pollCount >= OAUTH_CONFIG.MAX_POLLS) {
          this.stopPolling();
        }
      }
    }, OAUTH_CONFIG.POLL_INTERVAL);

    // Additional timeout safety
    setTimeout(() => {
      this.stopPolling();
    }, OAUTH_CONFIG.POLL_TIMEOUT);
  }

  // Stop polling for authentication
  private stopPolling(): void {
    if (this.pollingInterval) {
      clearInterval(this.pollingInterval);
      this.pollingInterval = null;
    }
  }

  // Check if error is authentication-related
  isAuthenticationError(error: any): boolean {
    const errorMessage = error.message || error.details || '';
    return errorMessage.toLowerCase().includes('authenticate') ||
           errorMessage.toLowerCase().includes('unauthorized') ||
           errorMessage.toLowerCase().includes('auth') ||
           error.status === 401 ||
           error.status === 403;
  }

  // Clean up resources
  destroy(): void {
    this.stopPolling();
  }
}