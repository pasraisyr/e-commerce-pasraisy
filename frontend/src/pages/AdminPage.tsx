import React, { useState, useEffect } from 'react'
import type { Order } from '../types'
import { authApi, ordersApi, productsApi } from '../api'
import { useCart } from '../context/CartContext'

interface AdminPageProps {
  onProductCreated: () => void
}

export const AdminPage: React.FC<AdminPageProps> = ({ onProductCreated }) => {
  const { showToast } = useCart()
  const [token, setToken] = useState<string>(() => localStorage.getItem('admin_token') || '')

  // Login Form
  const [email, setEmail] = useState<string>('admin@ecommerce.com')
  const [password, setPassword] = useState<string>('admin123')
  const [loggingIn, setLoggingIn] = useState<boolean>(false)

  // Admin Dashboard Tabs
  const [adminTab, setAdminTab] = useState<'orders' | 'add_product'>('orders')
  const [orders, setOrders] = useState<Order[]>([])
  const [loadingOrders, setLoadingOrders] = useState<boolean>(false)

  // Form Tambah Produk
  const [name, setName] = useState<string>('')
  const [sku, setSku] = useState<string>('')
  const [desc, setDesc] = useState<string>('')
  const [price, setPrice] = useState<string>('')
  const [stock, setStock] = useState<string>('20')
  const [weight, setWeight] = useState<string>('0.300')
  const [imageUrl, setImageUrl] = useState<string>(
    'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=500'
  )
  const [submittingProduct, setSubmittingProduct] = useState<boolean>(false)

  const fetchOrders = async () => {
    if (!token) return
    setLoadingOrders(true)
    try {
      const data = await ordersApi.getAdminOrders()
      setOrders(Array.isArray(data) ? data : [])
    } catch {
      showToast('[!] Gagal memuat pesanan')
    } finally {
      setLoadingOrders(false)
    }
  }

  useEffect(() => {
    if (token) {
      fetchOrders()
    }
  }, [token])

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoggingIn(true)
    try {
      const res = await authApi.login({ email, password })
      if (res.token) {
        setToken(res.token)
        localStorage.setItem('admin_token', res.token)
        showToast('[✓] Log masuk Admin berjaya!')
      }
    } catch (err: any) {
      showToast(`[!] ${err.message || 'Log masuk gagal'}`)
    } finally {
      setLoggingIn(false)
    }
  }

  const handleLogout = () => {
    setToken('')
    localStorage.removeItem('admin_token')
    showToast('[i] Admin telah log keluar')
  }

  const handleCreateProduct = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name || !price) {
      showToast('[!] Nama dan harga diperlukan')
      return
    }

    setSubmittingProduct(true)
    try {
      await productsApi.create({
        name,
        sku: sku || `SKU-${Date.now().toString().slice(-4)}`,
        description: desc,
        price: parseFloat(price),
        stock: parseInt(stock) || 0,
        weight: parseFloat(weight) || 0.5,
        imageUrl,
        isActive: true,
      })

      showToast('[✓] Produk baru berjaya ditambah!')
      setName('')
      setSku('')
      setDesc('')
      setPrice('')
      onProductCreated()
    } catch (err: any) {
      showToast(`[!] ${err.message || 'Gagal menambah produk'}`)
    } finally {
      setSubmittingProduct(false)
    }
  }

  if (!token) {
    return (
      <div
        style={{
          maxWidth: '440px',
          margin: '0 auto',
          background: 'var(--surface)',
          border: '1px solid var(--border)',
          padding: '28px',
        }}
      >
        <div className="font-mono" style={{ fontWeight: 700, fontSize: '16px', marginBottom: '8px' }}>
          🔒 LOG MASUK ADMINISTRATOR
        </div>
        <p style={{ fontSize: '12px', color: 'var(--text-muted)', marginBottom: '16px' }}>
          Hanya akaun dengan hak akses admin dibenarkan mengurus katalog dan pesanan kedai.
        </p>

        <form onSubmit={handleLogin} style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          <div>
            <label className="uppercase-label">EMEL ADMIN</label>
            <input
              type="email"
              required
              className="input-raw font-mono"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div>
            <label className="uppercase-label">KATA LALUAN</label>
            <input
              type="password"
              required
              className="input-raw font-mono"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>

          <button
            type="submit"
            disabled={loggingIn}
            className="btn-raw primary"
            style={{ marginTop: '8px', padding: '10px' }}
          >
            {loggingIn ? 'MEMERIKSA...' : 'LOG MASUK ADMIN →'}
          </button>

          <div
            className="font-mono"
            style={{
              background: 'var(--surface-alt)',
              border: '1px solid var(--border-muted)',
              padding: '8px',
              fontSize: '11px',
              marginTop: '8px',
            }}
          >
            Kredensial Default:<br />
            Emel: <code>admin@ecommerce.com</code><br />
            Kata laluan: <code>admin123</code>
          </div>
        </form>
      </div>
    )
  }

  return (
    <div>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          background: 'var(--surface)',
          border: '1px solid var(--border)',
          padding: '12px 18px',
          marginBottom: '16px',
          flexWrap: 'wrap',
          gap: '10px',
        }}
      >
        <div style={{ display: 'flex', gap: '8px' }}>
          <button
            onClick={() => setAdminTab('orders')}
            className={`btn-raw sm ${adminTab === 'orders' ? 'primary' : ''}`}
          >
            📦 REKOD PESANAN ({orders.length})
          </button>
          <button
            onClick={() => setAdminTab('add_product')}
            className={`btn-raw sm ${adminTab === 'add_product' ? 'primary' : ''}`}
          >
            + TAMBAH PRODUK BARU
          </button>
        </div>
        <button onClick={handleLogout} className="btn-raw sm">
          [LOG KELUAR ADMIN]
        </button>
      </div>

      {adminTab === 'orders' && (
        <div>
          {loadingOrders ? (
            <div
              className="font-mono"
              style={{
                background: '#FFF',
                padding: '30px',
                textAlign: 'center',
                border: '1px solid var(--border)',
              }}
            >
              MEMUAT REKOD PESANAN DARI DATABASE...
            </div>
          ) : orders.length === 0 ? (
            <div
              className="font-mono"
              style={{
                background: '#FFF',
                padding: '30px',
                textAlign: 'center',
                border: '1px solid var(--border)',
              }}
            >
              BELUM ADA PESANAN PELANGGAN DIBUAT LAGI.
            </div>
          ) : (
            <div style={{ overflowX: 'auto' }}>
              <table className="raw-table">
                <thead>
                  <tr>
                    <th>NO. PESANAN</th>
                    <th>PELANGGAN</th>
                    <th>NO. TEL</th>
                    <th>ALAMAT</th>
                    <th>KURIER</th>
                    <th>JUMLAH (RM)</th>
                    <th>STATUS</th>
                    <th>TARIKH</th>
                  </tr>
                </thead>
                <tbody>
                  {orders.map((o) => (
                    <tr key={o.id}>
                      <td className="font-mono" style={{ fontWeight: 700 }}>
                        {o.orderNumber}
                      </td>
                      <td>{o.customerName}</td>
                      <td className="font-mono">{o.customerPhone}</td>
                      <td style={{ fontSize: '12px' }}>
                        {o.city}, {o.state} ({o.postcode})
                      </td>
                      <td>{o.shipment?.courierName || 'J&T Express'}</td>
                      <td className="font-mono" style={{ fontWeight: 700 }}>
                        RM {o.totalAmount?.toFixed(2)}
                      </td>
                      <td>
                        <span
                          className={`badge-raw ${
                            o.status === 'paid' ? 'badge-paid' : 'badge-pending'
                          }`}
                        >
                          {o.status.toUpperCase()}
                        </span>
                      </td>
                      <td className="font-mono" style={{ fontSize: '11px' }}>
                        {o.createdAt ? new Date(o.createdAt).toLocaleDateString('ms-MY') : '-'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {adminTab === 'add_product' && (
        <div
          style={{
            maxWidth: '600px',
            margin: '0 auto',
            background: 'var(--surface)',
            border: '1px solid var(--border)',
            padding: '24px',
          }}
        >
          <div
            className="font-mono"
            style={{
              fontWeight: 700,
              fontSize: '15px',
              marginBottom: '14px',
              borderBottom: '1px solid var(--border)',
              paddingBottom: '8px',
            }}
          >
            BORANG MENAMBAH PRODUK BARU
          </div>

          <form onSubmit={handleCreateProduct} style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '10px' }}>
              <div>
                <label className="uppercase-label">NAMA PRODUK *</label>
                <input
                  type="text"
                  required
                  className="input-raw"
                  placeholder="Contoh: Baju Kemeja Oxford"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </div>
              <div>
                <label className="uppercase-label">SKU</label>
                <input
                  type="text"
                  className="input-raw font-mono"
                  placeholder="Contoh: OXF-001"
                  value={sku}
                  onChange={(e) => setSku(e.target.value)}
                />
              </div>
            </div>

            <div>
              <label className="uppercase-label">KETERANGAN PRODUK</label>
              <textarea
                rows={2}
                className="input-raw"
                placeholder="Spesifikasi & maklumat produk"
                value={desc}
                onChange={(e) => setDesc(e.target.value)}
              />
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '10px' }}>
              <div>
                <label className="uppercase-label">HARGA (RM) *</label>
                <input
                  type="number"
                  step="0.01"
                  required
                  className="input-raw font-mono"
                  placeholder="79.90"
                  value={price}
                  onChange={(e) => setPrice(e.target.value)}
                />
              </div>
              <div>
                <label className="uppercase-label">STOK UNIT</label>
                <input
                  type="number"
                  className="input-raw font-mono"
                  value={stock}
                  onChange={(e) => setStock(e.target.value)}
                />
              </div>
              <div>
                <label className="uppercase-label">BERAT (KG) *</label>
                <input
                  type="number"
                  step="0.01"
                  className="input-raw font-mono"
                  placeholder="0.300"
                  value={weight}
                  onChange={(e) => setWeight(e.target.value)}
                />
              </div>
            </div>

            <div>
              <label className="uppercase-label">URL GAMBAR PRODUK</label>
              <input
                type="url"
                className="input-raw font-mono"
                value={imageUrl}
                onChange={(e) => setImageUrl(e.target.value)}
              />
            </div>

            <button
              type="submit"
              disabled={submittingProduct}
              className="btn-raw primary"
              style={{ marginTop: '10px', padding: '12px' }}
            >
              {submittingProduct ? 'MENYIMPAN...' : '+ SIMPAN PRODUK KE DATABASE'}
            </button>
          </form>
        </div>
      )}
    </div>
  )
}
