import { request } from './client'

export interface HealthResponse {
  status: string
  database: string
  timestamp?: string
  service?: string
}

export const healthApi = {
  check: () => request<HealthResponse>('/health'),
}
