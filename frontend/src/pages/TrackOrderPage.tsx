import React, { useState } from 'react'
import type { Order } from '../types'
import { ordersApi } from '../api'

export const TrackOrderPage: React.FC = () => {
  const [orderNumber, setOrderNumber] = useState<string>('')
  const [order, setOrder] = useState<Order | null>(null)
  const [loading, setLoading] = useState<boolean>(false)
  const [errorMsg, setErrorMsg] = useState<string>('')

  const handleTrack = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!orderNumber.trim()) return

    setLoading(true)
    setErrorMsg('')
    setOrder(null)

    try {
      const data = await ordersApi.track(orderNumber.trim())
      setOrder(data)
    } catch {
      setErrorMsg('Pesanan tidak dijumpai. Sila pastikan nombor pesanan tepat (contoh: ORD-20261009-XXXX).')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: '700px', margin: '0 auto' }}>
      <div
        style={{
          background: 'var(--surface)',
          border: '1px solid var(--border)',
          padding: '24px',
        }}
      >
        <div className="font-mono" style={{ fontWeight: 700, fontSize: '16px', marginBottom: '8px' }}>
          🔍 JEJAK STATUS PESANAN ANDA
        </div>
        <p style={{ fontSize: '13px', color: 'var(--text-muted)', marginBottom: '16px' }}>
          Masukkan nombor pesanan yang anda terima semasa checkout (contoh: <code>ORD-20261009-1234</code>).
        </p>

        <form onSubmit={handleTrack} style={{ display: 'flex', gap: '8px', marginBottom: '20px' }}>
          <input
            type="text"
            required
            className="input-raw font-mono"
            placeholder="Contoh: ORD-20261009-XXXX"
            value={orderNumber}
            onChange={(e) => setOrderNumber(e.target.value)}
          />
          <button type="submit" disabled={loading} className="btn-raw primary">
            {loading ? 'MENYEMAK...' : 'CARI'}
          </button>
        </form>

        {errorMsg && (
          <div
            className="font-mono"
            style={{
              background: '#FDF2F2',
              border: '1px solid #E5989B',
              color: '#8A1F1D',
              padding: '12px',
              fontSize: '12px',
            }}
          >
            {errorMsg}
          </div>
        )}

        {order && (
          <div>
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: '12px',
                borderBottom: '2px solid var(--border)',
                paddingBottom: '8px',
              }}
            >
              <span className="font-mono" style={{ fontWeight: 700 }}>
                NO. PESANAN: {order.orderNumber}
              </span>
              <span
                className={`badge-raw ${order.status === 'paid' ? 'badge-paid' : 'badge-pending'}`}
              >
                STATUS: {order.status.toUpperCase()}
              </span>
            </div>

            <table className="raw-table" style={{ marginBottom: '16px' }}>
              <tbody>
                <tr>
                  <td className="font-mono" style={{ width: '160px', background: 'var(--surface-alt)' }}>
                    NAMA PELANGGAN
                  </td>
                  <td>{order.customerName}</td>
                </tr>
                <tr>
                  <td className="font-mono" style={{ background: 'var(--surface-alt)' }}>
                    ALAMAT PENGHANTARAN
                  </td>
                  <td>
                    {order.shippingAddress}, {order.postcode} {order.city}, {order.state}
                  </td>
                </tr>
                <tr>
                  <td className="font-mono" style={{ background: 'var(--surface-alt)' }}>
                    PENGHANTARAN (EASYPARCEL)
                  </td>
                  <td>
                    {order.shipment?.courierName || 'EasyParcel'} | No. Tracking:{' '}
                    {order.shipment?.trackingNumber || 'Akan dikemaskini setelah diproses'}
                  </td>
                </tr>
                <tr>
                  <td className="font-mono" style={{ background: 'var(--surface-alt)' }}>
                    JUMLAH BAYARAN
                  </td>
                  <td className="font-mono" style={{ fontWeight: 700 }}>
                    RM {order.totalAmount?.toFixed(2)}
                  </td>
                </tr>
              </tbody>
            </table>

            <span className="uppercase-label">ITEM DIPESAN:</span>
            <table className="raw-table" style={{ marginTop: '6px' }}>
              <thead>
                <tr>
                  <th>ITEM</th>
                  <th style={{ width: '80px' }}>KUANTITI</th>
                  <th style={{ width: '100px' }}>HARGA</th>
                  <th style={{ width: '100px' }}>JUMLAH</th>
                </tr>
              </thead>
              <tbody>
                {order.items?.map((it) => (
                  <tr key={it.id}>
                    <td>{it.productName}</td>
                    <td className="font-mono">{it.quantity}</td>
                    <td className="font-mono">RM {it.price.toFixed(2)}</td>
                    <td className="font-mono" style={{ fontWeight: 600 }}>
                      RM {it.subtotal.toFixed(2)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
