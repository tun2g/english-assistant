// OAuth authentication service for Chrome extension
export class OAuthService {
  private baseUrl: string;
  
  constructor(baseUrl: string = 'http://localhost:8000/api/v1') {
    this.baseUrl = baseUrl;
  }

  // Check if user is authenticated with YouTube OAuth
  async checkAuthStatus(): Promise<{ authenticated: boolean; error?: string }> {
    try {
      const response = await fetch(`${this.baseUrl}/oauth/youtube/status`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'omit', // Don't send cookies for OAuth status check
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();
      return { authenticated: data.authenticated };
    } catch (error) {
      console.error('Failed to check OAuth status:', error);
      return { 
        authenticated: false, 
        error: error instanceof Error ? error.message : 'Unknown error'
      };
    }
  }

  // Get OAuth authorization URL
  async getAuthUrl(): Promise<{ authUrl: string; state: string } | { error: string }> {
    try {
      const response = await fetch(`${this.baseUrl}/oauth/youtube/auth`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'omit', // Don't send cookies - not needed for OAuth URL generation
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();
      return { authUrl: data.authUrl, state: data.state };
    } catch (error) {
      console.error('Failed to get auth URL:', error);
      return { 
        error: error instanceof Error ? error.message : 'Failed to get authorization URL'
      };
    }
  }

  // Open OAuth flow in new tab
  async initiateOAuth(): Promise<void> {
    const result = await this.getAuthUrl();
    
    if ('error' in result) {
      throw new Error(result.error);
    }

    // Content scripts can't use chrome.tabs.create directly
    // Send message to background script to open the tab
    if (typeof chrome !== 'undefined' && chrome.runtime) {
      chrome.runtime.sendMessage({
        action: 'OPEN_TAB',
        url: result.authUrl
      });
    } else {
      // Fallback: open in same window (not ideal for OAuth)
      window.open(result.authUrl, '_blank');
    }
  }

  // Revoke OAuth authorization
  async revokeAuth(): Promise<{ success: boolean; error?: string }> {
    try {
      const response = await fetch(`${this.baseUrl}/oauth/youtube/revoke`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'omit', // Server-side OAuth management doesn't need browser cookies
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      return { success: true };
    } catch (error) {
      console.error('Failed to revoke OAuth:', error);
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Failed to revoke authorization'
      };
    }
  }
}

// Singleton instance
export const oauthService = new OAuthService();