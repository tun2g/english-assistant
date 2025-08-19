export const QUERY_KEY = {
  AUTH: {
    LOGIN: 'auth.login',
    LOGOUT: 'auth.logout',
    PROFILE: 'auth.profile',
    REFRESH: 'auth.refresh',
  },
  USER: {
    PROFILE: 'user.profile',
    LIST: 'user.list',
    DETAIL: 'user.detail',
  },
  LEARNING: {
    LESSONS: 'learning.lessons',
    PROGRESS: 'learning.progress',
    VOCABULARY: 'learning.vocabulary',
  },
} as const;

export type QueryKey = typeof QUERY_KEY[keyof typeof QUERY_KEY][keyof typeof QUERY_KEY[keyof typeof QUERY_KEY]];