import type { AxiosInstance } from 'axios';
import { createAxiosClient } from './axios-client';
import type {
  VideoInfo,
  VideoCapabilities,
  Language,
  TranscriptResponse,
  TranslationResponse,
  GetVideoInfoRequest,
  GetTranscriptRequest,
  TranslateTranscriptRequest,
  GetAvailableLanguagesRequest,
  GetCapabilitiesRequest,
  GetAvailableLanguagesResponse,
  GetSupportedProvidersResponse,
  GetSupportedLanguagesResponse,
  VideoApiConfig,
  VideoApiError,
  VideoProvider,
} from '../types/video-types';

/**
 * Video API client for backend communication
 * Handles video information, transcripts, and translations
 */
export class VideoApiClient {
  private client: AxiosInstance;
  private config: VideoApiConfig;

  constructor(config: VideoApiConfig) {
    this.config = config;
    this.client = createAxiosClient({
      baseURL: `${config.baseUrl}/video`,
      timeout: config.timeout || 30000,
    });

    // Set up request interceptor for authentication
    this.client.interceptors.request.use((requestConfig) => {
      if (this.config.authToken) {
        requestConfig.headers.Authorization = `Bearer ${this.config.authToken}`;
      }
      return requestConfig;
    });

    // Set up response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        throw this.createVideoApiError(error);
      }
    );
  }

  /**
   * Update authentication token
   */
  setAuthToken(token: string) {
    this.config.authToken = token;
  }

  /**
   * Get video information
   */
  async getVideoInfo(request: GetVideoInfoRequest): Promise<VideoInfo> {
    const videoId = this.extractVideoIdFromUrl(request.videoUrl);
    if (!videoId) {
      throw new Error('Unable to extract video ID from URL: ' + request.videoUrl);
    }
    const response = await this.client.get<VideoInfo>(`/${videoId}/info`);
    return response.data;
  }

  /**
   * Get video transcript
   */
  async getTranscript(request: GetTranscriptRequest): Promise<TranscriptResponse> {
    const videoId = this.extractVideoIdFromUrl(request.videoUrl);
    if (!videoId) {
      throw new Error('Unable to extract video ID from URL: ' + request.videoUrl);
    }
    const queryParam = request.language ? `?lang=${request.language}` : '';
    const response = await this.client.get<TranscriptResponse>(`/${videoId}/transcript${queryParam}`);
    return response.data;
  }

  /**
   * Translate video transcript
   */
  async translateTranscript(request: TranslateTranscriptRequest): Promise<TranslationResponse> {
    const videoId = this.extractVideoIdFromUrl(request.videoUrl);
    if (!videoId) {
      throw new Error('Unable to extract video ID from URL: ' + request.videoUrl);
    }
    const response = await this.client.post<TranslationResponse>(`/${videoId}/translate`, {
      targetLang: request.targetLang,
      sourceLang: request.sourceLang,
      cacheResult: request.cacheResult ?? true,
    });
    return response.data;
  }

  /**
   * Get available transcript languages for a video
   */
  async getAvailableLanguages(request: GetAvailableLanguagesRequest): Promise<Language[]> {
    const videoId = this.extractVideoIdFromUrl(request.videoUrl);
    if (!videoId) {
      throw new Error('Unable to extract video ID from URL: ' + request.videoUrl);
    }
    const response = await this.client.get<GetAvailableLanguagesResponse>(`/${videoId}/languages`);
    return response.data.languages;
  }

  /**
   * Get video capabilities
   */
  async getCapabilities(request: GetCapabilitiesRequest): Promise<VideoCapabilities> {
    const videoId = this.extractVideoIdFromUrl(request.videoUrl);
    if (!videoId) {
      throw new Error('Unable to extract video ID from URL: ' + request.videoUrl);
    }
    const response = await this.client.get<VideoCapabilities>(`/${videoId}/capabilities`);
    return response.data;
  }

  /**
   * Get supported video providers
   */
  async getSupportedProviders(): Promise<VideoProvider[]> {
    const response = await this.client.get<GetSupportedProvidersResponse>('/providers');
    return response.data.providers;
  }

  /**
   * Get supported translation languages
   */
  async getSupportedLanguages(): Promise<Language[]> {
    const response = await this.client.get<GetSupportedLanguagesResponse>('/languages');
    return response.data.languages;
  }

  /**
   * Extract video ID from URL (instance method)
   */
  private extractVideoIdFromUrl(url: string): string {
    return VideoApiClient.extractVideoId(url);
  }

  /**
   * Utility method to extract video ID from URL
   */
  static extractVideoId(url: string): string {
    // Handle different YouTube URL formats
    const patterns = [
      /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([a-zA-Z0-9_-]{11})/,
    ];

    for (const pattern of patterns) {
      const match = url.match(pattern);
      if (match) {
        return match[1];
      }
    }

    // If no pattern matches, assume it's already a video ID
    return url;
  }

  /**
   * Check if URL is a supported video URL
   */
  static isSupportedVideoUrl(url: string): boolean {
    return url.includes('youtube.com') || 
           url.includes('youtu.be') || 
           url.includes('youtube-nocookie.com');
  }

  /**
   * Validate video ID format
   */
  static validateVideoId(videoId: string, provider: VideoProvider = 'youtube'): boolean {
    switch (provider) {
      case 'youtube':
        // YouTube video IDs are 11 characters long
        return /^[a-zA-Z0-9_-]{11}$/.test(videoId);
      default:
        return false;
    }
  }

  /**
   * Create a standardized API error
   */
  private createVideoApiError(error: any): VideoApiError {
    const apiError = new Error() as VideoApiError;
    
    if (error.response) {
      // Server responded with error status
      const data = error.response.data;
      apiError.message = data?.error || `HTTP ${error.response.status}: ${error.response.statusText}`;
      apiError.code = data?.code;
      apiError.details = data?.details;
      apiError.status = error.response.status;
    } else if (error.request) {
      // Request was made but no response received
      apiError.message = 'Network error: No response from server';
      apiError.code = 'NETWORK_ERROR';
    } else {
      // Something else happened
      apiError.message = error.message || 'Unknown error occurred';
      apiError.code = 'UNKNOWN_ERROR';
    }

    return apiError;
  }
}

/**
 * Factory function to create video API client
 */
export function createVideoApiClient(config: VideoApiConfig): VideoApiClient {
  return new VideoApiClient(config);
}

/**
 * Default video API client instance
 * Can be configured with setConfig()
 */
let defaultClient: VideoApiClient | null = null;

export function getVideoApiClient(): VideoApiClient {
  if (!defaultClient) {
    throw new Error('Video API client not configured. Call setVideoApiConfig() first.');
  }
  return defaultClient;
}

export function setVideoApiConfig(config: VideoApiConfig): void {
  defaultClient = new VideoApiClient(config);
}

// Re-export types for convenience
export type {
  VideoInfo,
  VideoCapabilities,
  Language,
  TranscriptResponse,
  TranslationResponse,
  VideoApiConfig,
  VideoApiError,
  VideoProvider,
} from '../types/video-types';