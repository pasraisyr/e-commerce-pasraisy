export type PaymentStatus = 'pending' | 'successful' | 'failed'

export interface PaymentInfo {
  billCode: string
  amount: number
  status: string
  paymentUrl: string
}
