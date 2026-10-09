import React, { useState } from 'react'
import type { CourierRate, CreateOrderPayload, CreateOrderResponse } from '../types'
import { useCart } from '../context/CartContext'
import { ordersApi, shippingApi } from '../api'

interface CartPageProps {
  onBackToCatalog: () => void
}

export const CartPage: React.FC<CartPageProps> = ({ onBackToCatalog }) => {
  const { cart, updateQty, subtotal, totalWeight, clearCart, showToast } = useCart()

  // Shipping & EasyParcel
  const [postcode, setPostcode] = useState<string>('50450')
  const [courierRates, setCourierRates] = useState<CourierRate[]>([])
  const [selectedCourier, setSelectedCourier] = useState<CourierRate | null>(null)
  const [checkingRates, setCheckingRates] = useState<boolean>(false)

  // Form Fields
  const [isGuest, setIsGuest] = useState<boolean>(true)
  const [customerName, setCustomerName] = useState<string>('')
  const [customerEmail, setCustomerEmail] = useState<string>('')
  const [customerPhone, setCustomerPhone] = useState<string>('')
  const [shippingAddress, setShippingAddress] = useState<string>('')
  const [city, setCity] = useState<string>('Kuala Lumpur')
  const [stateName, setStateName] = useState<string>('Kuala Lumpur')
  const [checkoutNotes, setCheckoutNotes] = useState<string>('')
  const [submittingOrder, setSubmittingOrder] = useState<boolean>(false)
  const [placedOrder, setPlacedOrder] = useState<CreateOrderResponse | null>(null)

  const shippingFee = selectedCourier ? selectedCourier.price : 8.0
  const grandTotal = subtotal + (cart.length > 0 ? shippingFee : 0)

  const handleCheckShipping = async () => {
    if (!postcode) {
      showToast('[!] Sila masukkan poskod penghantaran')
      return
    }
    setCheckingRates(true)
    try {
      const data = await shippingApi.checkRates(postcode, Math.max(totalWeight, 0.2))
      setCourierRates(data)
      if (data.length > 0) {
        setSelectedCourier(data[0])
        showToast(`[✓] Kadar EasyParcel dikemaskini (${data.length} kurier)`)
      }
    } catch {
      showToast('[!] Gagal menyemak kadar pos')
    } finally {
      setCheckingRates(false)
    }
  }

  const handleCheckout = async (e: React.FormEvent) => {
    e.preventDefault()
    if (cart.length === 0) {
      showToast('[!] Troli anda masih kosong')
      return
    }
    if (!customerName || !customerEmail || !customerPhone || !shippingAddress || !postcode) {
      showToast('[!] Sila lengkapkan maklumat bertanda *')
      return
    }

    setSubmittingOrder(true)
    try {
      const payload: CreateOrderPayload = {
        customerName,
        customerEmail,
        customerPhone,
        shippingAddress,
        city,
        state: stateName,
        postcode,
        courierName: selectedCourier ? selectedCourier.courierName : 'J&T Express',
        shippingRate: shippingFee,
        notes: checkoutNotes,
        items: cart.map((i) => ({ productId: i.product.id, quantity: i.quantity })),
      }

      const res = await ordersApi.create(payload)
      setPlacedOrder(res)
      clearCart()
      showToast(`[✓] Pesanan berjaya dijana: ${res.order.orderNumber}`)
    } catch (err: any) {
      showToast(`[!] ${err.message || 'Gagal membuat pesanan'}`)
    } finally {
      setSubmittingOrder(false)
    }
  }

  if (placedOrder) {
    return (
      <div
        style={{
          background: 'var(--surface)',
          border: '2px solid var(--border)',
          padding: '30px',
          maxWidth: '680px',
          margin: '0 auto',
        }}
      >
        <div
          className="font-mono"
          style={{
            borderBottom: '2px solid var(--border)',
            paddingBottom: '12px',
            marginBottom: '16px',
          }}
        >
          <span className="badge-raw badge-paid">PESANAN DITERIMA: {placedOrder.order.orderNumber}</span>
          <h2 style={{ fontSize: '18px', marginTop: '10px' }}>TERIMA KASIH ATAS PESANAN ANDA</h2>
        </div>

        <p style={{ fontSize: '13px', marginBottom: '16px' }}>
          Pesanan anda telah direkodkan dalam sistem. Sila buat pembayaran melalui ToyyibPay untuk melengkapkan pesanan.
        </p>

        <table className="raw-table" style={{ marginBottom: '20px' }}>
          <tbody>
            <tr>
              <td className="font-mono" style={{ width: '180px', background: 'var(--surface-alt)' }}>PELANGGAN</td>
              <td>{placedOrder.order.customerName} ({placedOrder.order.customerEmail})</td>
            </tr>
            <tr>
              <td className="font-mono" style={{ background: 'var(--surface-alt)' }}>ALAMAT HANTAR</td>
              <td>{placedOrder.order.shippingAddress}, {placedOrder.order.postcode} {placedOrder.order.city}</td>
            </tr>
            <tr>
              <td className="font-mono" style={{ background: 'var(--surface-alt)' }}>JUMLAH BAYARAN</td>
              <td className="font-mono" style={{ fontWeight: 700, fontSize: '16px' }}>
                RM {placedOrder.order.totalAmount?.toFixed(2)}
              </td>
            </tr>
            <tr>
              <td className="font-mono" style={{ background: 'var(--surface-alt)' }}>KURIER (EASYPARCEL)</td>
              <td>{placedOrder.order.shipment?.courierName || 'J&T Express'} (Kadar: RM {placedOrder.order.shippingFee?.toFixed(2)})</td>
            </tr>
          </tbody>
        </table>

        <div style={{ display: 'flex', gap: '10px', flexWrap: 'wrap' }}>
          <a
            href={placedOrder.paymentUrl}
            target="_blank"
            rel="noreferrer"
            className="btn-raw primary"
            style={{ textDecoration: 'none', padding: '12px 20px', flex: 1, textAlign: 'center' }}
          >
            💳 BAYAR SEKARANG VIA TOYYIBPAY (FPX / KAD) →
          </a>
          <button
            onClick={() => {
              setPlacedOrder(null)
              onBackToCatalog()
            }}
            className="btn-raw"
          >
            KEMBALI KE KEDAI
          </button>
        </div>
      </div>
    )
  }

  if (cart.length === 0) {
    return (
      <div
        className="font-mono"
        style={{
          background: 'var(--surface)',
          border: '1px solid var(--border)',
          padding: '50px',
          textAlign: 'center',
        }}
      >
        <p style={{ fontSize: '14px', marginBottom: '16px' }}>TROLI ANDA MASIH KOSONG.</p>
        <button onClick={onBackToCatalog} className="btn-raw primary">
          ← PILIH BARANG DARI KATALOG
        </button>
      </div>
    )
  }

  return (
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(360px, 1fr))', gap: '20px' }}>
      {/* KIRI: ITEM DALAM TROLI & EASYPARCEL */}
      <div>
        <div
          style={{
            background: 'var(--surface-alt)',
            border: '1px solid var(--border)',
            borderBottom: 'none',
            padding: '12px 16px',
            fontWeight: 700,
          }}
          className="font-mono"
        >
          SENARAI BARANG DALAM TROLI ({cart.length})
        </div>

        <div style={{ overflowX: 'auto' }}>
          <table className="raw-table">
            <thead>
              <tr>
                <th>PRODUK</th>
                <th style={{ width: '90px' }}>HARGA</th>
                <th style={{ width: '100px' }}>KUANTITI</th>
                <th style={{ width: '90px' }}>JUMLAH</th>
              </tr>
            </thead>
            <tbody>
              {cart.map((item) => (
                <tr key={item.product.id}>
                  <td>
                    <div style={{ fontWeight: 600 }}>{item.product.name}</div>
                    <div className="font-mono" style={{ fontSize: '11px', color: 'var(--text-muted)' }}>
                      SKU: {item.product.sku} | {((item.product.weight || 0.5) * item.quantity).toFixed(3)} KG
                    </div>
                  </td>
                  <td className="font-mono">RM {item.product.price.toFixed(2)}</td>
                  <td>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
                      <button
                        onClick={() => updateQty(item.product.id, -1)}
                        className="btn-raw sm"
                        style={{ padding: '2px 6px' }}
                      >
                        -
                      </button>
                      <span className="font-mono" style={{ minWidth: '20px', textAlign: 'center' }}>
                        {item.quantity}
                      </span>
                      <button
                        onClick={() => updateQty(item.product.id, 1)}
                        className="btn-raw sm"
                        style={{ padding: '2px 6px' }}
                      >
                        +
                      </button>
                    </div>
                  </td>
                  <td className="font-mono" style={{ fontWeight: 700 }}>
                    RM {(item.product.price * item.quantity).toFixed(2)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* KIRAAN KOS POS EASYPARCEL */}
        <div
          style={{
            marginTop: '20px',
            background: 'var(--surface)',
            border: '1px solid var(--border)',
            padding: '16px',
          }}
        >
          <div className="font-mono" style={{ fontWeight: 700, marginBottom: '10px' }}>
            🚚 KIRAAN KOS PENGHANTARAN (EASYPARCEL)
          </div>
          <p style={{ fontSize: '12px', color: 'var(--text-muted)', marginBottom: '12px' }}>
            Jumlah berat pesanan: <strong className="font-mono">{totalWeight.toFixed(3)} KG</strong>. Masukkan poskod untuk semak kurier sebenar.
          </p>

          <div style={{ display: 'flex', gap: '8px', marginBottom: '14px' }}>
            <input
              type="text"
              className="input-raw font-mono"
              placeholder="Poskod (cth: 50450)"
              value={postcode}
              onChange={(e) => setPostcode(e.target.value)}
              style={{ maxWidth: '160px' }}
            />
            <button
              type="button"
              onClick={handleCheckShipping}
              disabled={checkingRates}
              className="btn-raw"
            >
              {checkingRates ? 'MENYEMAK...' : 'KIRA KOS POS'}
            </button>
          </div>

          {courierRates.length > 0 && (
            <div style={{ borderTop: '1px solid var(--border-muted)', paddingTop: '10px' }}>
              <span className="uppercase-label">PILIH KURIER:</span>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '6px', marginTop: '6px' }}>
                {courierRates.map((cr) => (
                  <label
                    key={cr.serviceCode}
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'space-between',
                      padding: '8px 10px',
                      border: selectedCourier?.serviceCode === cr.serviceCode ? '2px solid var(--border)' : '1px solid var(--border-muted)',
                      background: selectedCourier?.serviceCode === cr.serviceCode ? 'var(--surface-alt)' : '#FFF',
                      cursor: 'pointer',
                    }}
                  >
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                      <input
                        type="radio"
                        name="courier"
                        checked={selectedCourier?.serviceCode === cr.serviceCode}
                        onChange={() => setSelectedCourier(cr)}
                      />
                      <div>
                        <strong style={{ fontSize: '13px' }}>{cr.courierName}</strong>
                        <span style={{ fontSize: '11px', color: 'var(--text-muted)', marginLeft: '6px' }}>
                          ({cr.deliveryEta})
                        </span>
                      </div>
                    </div>
                    <span className="font-mono" style={{ fontWeight: 700 }}>RM {cr.price.toFixed(2)}</span>
                  </label>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* KANAN: CHECKOUT FORM */}
      <div
        style={{
          background: 'var(--surface)',
          border: '1px solid var(--border)',
          padding: '20px',
        }}
      >
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            marginBottom: '16px',
            borderBottom: '1px solid var(--border-muted)',
            paddingBottom: '10px',
          }}
        >
          <span className="font-mono" style={{ fontWeight: 700 }}>MAKLUMAT PENGHANTARAN</span>
          <div style={{ display: 'flex', gap: '4px' }}>
            <button
              type="button"
              onClick={() => setIsGuest(true)}
              className={`btn-raw sm ${isGuest ? 'primary' : ''}`}
            >
              BELI TETAMU (GUEST)
            </button>
            <button
              type="button"
              onClick={() => {
                setIsGuest(false)
                showToast('[i] Anda boleh terus isi maklumat di bawah')
              }}
              className={`btn-raw sm ${!isGuest ? 'primary' : ''}`}
            >
              AHLI
            </button>
          </div>
        </div>

        <form onSubmit={handleCheckout}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <div>
              <label className="uppercase-label">NAMA PENUH *</label>
              <input
                type="text"
                required
                className="input-raw"
                placeholder="Contoh: Muhammad Ali"
                value={customerName}
                onChange={(e) => setCustomerName(e.target.value)}
              />
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '10px' }}>
              <div>
                <label className="uppercase-label">EMEL *</label>
                <input
                  type="email"
                  required
                  className="input-raw"
                  placeholder="nama@email.com"
                  value={customerEmail}
                  onChange={(e) => setCustomerEmail(e.target.value)}
                />
              </div>
              <div>
                <label className="uppercase-label">NO. TELEFON *</label>
                <input
                  type="tel"
                  required
                  className="input-raw"
                  placeholder="0123456789"
                  value={customerPhone}
                  onChange={(e) => setCustomerPhone(e.target.value)}
                />
              </div>
            </div>

            <div>
              <label className="uppercase-label">ALAMAT LENGKAP PENGHANTARAN *</label>
              <textarea
                required
                rows={2}
                className="input-raw"
                placeholder="No rumah, jalan, taman/apartmen"
                value={shippingAddress}
                onChange={(e) => setShippingAddress(e.target.value)}
              />
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '8px' }}>
              <div>
                <label className="uppercase-label">POSKOD *</label>
                <input
                  type="text"
                  required
                  className="input-raw font-mono"
                  value={postcode}
                  onChange={(e) => setPostcode(e.target.value)}
                />
              </div>
              <div>
                <label className="uppercase-label">BANDAR *</label>
                <input
                  type="text"
                  required
                  className="input-raw"
                  value={city}
                  onChange={(e) => setCity(e.target.value)}
                />
              </div>
              <div>
                <label className="uppercase-label">NEGERI *</label>
                <input
                  type="text"
                  required
                  className="input-raw"
                  value={stateName}
                  onChange={(e) => setStateName(e.target.value)}
                />
              </div>
            </div>

            <div>
              <label className="uppercase-label">NOTA PESANAN (PILIHAN)</label>
              <input
                type="text"
                className="input-raw"
                placeholder="Arahan khas pos atau nota pembungkusan"
                value={checkoutNotes}
                onChange={(e) => setCheckoutNotes(e.target.value)}
              />
            </div>

            <table className="raw-table" style={{ marginTop: '12px' }}>
              <tbody>
                <tr>
                  <td className="font-mono">Subtotal ({cart.length} item)</td>
                  <td className="font-mono" style={{ textAlign: 'right' }}>RM {subtotal.toFixed(2)}</td>
                </tr>
                <tr>
                  <td className="font-mono">Kos Pos ({selectedCourier?.courierName || 'Standard'})</td>
                  <td className="font-mono" style={{ textAlign: 'right' }}>RM {shippingFee.toFixed(2)}</td>
                </tr>
                <tr style={{ background: 'var(--surface-alt)', fontWeight: 700 }}>
                  <td className="font-mono">JUMLAH KESELURUHAN</td>
                  <td className="font-mono" style={{ textAlign: 'right', fontSize: '16px' }}>
                    RM {grandTotal.toFixed(2)}
                  </td>
                </tr>
              </tbody>
            </table>

            <button
              type="submit"
              disabled={submittingOrder}
              className="btn-raw primary"
              style={{ padding: '14px', marginTop: '10px', fontSize: '13px' }}
            >
              {submittingOrder ? 'MENJANA PESANAN...' : '💳 SAHKAN & BAYAR VIA TOYYIBPAY →'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
