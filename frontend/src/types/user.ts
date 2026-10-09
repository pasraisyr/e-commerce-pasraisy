export type UserRole = 'admin' | 'customer'

export interface User {
  id: number
  name: string
  email: string
  phone?: string
  role: UserRole
  address?: string
  city?: string
  state?: string
  postcode?: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface LoginPayload {
  email: string
  password: string
}

export interface RegisterPayload {
  name: string
  email: string
  password: string
  phone?: string
  address?: string
  city?: string
  state?: string
  postcode?: string
}
