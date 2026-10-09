import React from 'react'

export const Footer: React.FC = () => {
  return (
    <footer
      style={{
        marginTop: '40px',
        padding: '16px 20px',
        border: '1px solid var(--border)',
        background: 'var(--surface)',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        flexWrap: 'wrap',
        gap: '12px',
      }}
      className="font-mono"
    >
      <span style={{ fontSize: '11px', color: 'var(--text-muted)' }}>
        © PAS RAISY E-COMMERCE • MODAL SISTEM: RAW TABLE & BOX MINIMALIST
      </span>
      <div style={{ display: 'flex', gap: '12px', fontSize: '11px' }}>
        <span>TOYYIBPAY FPX</span>
        <span>•</span>
        <span>EASYPARCEL SHIPPING</span>
        <span>•</span>
        <span>POSTGRESQL & GORM</span>
      </div>
    </footer>
  )
}
