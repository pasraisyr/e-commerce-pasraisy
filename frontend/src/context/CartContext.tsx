import React, { createContext, useContext, useState, type ReactNode } from 'react'
import type { CartItem, Product } from '../types'

interface CartContextType {
  cart: CartItem[]
  addToCart: (product: Product) => void
  updateQty: (productId: number, delta: number) => void
  removeFromCart: (productId: number) => void
  clearCart: () => void
  totalItems: number
  subtotal: number
  totalWeight: number
  toastMsg: string
  showToast: (msg: string) => void
}

const CartContext = createContext<CartContextType | undefined>(undefined)

export const CartProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [cart, setCart] = useState<CartItem[]>([])
  const [toastMsg, setToastMsg] = useState<string>('')

  const showToast = (msg: string) => {
    setToastMsg(msg)
    setTimeout(() => setToastMsg(''), 3500)
  }

  const addToCart = (product: Product) => {
    setCart((prev) => {
      const existing = prev.find((item) => item.product.id === product.id)
      if (existing) {
        if (existing.quantity >= product.stock) {
          showToast(`[!] Had stok maksimum tercapai (${product.stock} unit)`)
          return prev
        }
        return prev.map((item) =>
          item.product.id === product.id ? { ...item, quantity: item.quantity + 1 } : item
        )
      }
      return [...prev, { product, quantity: 1 }]
    })
    showToast(`[✓] ${product.name.toUpperCase()} DITAMBAH KE TROLI`)
  }

  const updateQty = (productId: number, delta: number) => {
    setCart((prev) =>
      prev
        .map((item) => {
          if (item.product.id === productId) {
            const nextQty = item.quantity + delta
            if (nextQty <= 0) return null
            if (nextQty > item.product.stock) {
              showToast(`[!] Melebihi baki stok yang ada`)
              return item
            }
            return { ...item, quantity: nextQty }
          }
          return item
        })
        .filter(Boolean) as CartItem[]
    )
  }

  const removeFromCart = (productId: number) => {
    setCart((prev) => prev.filter((item) => item.product.id !== productId))
  }

  const clearCart = () => setCart([])

  const totalItems = cart.reduce((sum, item) => sum + item.quantity, 0)
  const subtotal = cart.reduce((sum, item) => sum + item.product.price * item.quantity, 0)
  const totalWeight = cart.reduce(
    (sum, item) => sum + (item.product.weight || 0.5) * item.quantity,
    0
  )

  return (
    <CartContext.Provider
      value={{
        cart,
        addToCart,
        updateQty,
        removeFromCart,
        clearCart,
        totalItems,
        subtotal,
        totalWeight,
        toastMsg,
        showToast,
      }}
    >
      {children}
    </CartContext.Provider>
  )
}

export const useCart = (): CartContextType => {
  const context = useContext(CartContext)
  if (!context) {
    throw new Error('useCart must be used within a CartProvider')
  }
  return context
}
