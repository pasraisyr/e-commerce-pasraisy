import { request } from './client'
import type { Product } from '../types/product'

export const productsApi = {
  getAll: () => request<Product[]>('/products'),
  getById: (id: number) => request<Product>(`/products/${id}`),
  create: (data: Partial<Product>) =>
    request<Product>('/admin/products', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: number, data: Partial<Product>) =>
    request<Product>(`/admin/products/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: number) =>
    request<{ message: string }>(`/admin/products/${id}`, {
      method: 'DELETE',
    }),
}
