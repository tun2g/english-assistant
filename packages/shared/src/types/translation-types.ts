// Translation API types for backend integration

export interface TranslateTextsRequest {
  texts: string[];
  targetLang: string;
  sourceLang?: string;
}

export interface TranslateTextsResponse {
  translations: string[];
  sourceLang: string;
  targetLang: string;
}

export interface TranslationApiConfig {
  baseUrl: string;
  timeout?: number;
  retries?: number;
  authToken?: string;
}

export interface TranslationApiError extends Error {
  code?: string;
  details?: string;
  status?: number;
}

// Language detection types
export interface Language {
  code: string;
  name: string;
}

export interface GetSupportedLanguagesResponse {
  languages: Language[];
}
