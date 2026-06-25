import { useState, useEffect } from 'react'
import './App.css'

interface Product {
  sku: string
  name: string
  current_stock: number
}

interface Event {
  event_id: string
  sku: string
  type: string
  quantity: number
  occurred_at: string
}

function App() {
  const [products, setProducts] = useState<Product[]>([])
  const [selectedSku, setSelectedSku] = useState<string | null>(null)
  const [movements, setMovements] = useState<Event[]>([])
  const [page, setPage] = useState<number>(1)
  const limit = 50

  useEffect(() => {
    fetch('http://localhost:8080/api/products')
      .then(res => res.json())
      .then(data => setProducts(data || []))
      .catch(err => console.error("Error fetching products:", err))
  }, [])

  useEffect(() => {
    if (selectedSku) {
      fetch(`http://localhost:8080/api/products/${selectedSku}/movements?page=${page}&limit=${limit}`)
        .then(res => res.json())
        .then(data => setMovements(data || []))
        .catch(err => console.error("Error fetching movements:", err))
    } else {
      setMovements([])
    }
  }, [selectedSku, page])

  const handleSelectProduct = (sku: string) => {
    if (sku !== selectedSku) {
      setSelectedSku(sku)
      setPage(1)
    }
  }

  return (
    <div className="container">
      <div className="sidebar">
        <h2>Products</h2>
        <ul className="product-list">
          {products.map(p => (
            <li 
              key={p.sku} 
              className={selectedSku === p.sku ? 'selected' : ''}
              onClick={() => handleSelectProduct(p.sku)}
            >
              <div className="product-info">
                <strong>{p.name}</strong> <span className="sku">{p.sku}</span>
              </div>
              <div className="stock">Stock: {p.current_stock}</div>
            </li>
          ))}
        </ul>
      </div>
      <div className="content">
        <h2>Movements {selectedSku && `for ${selectedSku}`}</h2>
        {!selectedSku ? (
          <p>Select a product to view its movements.</p>
        ) : (
          <div>
            <div className="table-container">
              <table>
                <thead>
                  <tr>
                    <th>Event ID</th>
                    <th>Type</th>
                    <th>Quantity</th>
                    <th>Date</th>
                  </tr>
                </thead>
                <tbody>
                  {movements.map(m => (
                    <tr key={m.event_id}>
                      <td>{m.event_id}</td>
                      <td className={m.type === 'IN' ? 'text-green' : 'text-red'}>{m.type}</td>
                      <td>{m.quantity}</td>
                      <td>{new Date(m.occurred_at).toLocaleString()}</td>
                    </tr>
                  ))}
                  {movements.length === 0 && (
                    <tr>
                      <td colSpan={4} className="text-center">No movements found.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
            <div className="pagination">
              <button 
                onClick={() => setPage(p => Math.max(1, p - 1))} 
                disabled={page === 1}
              >
                Previous
              </button>
              <span className="page-info">Page {page}</span>
              <button 
                onClick={() => setPage(p => p + 1)} 
                disabled={movements.length < limit}
              >
                Next
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default App
