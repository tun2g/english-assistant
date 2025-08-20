// Notification types and timing
export const NOTIFICATION_TYPES = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  LOADING: 'loading',
} as const;

export const NOTIFICATION_DURATIONS = {
  SUCCESS: 3000,
  ERROR: 5000,
  WARNING: 4000,
  LOADING: 0, // Don't auto-hide loading notifications
} as const;
