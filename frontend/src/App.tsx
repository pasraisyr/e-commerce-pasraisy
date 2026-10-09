import { useState, useEffect } from 'react'
import type { Product, TabType } from './types'
import { productsApi } from './api'
import { CartProvider } from './context/CartContext'
import { Header } from './components/Header'
import { Footer } from './components/Footer'
import { Toast } from './components/Toast'
import { CatalogPage } from './pages/CatalogPage'
import { CartPage } from './pages/CartPage'
import { TrackOrderPage } from './pages/TrackOrderPage'
import { AdminPage } from './pages/AdminPage'

export default function App() {
  const [activeTab, setActiveTab] = useState<TabType>('katalog')
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState<boolean>(true)

  const fetchProducts = async () => {
    setLoading(true)
    try {
      const data = await productsApi.getAll()
      setProducts(Array.isArray(data) ? data : [])
    } catch (err) {
      console.error('Gagal muat produk:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchProducts()
  }, [])

  return (
    <CartProvider>
      <div style={{ maxWidth: '1240px', margin: '0 auto', padding: '24px 16px', minHeight: '100vh' }}>
        {/* Header Manifest */}
        <Header activeTab={activeTab} setActiveTab={setActiveTab} productCount={products.length} />

        {/* Dynamic Page Rendering */}
        <main>
          {activeTab === 'katalog' && <CatalogPage products={products} loading={loading} />}

          {activeTab === 'troli' && <CartPage onBackToCatalog={() => setActiveTab('katalog')} />}

          {activeTab === 'jejak' && <TrackOrderPage />}

          {activeTab === 'admin' && <AdminPage onProductCreated={fetchProducts} />}
        </main>

        {/* Footer & Notifications */}
        <Footer />
        <Toast />
      </div>
    </CartProvider>
  )
}
