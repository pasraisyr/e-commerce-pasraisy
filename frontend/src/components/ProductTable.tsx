import React from 'react'
import type { Product } from '../types'
import { useCart } from '../context/CartContext'

interface ProductTableProps {
  products: Product[]
}

export const ProductTable: React.FC<ProductTableProps> = ({ products }) => {
  const { addToCart } = useCart()

  return (
    <div style={{ overflowX: 'auto' }}>
      <table className="raw-table">
        <thead>
          <tr>
            <th style={{ width: '60px' }}>NO</th>
            <th style={{ width: '110px' }}>SKU</th>
            <th style={{ width: '80px' }}>GAMBAR</th>
            <th>NAMA PRODUK & KETERANGAN</th>
            <th style={{ width: '120px' }}>BERAT (KG)</th>
            <th style={{ width: '100px' }}>STOK</th>
            <th style={{ width: '120px' }}>HARGA</th>
            <th style={{ width: '150px' }}>TINDAKAN</th>
          </tr>
        </thead>
        <tbody>
          {products.map((p, idx) => (
            <tr key={p.id}>
              <td className="font-mono">{idx + 1}</td>
              <td className="font-mono" style={{ fontWeight: 600 }}>
                {p.sku || `PRD-${p.id}`}
              </td>
              <td>
                <img
                  src={p.imageUrl}
                  alt=""
                  style={{
                    width: '48px',
                    height: '48px',
                    objectFit: 'cover',
                    border: '1px solid var(--border)',
                  }}
                />
              </td>
              <td>
                <div style={{ fontWeight: 600, fontSize: '14px' }}>{p.name}</div>
                <div style={{ fontSize: '12px', color: 'var(--text-muted)' }}>{p.description}</div>
              </td>
              <td className="font-mono">{(p.weight || 0.5).toFixed(3)} KG</td>
              <td>
                <span className="badge-raw">{p.stock} UNIT</span>
              </td>
              <td className="font-mono" style={{ fontWeight: 700, fontSize: '14px' }}>
                RM {p.price.toFixed(2)}
              </td>
              <td>
                <button
                  onClick={() => addToCart(p)}
                  disabled={p.stock <= 0}
                  className="btn-raw sm primary"
                >
                  + TROLI
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
