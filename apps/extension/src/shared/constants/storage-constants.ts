// Storage keys for extension settings and data
export const EXTENSION_STORAGE_KEYS = {
  API_BASE_URL: 'apiBaseUrl',
  AUTH_TOKEN: 'authToken',
  AUTO_TRANSLATE_ENABLED: 'autoTranslateEnabled',
  AUTO_FETCH_TRANSCRIPT: 'autoFetchTranscript',
  DUAL_LANGUAGE_ENABLED: 'dualLanguageEnabled',
  PRIMARY_LANGUAGE: 'primaryLanguage',
  SECONDARY_LANGUAGE: 'secondaryLanguage',
  TRANSCRIPT_OVERLAY_ENABLED: 'transcriptOverlayEnabled',
  TRANSCRIPT_CACHE: 'transcriptCache',
} as const;
