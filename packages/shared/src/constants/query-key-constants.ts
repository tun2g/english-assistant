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
  VIDEO: {
    INFO: 'video.info',
    TRANSCRIPT: 'video.transcript',
    PROGRESS: 'video.progress',
    CAPABILITIES: 'video.capabilities',
    LANGUAGES: 'video.languages',
  },
  OAUTH: {
    STATUS: 'oauth.status',
    TOKEN: 'oauth.token',
  },
  PAGE_INFO: 'page.info',
  AUTO_TRANSLATE: 'auto.translate',
  YOUTUBE: {
    VIDEO_INFO: 'youtube.video.info',
    TRANSCRIPT: 'youtube.transcript',
  },
} as const;

export type QueryKey = (typeof QUERY_KEY)[keyof typeof QUERY_KEY][keyof (typeof QUERY_KEY)[keyof typeof QUERY_KEY]];
