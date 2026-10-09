import React from 'react'
import { useCart } from '../context/CartContext'

export const Toast: React.FC = () => {
  const { toastMsg } = useCart()

  if (!toastMsg) return null

  return (
    <div
      className="font-mono"
      style={{
        position: 'fixed',
        bottom: '24px',
        right: '24px',
        background: 'var(--border)',
        color: '#FFFFFF',
        padding: '12px 20px',
        fontSize: '12px',
        fontWeight: 600,
        border: '1px solid #000',
        boxShadow: '4px 4px 0px rgba(0,0,0,0.2)',
        zIndex: 9999,
      }}
    >
      {toastMsg}
    </div>
  )
}
