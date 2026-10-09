import type { PaymentInfo } from './payment'
import type { ShipmentInfo } from './shipping'

export type OrderStatus =
  | 'pending_payment'
  | 'paid'
  | 'processing'
  | 'shipped'
  | 'completed'
  | 'cancelled'

export interface OrderItem {
  id: number
  productName: string
  price: number
  quantity: number
  subtotal: number
  weight: number
}

export interface Order {
  id: number
  orderNumber: string
  customerName: string
  customerEmail: string
  customerPhone: string
  shippingAddress: string
  city: string
  state: string
  postcode: string
  subtotal: number
  shippingFee: number
  totalAmount: number
  totalWeight: number
  status: string
  notes?: string
  items: OrderItem[]
  payment?: PaymentInfo
  shipment?: ShipmentInfo
  createdAt: string
}

export interface CreateOrderPayload {
  customerName: string
  customerEmail: string
  customerPhone: string
  shippingAddress: string
  city: string
  state: string
  postcode: string
  courierName: string
  shippingRate: number
  notes?: string
  items: Array<{ productId: number; quantity: number }>
}

export interface CreateOrderResponse {
  order: Order
  paymentUrl: string
}
