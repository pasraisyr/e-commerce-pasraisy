export interface Product {
  id: number
  sku: string
  name: string
  description: string
  price: number
  stock: number
  weight: number
  imageUrl: string
  isActive: boolean
}

export type ProductInput = Omit<Product, 'id'>
