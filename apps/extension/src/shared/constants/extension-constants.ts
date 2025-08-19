// Extension constants for consistent usage across modules
export const EXTENSION_STORAGE_KEYS = {
  API_BASE_URL: 'apiBaseUrl',
  AUTH_TOKEN: 'authToken',
  AUTO_TRANSLATE_ENABLED: 'autoTranslateEnabled'
} as const;

export const EXTENSION_MESSAGES = {
  GET_PAGE_INFO: 'GET_PAGE_INFO',
  TOGGLE_TRANSLATION: 'TOGGLE_TRANSLATION',
  GET_TRANSCRIPT_WITH_AUTH: 'GET_TRANSCRIPT_WITH_AUTH',
  OPEN_TAB: 'OPEN_TAB'
} as const;

// Backend API configuration
export const DEFAULT_API_BASE_URL = 'http://localhost:8000/api/v1';

// Video URL patterns for different providers
export const VIDEO_URL_PATTERNS: Record<string, RegExp[]> = {
  YOUTUBE: [/(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]{11})/]
};

// YouTube selectors for extension integration
export const YOUTUBE_SELECTORS = {
  RIGHT_CONTROLS: [
    '#movie_player .ytp-right-controls',
    '.ytp-chrome-bottom .ytp-right-controls', 
    '.ytp-chrome-controls .ytp-right-controls',
    '.html5-video-player .ytp-right-controls',
    'div.ytp-right-controls'
  ],
  SETTINGS_BUTTON: '.ytp-settings-button',
  MOVIE_PLAYER: '#movie_player',
  VIDEO_PLAYER: '.html5-video-player'
} as const;

// Extension specific CSS classes
export const EXTENSION_CLASSES = {
  EXTENSION_BTN: 'english-learning-extension-btn',
  PANEL: 'english-learning-panel',
  MODAL: 'english-learning-transcript-modal',
  NOTIFICATION: 'english-learning-notification',
  AUTH_MODAL: 'oauth-authentication-modal'
} as const;

// Notification types and timing
export const NOTIFICATION_TYPES = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  LOADING: 'loading'
} as const;

export const NOTIFICATION_DURATIONS = {
  SUCCESS: 3000,
  ERROR: 5000,
  WARNING: 4000,
  LOADING: 0 // Don't auto-hide loading notifications
} as const;

// OAuth polling configuration
export const OAUTH_CONFIG = {
  POLL_INTERVAL: 1000,
  MAX_POLLS: 60,
  POLL_TIMEOUT: 60000
} as const;