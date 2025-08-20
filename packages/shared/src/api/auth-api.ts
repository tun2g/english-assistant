import { API_ENDPOINTS } from '../constants/api-constants';
import type {
  AuthResponse,
  LoginRequest,
  RefreshTokenRequest,
  RegisterRequest,
  UpdateProfileRequest,
  User,
} from '../types/auth-types';
import { apiGet, apiPatch, apiPost } from './axios-client';

export async function loginUser(credentials: LoginRequest): Promise<AuthResponse> {
  return apiPost<AuthResponse>(API_ENDPOINTS.AUTH.LOGIN, credentials);
}

export async function registerUser(userData: RegisterRequest): Promise<AuthResponse> {
  return apiPost<AuthResponse>(API_ENDPOINTS.AUTH.REGISTER, userData);
}

export async function logoutUser(): Promise<null> {
  return apiPost<null>(API_ENDPOINTS.AUTH.LOGOUT, {});
}

export async function refreshUserToken(tokenData: RefreshTokenRequest): Promise<AuthResponse> {
  return apiPost<AuthResponse>(API_ENDPOINTS.AUTH.REFRESH, tokenData);
}

export async function getUserProfile(): Promise<User> {
  return apiGet<User>(API_ENDPOINTS.AUTH.PROFILE);
}

export async function updateUserProfile(profileData: UpdateProfileRequest): Promise<User> {
  return apiPatch<User>(API_ENDPOINTS.USER.UPDATE, profileData);
}
