// Types for extension functionality
export interface PageInfo {
  url: string;
  title: string;
  isYouTube: boolean;
  videoId: string | null;
}

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

export interface TranscriptSegment {
  text: string;
  start: number;
  end: number;
  duration: number;
}

export interface VideoTranscript {
  videoId: string;
  provider: string;
  language: string;
  segments: TranscriptSegment[];
  available: boolean;
  source: string;
}

export interface ExtensionMessage {
  action: string;
  data?: any;
  enabled?: boolean;
  url?: string;
}

export interface MessageResponse {
  success: boolean;
  data?: any;
  message?: string;
  error?: string;
}

export type NotificationType = 'success' | 'error' | 'warning' | 'loading';

export interface NotificationConfig {
  type: NotificationType;
  message: string;
  duration?: number;
}