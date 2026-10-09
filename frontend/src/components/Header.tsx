import React, { useEffect, useState } from 'react'
import { healthApi } from '../api'
import type { TabType } from '../types'
import { useCart } from '../context/CartContext'

interface HeaderProps {
  activeTab: TabType
  setActiveTab: (tab: TabType) => void
  productCount: number
}

export const Header: React.FC<HeaderProps> = ({ activeTab, setActiveTab, productCount }) => {
  const { totalItems } = useCart()
  const [serverStatus, setServerStatus] = useState<'online' | 'offline' | 'checking'>('checking')
  const isAdmin = Boolean(localStorage.getItem('admin_token'))

  useEffect(() => {
    const checkServer = async () => {
      try {
        const res = await healthApi.check()
        if (res.status === 'OK') setServerStatus('online')
        else setServerStatus('offline')
      } catch {
        setServerStatus('offline')
      }
    }
    checkServer()
  }, [])

  return (
    <header
      style={{
        border: '1px solid var(--border)',
        background: 'var(--surface)',
        padding: '16px 20px',
        marginBottom: '20px',
      }}
    >
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          flexWrap: 'wrap',
          gap: '16px',
        }}
      >
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <span className="font-mono" style={{ fontSize: '20px', fontWeight: 700, letterSpacing: '-0.02em' }}>
              PAS RAISY / KEDAI RUNCIT
            </span>
            <span
              className="badge-raw font-mono"
              style={{
                background: serverStatus === 'online' ? '#E3EDE5' : '#FCE8E6',
                color: serverStatus === 'online' ? '#1B4723' : '#8A1F1D',
                borderColor: serverStatus === 'online' ? '#7CA986' : '#E29795',
              }}
            >
              ● SERVER: {serverStatus.toUpperCase()}
            </span>
          </div>
          <p className="font-mono" style={{ fontSize: '11px', color: 'var(--text-muted)', marginTop: '4px' }}>
            INVENTORI & PEMBAYARAN RAW • INTEGRASI TOYYIBPAY (FPX) & EASYPARCEL
          </p>
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span className="badge-raw">MATA WANG: MYR (RM)</span>
          {isAdmin && (
            <span className="badge-raw" style={{ background: '#1A1A1A', color: '#FFF' }}>
              MOD ADMIN AKTIF
            </span>
          )}
        </div>
      </div>

      {/* Navigation Tabs */}
      <nav
        style={{
          display: 'flex',
          gap: '8px',
          marginTop: '16px',
          paddingTop: '14px',
          borderTop: '1px solid var(--border-muted)',
          flexWrap: 'wrap',
        }}
      >
        <button
          onClick={() => setActiveTab('katalog')}
          className={`btn-raw ${activeTab === 'katalog' ? 'primary' : ''}`}
        >
          [01] KATALOG PRODUK ({productCount})
        </button>
        <button
          onClick={() => setActiveTab('troli')}
          className={`btn-raw ${activeTab === 'troli' ? 'primary' : ''}`}
        >
          [02] TROLI BELIAN [{totalItems}]
        </button>
        <button
          onClick={() => setActiveTab('jejak')}
          className={`btn-raw ${activeTab === 'jejak' ? 'primary' : ''}`}
        >
          [03] JEJAK PESANAN
        </button>
        <button
          onClick={() => setActiveTab('admin')}
          className={`btn-raw ${activeTab === 'admin' ? 'primary' : ''}`}
        >
          [04] {isAdmin ? 'PANEL ADMIN' : 'LOGIN ADMIN'}
        </button>
      </nav>
    </header>
  )
}
