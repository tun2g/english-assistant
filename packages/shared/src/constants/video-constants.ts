// Video service constants

export const VIDEO_PROVIDERS = {
  YOUTUBE: 'youtube',
  VIMEO: 'vimeo',
  TWITCH: 'twitch',
} as const;

export const VIDEO_ACTIONS = {
  GET_INFO: 'GET_VIDEO_INFO',
  GET_TRANSCRIPT: 'GET_TRANSCRIPT',
  TRANSLATE: 'TRANSLATE_TRANSCRIPT',
  GET_LANGUAGES: 'GET_AVAILABLE_LANGUAGES',
  GET_CAPABILITIES: 'GET_CAPABILITIES',
} as const;

export const TRANSCRIPT_SOURCES = {
  MANUAL: 'manual',
  AUTO_GENERATED: 'auto-generated',
  FORCED: 'forced',
} as const;

export const SUPPORTED_LANGUAGES = [
  { code: 'en', name: 'English', flag: '🇺🇸' },
  { code: 'es', name: 'Spanish', flag: '🇪🇸' },
  { code: 'fr', name: 'French', flag: '🇫🇷' },
  { code: 'de', name: 'German', flag: '🇩🇪' },
  { code: 'ja', name: 'Japanese', flag: '🇯🇵' },
  { code: 'ko', name: 'Korean', flag: '🇰🇷' },
  { code: 'zh', name: 'Chinese', flag: '🇨🇳' },
  { code: 'pt', name: 'Portuguese', flag: '🇵🇹' },
  { code: 'ru', name: 'Russian', flag: '🇷🇺' },
  { code: 'it', name: 'Italian', flag: '🇮🇹' },
  { code: 'ar', name: 'Arabic', flag: '🇸🇦' },
  { code: 'hi', name: 'Hindi', flag: '🇮🇳' },
  { code: 'th', name: 'Thai', flag: '🇹🇭' },
  { code: 'vi', name: 'Vietnamese', flag: '🇻🇳' },
  { code: 'nl', name: 'Dutch', flag: '🇳🇱' },
  { code: 'sv', name: 'Swedish', flag: '🇸🇪' },
  { code: 'no', name: 'Norwegian', flag: '🇳🇴' },
  { code: 'da', name: 'Danish', flag: '🇩🇰' },
  { code: 'fi', name: 'Finnish', flag: '🇫🇮' },
  { code: 'pl', name: 'Polish', flag: '🇵🇱' },
] as const;

export const VIDEO_API_ENDPOINTS = {
  VIDEO_INFO: (videoUrl: string) => `/${encodeURIComponent(videoUrl)}/info`,
  TRANSCRIPT: (videoUrl: string, lang?: string) => 
    `/${encodeURIComponent(videoUrl)}/transcript${lang ? `?lang=${lang}` : ''}`,
  TRANSLATE: (videoUrl: string) => `/${encodeURIComponent(videoUrl)}/translate`,
  LANGUAGES: (videoUrl: string) => `/${encodeURIComponent(videoUrl)}/languages`,
  CAPABILITIES: (videoUrl: string) => `/${encodeURIComponent(videoUrl)}/capabilities`,
  PROVIDERS: '/providers',
  SUPPORTED_LANGUAGES: '/languages',
} as const;

export const VIDEO_API_CONFIG = {
  DEFAULT_TIMEOUT: 30000,
  MAX_RETRIES: 3,
  RETRY_DELAY: 1000,
  CACHE_TTL: 3600000, // 1 hour in milliseconds
} as const;

export const VIDEO_URL_PATTERNS = {
  YOUTUBE: [
    /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([a-zA-Z0-9_-]{11})/,
    /youtube\.com\/watch\?.*v=([a-zA-Z0-9_-]{11})/,
  ],
  VIMEO: [
    /vimeo\.com\/(\d+)/,
  ],
  TWITCH: [
    /twitch\.tv\/videos\/(\d+)/,
  ],
} as const;

export const VIDEO_ERROR_CODES = {
  NETWORK_ERROR: 'NETWORK_ERROR',
  UNAUTHORIZED: 'UNAUTHORIZED',
  NOT_FOUND: 'NOT_FOUND',
  INVALID_URL: 'INVALID_URL',
  UNSUPPORTED_PROVIDER: 'UNSUPPORTED_PROVIDER',
  NO_TRANSCRIPT: 'NO_TRANSCRIPT',
  TRANSLATION_FAILED: 'TRANSLATION_FAILED',
  RATE_LIMITED: 'RATE_LIMITED',
  UNKNOWN_ERROR: 'UNKNOWN_ERROR',
} as const;

export const EXTENSION_STORAGE_KEYS = {
  AUTO_TRANSLATE_ENABLED: 'autoTranslateEnabled',
  DEFAULT_TARGET_LANGUAGE: 'defaultTargetLanguage',
  API_BASE_URL: 'apiBaseUrl',
  AUTH_TOKEN: 'authToken',
  TRANSLATION_CACHE: 'translationCache',
  USER_PREFERENCES: 'userPreferences',
} as const;

export const EXTENSION_MESSAGES = {
  GET_PAGE_INFO: 'GET_PAGE_INFO',
  TOGGLE_TRANSLATION: 'TOGGLE_TRANSLATION',
  SET_TARGET_LANGUAGE: 'SET_TARGET_LANGUAGE',
  CLEAR_CACHE: 'CLEAR_CACHE',
  UPDATE_SETTINGS: 'UPDATE_SETTINGS',
} as const;