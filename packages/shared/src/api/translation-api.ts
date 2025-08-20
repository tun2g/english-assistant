import { API_ENDPOINTS } from '../constants/api-constants';
import type {
  GetSupportedLanguagesResponse,
  TranslateTextsRequest,
  TranslateTextsResponse,
} from '../types/translation-types';
import { apiGet, apiPost } from './axios-client';

export async function translateTexts(request: TranslateTextsRequest): Promise<TranslateTextsResponse> {
  return apiPost<TranslateTextsResponse>(API_ENDPOINTS.TRANSLATION.TRANSLATE_TEXTS, request);
}

export async function getSupportedTranslationLanguages(): Promise<GetSupportedLanguagesResponse> {
  return apiGet<GetSupportedLanguagesResponse>(API_ENDPOINTS.TRANSLATION.SUPPORTED_LANGUAGES);
}
