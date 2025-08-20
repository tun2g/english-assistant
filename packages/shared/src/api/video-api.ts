import { apiGet } from './axios-client';
import { API_ENDPOINTS } from '../constants/api-constants';
import type {
  VideoInfo,
  VideoCapabilities,
  TranscriptResponse,
  GetVideoInfoRequest,
  GetTranscriptRequest,
  GetAvailableLanguagesRequest,
  GetCapabilitiesRequest,
  GetAvailableLanguagesResponse,
  GetSupportedProvidersResponse,
  GetSupportedLanguagesResponse,
} from '../types/video-types';
import type {} from '../types/api-types';
import { extractVideoId } from '../utils';

export async function getVideoInfo(request: GetVideoInfoRequest): Promise<VideoInfo> {
  const videoId = extractVideoId(request.videoUrl);
  return apiGet<VideoInfo>(API_ENDPOINTS.VIDEO.INFO(videoId));
}

export async function getVideoTranscript(request: GetTranscriptRequest): Promise<TranscriptResponse> {
  const videoId = extractVideoId(request.videoUrl);
  const queryParam = request.language ? `?lang=${request.language}` : '';
  return apiGet<TranscriptResponse>(`${API_ENDPOINTS.VIDEO.TRANSCRIPT(videoId)}${queryParam}`);
}

export async function getAvailableLanguages(
  request: GetAvailableLanguagesRequest
): Promise<GetAvailableLanguagesResponse> {
  const videoId = extractVideoId(request.videoUrl);
  return apiGet<GetAvailableLanguagesResponse>(API_ENDPOINTS.VIDEO.AVAILABLE_LANGUAGES(videoId));
}

export async function getVideoCapabilities(request: GetCapabilitiesRequest): Promise<VideoCapabilities> {
  const videoId = extractVideoId(request.videoUrl);
  return apiGet<VideoCapabilities>(API_ENDPOINTS.VIDEO.CAPABILITIES(videoId));
}

export async function getSupportedVideoProviders(): Promise<GetSupportedProvidersResponse> {
  return apiGet<GetSupportedProvidersResponse>(API_ENDPOINTS.VIDEO.PROVIDERS);
}

export async function getSupportedVideoLanguages(): Promise<GetSupportedLanguagesResponse> {
  return apiGet<GetSupportedLanguagesResponse>(API_ENDPOINTS.VIDEO.LANGUAGES);
}

// Additional functions for hooks
export async function saveVideoProgress(data: import('../types').VideoProgressRequest): Promise<void> {
  // This would typically make an API call to save progress
  // For now, just a placeholder
  console.log('Saving video progress:', data);
}
