import React from 'react'
import type { Product } from '../types'
import { useCart } from '../context/CartContext'

interface ProductCardProps {
  product: Product
}

export const ProductCard: React.FC<ProductCardProps> = ({ product }) => {
  const { addToCart } = useCart()

  return (
    <div className="product-box">
      {/* Header Box Petak */}
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          padding: '8px 12px',
          background: 'var(--surface-alt)',
          borderBottom: '1px solid var(--border)',
          fontSize: '11px',
        }}
        className="font-mono"
      >
        <span>SKU: {product.sku || `PRD-${product.id}`}</span>
        <span>BERAT: {product.weight ? `${product.weight.toFixed(3)} KG` : '0.500 KG'}</span>
      </div>

      {/* Gambar */}
      <div className="img-wrapper">
        <img
          src={product.imageUrl || 'https://placehold.co/400x400/png?text=Tiada+Gambar'}
          alt={product.name}
        />
      </div>

      {/* Kandungan Butiran */}
      <div style={{ padding: '14px', flex: 1, display: 'flex', flexDirection: 'column' }}>
        <h3 style={{ fontSize: '15px', fontWeight: 600, marginBottom: '6px' }}>{product.name}</h3>
        <p style={{ fontSize: '12px', color: 'var(--text-muted)', marginBottom: '14px', flex: 1 }}>
          {product.description || 'Tiada keterangan tambahan.'}
        </p>

        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'baseline',
            paddingTop: '10px',
            borderTop: '1px solid var(--border-muted)',
            marginBottom: '12px',
          }}
          className="font-mono"
        >
          <div>
            <span style={{ fontSize: '11px', color: 'var(--text-muted)' }}>HARGA: </span>
            <strong style={{ fontSize: '16px' }}>RM {product.price.toFixed(2)}</strong>
          </div>
          <span className="badge-raw">STOK: {product.stock}</span>
        </div>

        <button
          onClick={() => addToCart(product)}
          disabled={product.stock <= 0}
          className="btn-raw primary"
          style={{ width: '100%' }}
        >
          {product.stock > 0 ? '+ TAMBAH KE TROLI' : 'STOK HABIS'}
        </button>
      </div>
    </div>
  )
}
