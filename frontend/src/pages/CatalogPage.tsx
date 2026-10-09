import React, { useState } from 'react'
import type { Product, ViewMode } from '../types'
import { ProductCard } from '../components/ProductCard'
import { ProductTable } from '../components/ProductTable'

interface CatalogPageProps {
  products: Product[]
  loading: boolean
}

export const CatalogPage: React.FC<CatalogPageProps> = ({ products, loading }) => {
  const [viewMode, setViewMode] = useState<ViewMode>('grid')
  const [searchTerm, setSearchTerm] = useState<string>('')

  const filteredProducts = products.filter(
    (p) =>
      p.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (p.sku && p.sku.toLowerCase().includes(searchTerm.toLowerCase()))
  )

  return (
    <div>
      {/* Subheader Toolbar */}
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          background: 'var(--surface)',
          border: '1px solid var(--border)',
          padding: '10px 16px',
          marginBottom: '16px',
          flexWrap: 'wrap',
          gap: '12px',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px', flexWrap: 'wrap' }}>
          <span className="uppercase-label">MOD PAPARAN:</span>
          <button
            onClick={() => setViewMode('grid')}
            className={`btn-raw sm ${viewMode === 'grid' ? 'primary' : ''}`}
          >
            ■ JADUAL PETAK (BOX)
          </button>
          <button
            onClick={() => setViewMode('table')}
            className={`btn-raw sm ${viewMode === 'table' ? 'primary' : ''}`}
          >
            ☰ SENARAI RAW (TABLE)
          </button>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <input
            type="text"
            className="input-raw font-mono"
            placeholder="Cari produk / SKU..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            style={{ padding: '4px 8px', fontSize: '12px', maxWidth: '180px' }}
          />
          <div className="font-mono" style={{ fontSize: '12px', color: 'var(--text-muted)' }}>
            TERSEDIA: {filteredProducts.length} UNIT
          </div>
        </div>
      </div>

      {loading ? (
        <div
          className="font-mono"
          style={{
            background: 'var(--surface)',
            border: '1px solid var(--border)',
            padding: '40px',
            textAlign: 'center',
          }}
        >
          MEMUAT DATA INVENTORI DARI POSTGRESQL...
        </div>
      ) : filteredProducts.length === 0 ? (
        <div
          className="font-mono"
          style={{
            background: 'var(--surface)',
            border: '1px solid var(--border)',
            padding: '40px',
            textAlign: 'center',
          }}
        >
          TIADA PRODUK DITEMUI.
        </div>
      ) : viewMode === 'grid' ? (
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
            gap: '16px',
          }}
        >
          {filteredProducts.map((p) => (
            <ProductCard key={p.id} product={p} />
          ))}
        </div>
      ) : (
        <ProductTable products={filteredProducts} />
      )}
    </div>
  )
}
