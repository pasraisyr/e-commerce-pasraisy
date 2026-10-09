import { request } from './client'
import type { CreateOrderPayload, CreateOrderResponse, Order } from '../types/order'

export const ordersApi = {
  create: (payload: CreateOrderPayload) =>
    request<CreateOrderResponse>('/orders', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),
  track: (orderNumber: string) =>
    request<Order>(`/orders/track?orderNumber=${encodeURIComponent(orderNumber)}`),
  getAdminOrders: () => request<Order[]>('/admin/orders'),
  getMyOrders: () => request<Order[]>('/orders/my-orders'),
}
