import { STORAGE_KEYS } from '../config/app-config';
import type { AuthResponse } from '../types/auth-types';
import { getStorageItem, setStorageItem, removeStorageItem } from '../storage/universal-storage';

// Token state
let accessToken: string | null = null;
let refreshToken: string | null = null;

export async function initializeTokens(): Promise<void> {
  try {
    accessToken = await getStorageItem<string>(STORAGE_KEYS.AUTH_TOKEN);
    refreshToken = await getStorageItem<string>(STORAGE_KEYS.REFRESH_TOKEN);
  } catch (error) {
    console.error('Failed to initialize tokens:', error);
    accessToken = null;
    refreshToken = null;
  }
}

export async function saveTokens(authResponse: AuthResponse): Promise<void> {
  try {
    accessToken = authResponse.accessToken;
    refreshToken = authResponse.refreshToken;
    
    await Promise.all([
      setStorageItem(STORAGE_KEYS.AUTH_TOKEN, authResponse.accessToken),
      setStorageItem(STORAGE_KEYS.REFRESH_TOKEN, authResponse.refreshToken),
    ]);
  } catch (error) {
    console.error('Failed to save tokens:', error);
    throw error;
  }
}

export async function clearTokens(): Promise<void> {
  try {
    accessToken = null;
    refreshToken = null;
    
    await Promise.all([
      removeStorageItem(STORAGE_KEYS.AUTH_TOKEN),
      removeStorageItem(STORAGE_KEYS.REFRESH_TOKEN),
    ]);
  } catch (error) {
    console.error('Failed to clear tokens:', error);
    throw error;
  }
}

export function getAccessToken(): string | null {
  return accessToken;
}

export function getRefreshToken(): string | null {
  return refreshToken;
}

export function hasValidToken(): boolean {
  return !!accessToken;
}

// Initialize tokens on module load
initializeTokens();