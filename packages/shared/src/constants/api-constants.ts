export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/api/v1/auth/login',
    LOGOUT: '/api/v1/auth/logout',
    REGISTER: '/api/v1/auth/register',
    REFRESH: '/api/v1/auth/refresh',
    PROFILE: '/api/v1/auth/profile',
  },
  USER: {
    PROFILE: '/api/v1/user/profile',
    UPDATE: '/api/v1/user/update',
    DELETE: '/api/v1/user/delete',
  },
  VIDEO: {
    PROVIDERS: '/api/v1/video/providers',
    LANGUAGES: '/api/v1/video/languages',
    INFO: (videoId: string) => `/api/v1/video/${videoId}/info`,
    TRANSCRIPT: (videoId: string) => `/api/v1/video/${videoId}/transcript`,
    AVAILABLE_LANGUAGES: (videoId: string) => `/api/v1/video/${videoId}/languages`,
    CAPABILITIES: (videoId: string) => `/api/v1/video/${videoId}/capabilities`,
  },
  TRANSLATION: {
    TRANSLATE_TEXTS: '/api/v1/translate',
    SUPPORTED_LANGUAGES: '/api/v1/translate/languages',
  },
} as const;

export const HTTP_METHODS = {
  GET: 'GET',
  POST: 'POST',
  PUT: 'PUT',
  PATCH: 'PATCH',
  DELETE: 'DELETE',
} as const;

export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500,
} as const;

export const API_CONFIG = {
  TIMEOUT: 30000,
  RETRY_ATTEMPTS: 3,
  RETRY_DELAY: 1000,
} as const;

export type HttpMethod = (typeof HTTP_METHODS)[keyof typeof HTTP_METHODS];
export type HttpStatus = (typeof HTTP_STATUS)[keyof typeof HTTP_STATUS];
export type ApiEndpoint =
  (typeof API_ENDPOINTS)[keyof typeof API_ENDPOINTS][keyof (typeof API_ENDPOINTS)[keyof typeof API_ENDPOINTS]];
