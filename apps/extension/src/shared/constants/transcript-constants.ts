// Transcript overlay configuration
export const TRANSCRIPT_CONFIG = {
  SYNC_INTERVAL: 100, // ms - how often to check video time
  TRANSLATION_QUEUE_SIZE: 10, // Number of segments to translate ahead
  CACHE_EXPIRY: 3600000, // 1 hour in ms
  SEGMENT_HIGHLIGHT_OFFSET: 100, // ms - reduced offset for faster real-time sync
} as const;
