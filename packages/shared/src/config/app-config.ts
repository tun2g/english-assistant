// Check if running in Chrome extension context
const isExtension =
  typeof globalThis !== 'undefined' &&
  typeof globalThis.chrome !== 'undefined' &&
  globalThis.chrome.runtime &&
  globalThis.chrome.runtime.id;

export const APP_CONFIG = {
  NAME: 'English Learning Platform',
  VERSION: '1.0.0',
  DESCRIPTION: 'A comprehensive platform for learning English',
  AUTHOR: 'English Learning Team',
  BASE_URL:
    typeof window !== 'undefined' && !isExtension
      ? window.location.origin
      : process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000',
  API_BASE_URL: isExtension
    ? 'http://localhost:8000' // Extension always uses backend server
    : typeof window !== 'undefined'
      ? window.location.origin.replace('3000', '8000')
      : process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000',
} as const;

export const THEME_CONFIG = {
  DEFAULT_THEME: 'light',
  THEMES: {
    LIGHT: 'light',
    DARK: 'dark',
    SYSTEM: 'system',
  },
} as const;

export const PAGINATION_CONFIG = {
  DEFAULT_PAGE_SIZE: 10,
  MAX_PAGE_SIZE: 100,
  DEFAULT_PAGE: 1,
} as const;

export const VALIDATION_CONFIG = {
  PASSWORD_MIN_LENGTH: 8,
  USERNAME_MIN_LENGTH: 3,
  USERNAME_MAX_LENGTH: 50,
  EMAIL_MAX_LENGTH: 255,
} as const;

export const STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  REFRESH_TOKEN: 'refresh_token',
  USER_PREFERENCES: 'user_preferences',
  THEME: 'theme',
} as const;

export type Theme = (typeof THEME_CONFIG.THEMES)[keyof typeof THEME_CONFIG.THEMES];
export type StorageKey = (typeof STORAGE_KEYS)[keyof typeof STORAGE_KEYS];
