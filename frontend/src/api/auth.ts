import { request } from './client'
import type { AuthResponse, LoginPayload, RegisterPayload, User } from '../types/user'

export const authApi = {
  login: (payload: LoginPayload) =>
    request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),
  register: (payload: RegisterPayload) =>
    request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),
  getMe: () => request<User>('/auth/me'),
}
