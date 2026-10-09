import { request } from './client'
import type { CourierRate } from '../types/shipping'

export const shippingApi = {
  checkRates: (receiverPostcode: string, weightKg: number) =>
    request<CourierRate[]>('/shipping/rates', {
      method: 'POST',
      body: JSON.stringify({ receiverPostcode, weightKg }),
    }),
}
