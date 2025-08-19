import { apiGet, apiPost, apiPatch } from './axios-client';
import { API_ENDPOINTS } from '../constants/api-constants';
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  RefreshTokenRequest,
  UpdateProfileRequest,
  User,
} from '../types/auth-types';
import type { ApiResponse } from '../types/api-types';

export async function loginUser(credentials: LoginRequest): Promise<ApiResponse<AuthResponse>> {
  return apiPost<AuthResponse>(API_ENDPOINTS.AUTH.LOGIN, credentials);
}

export async function registerUser(userData: RegisterRequest): Promise<ApiResponse<AuthResponse>> {
  return apiPost<AuthResponse>(API_ENDPOINTS.AUTH.REGISTER, userData);
}

export async function logoutUser(): Promise<ApiResponse<null>> {
  return apiPost<null>(API_ENDPOINTS.AUTH.LOGOUT, {});
}

export async function refreshUserToken(tokenData: RefreshTokenRequest): Promise<ApiResponse<AuthResponse>> {
  return apiPost<AuthResponse>(API_ENDPOINTS.AUTH.REFRESH, tokenData);
}

export async function getUserProfile(): Promise<ApiResponse<User>> {
  return apiGet<User>(API_ENDPOINTS.AUTH.PROFILE);
}

export async function updateUserProfile(profileData: UpdateProfileRequest): Promise<ApiResponse<User>> {
  return apiPatch<User>(API_ENDPOINTS.USER.UPDATE, profileData);
}