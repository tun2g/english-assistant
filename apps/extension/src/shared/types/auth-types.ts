// Types for OAuth authentication
export interface OAuthAuthenticationResult {
  authenticated: boolean;
  error?: string;
}

export interface OAuthAuthUrlResult {
  authUrl: string;
  state: string;
}

export interface OAuthRevokeResult {
  success: boolean;
  error?: string;
}
