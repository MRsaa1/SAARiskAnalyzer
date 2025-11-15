import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../lib/api'

export default function PortfolioImport() {
  const navigate = useNavigate()
  const [portfolioName, setPortfolioName] = useState('My Portfolio')
  const [positionsFile, setPositionsFile] = useState<File | null>(null)
  const [pricesFile, setPricesFile] = useState<File | null>(null)
  
  // Manual input - with current market prices (Nov 13, 2025)
  const [positions, setPositions] = useState([
    { symbol: 'SPY', quantity: 100, avgPrice: 595.12 },
    { symbol: 'TLT', quantity: 500, avgPrice: 94.23 },
    { symbol: 'GLD', quantity: 50, avgPrice: 234.56 },
    { symbol: 'BTC', quantity: 1, avgPrice: 99886.29 },
  ])

  const addPosition = () => {
    setPositions([...positions, { symbol: '', quantity: 0, avgPrice: 0 }])
  }

  const removePosition = (idx: number) => {
    setPositions(positions.filter((_, i) => i !== idx))
  }

  const updatePosition = async (idx: number, field: string, value: any) => {
    const updated = [...positions]
    updated[idx] = { ...updated[idx], [field]: value }
    setPositions(updated)
    
    // Auto-fetch price when symbol changes
    if (field === 'symbol' && value && value.length > 0) {
      try {
        const res = await api.get(`/market/price/${value.toUpperCase()}`)
        if (res.data && res.data.price) {
          updated[idx].avgPrice = parseFloat(res.data.price.toFixed(2))
          setPositions([...updated])
          console.log(`✅ Auto-filled price for ${value}: $${res.data.price}`)
        }
      } catch (err) {
        console.log(`⚠️ Could not fetch price for ${value}`)
      }
    }
  }

  const handleCSVUpload = async () => {
    if (!positionsFile || !pricesFile) {
      alert('Please select both files')
      return
    }

    try {
      const portfolio = await api.post('/portfolios', {
        name: portfolioName,
      })

      const positionsFormData = new FormData()
      positionsFormData.append('file', positionsFile)
      await api.post(`/portfolios/${portfolio.data.id}/positions:import`, positionsFormData)

      const pricesFormData = new FormData()
      pricesFormData.append('file', pricesFile)
      await api.post(`/portfolios/${portfolio.data.id}/prices:import`, pricesFormData)

      alert(`✅ Portfolio "${portfolioName}" created from CSV files!`)
      navigate('/', { state: { reload: true } })
    } catch (err: any) {
      alert('Error: ' + (err.response?.data?.error || err.message))
    }
  }

  const handleManualUpload = async () => {
    if (positions.length === 0) {
      alert('Add at least one position')
      return
    }

    try {
      const portfolio = await api.post('/portfolios', {
        name: portfolioName,
      })

      await api.post(`/portfolios/${portfolio.data.id}/positions`, {
        positions: positions.map(p => ({
          symbol: p.symbol,
          quantity: p.quantity,
          avg_price: p.avgPrice,
        })),
      })

      alert(`✅ Portfolio "${portfolioName}" created with ${positions.length} positions!`)
      navigate('/', { state: { reload: true } })
    } catch (err: any) {
      alert('Error: ' + (err.response?.data?.error || err.message))
    }
  }

  return (
    <div className="min-h-screen bg-bg text-fg p-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <button
            onClick={() => navigate('/')}
            className="text-accent hover:underline mb-4"
          >
            ← Back to Dashboard
          </button>
          <h1 className="text-3xl font-bold">SAA Risk Analyzer — Data Import</h1>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* CSV Import */}
          <div className="sa-card p-6">
            <h2 className="text-xl font-semibold mb-4">📂 CSV Import</h2>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">
                  Portfolio Name
                </label>
                <input
                  type="text"
                  value={portfolioName}
                  onChange={(e) => setPortfolioName(e.target.value)}
                  className="w-full px-3 py-2 bg-bg border border-border rounded"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">
                  Positions File (positions.csv)
                </label>
                <input
                  type="file"
                  accept=".csv"
                  onChange={(e) => setPositionsFile(e.target.files?.[0] || null)}
                  className="w-full"
                />
                <p className="text-xs text-muted mt-1">
                  Format: symbol,quantity,avg_price
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">
                  Prices File (prices.csv)
                </label>
                <input
                  type="file"
                  accept=".csv"
                  onChange={(e) => setPricesFile(e.target.files?.[0] || null)}
                  className="w-full"
                />
                <p className="text-xs text-muted mt-1">
                  Format: date,symbol,close
                </p>
              </div>

              <button
                onClick={handleCSVUpload}
                className="w-full bg-accent text-black py-2 px-4 rounded hover:opacity-90"
              >
                📊 Upload CSV
              </button>
            </div>
          </div>

          {/* Manual Input */}
          <div className="sa-card p-6">
            <h2 className="text-xl font-semibold mb-4">✍️ Manual Input</h2>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">
                  Portfolio Name
                </label>
                <input
                  type="text"
                  value={portfolioName}
                  onChange={(e) => setPortfolioName(e.target.value)}
                  className="w-full px-3 py-2 bg-bg border border-border rounded"
                />
              </div>

              <div>
                <h3 className="font-medium mb-2">Positions:</h3>
                <div className="space-y-2">
                  <div className="grid grid-cols-4 gap-2 text-sm font-medium text-muted">
                    <div>Symbol</div>
                    <div>Quantity</div>
                    <div>Avg Price</div>
                    <div>Action</div>
                  </div>
                  
                  {positions.map((pos, idx) => (
                    <div key={idx} className="grid grid-cols-4 gap-2 items-center">
                      <input
                        type="text"
                        value={pos.symbol}
                        onChange={(e) => updatePosition(idx, 'symbol', e.target.value)}
                        placeholder="SPY"
                        className="px-2 py-1 bg-bg border border-border rounded"
                      />
                      <input
                        type="number"
                        step="0.001"
                        value={pos.quantity}
                        onChange={(e) => updatePosition(idx, 'quantity', parseFloat(e.target.value) || 0)}
                        placeholder="1000"
                        className="px-2 py-1 bg-bg border border-border rounded"
                      />
                      <input
                        type="number"
                        step="0.01"
                        value={pos.avgPrice}
                        onChange={(e) => updatePosition(idx, 'avgPrice', parseFloat(e.target.value) || 0)}
                        placeholder="370"
                        className="px-2 py-1 bg-bg border border-border rounded"
                      />
                      <button
                        onClick={() => removePosition(idx)}
                        className="text-red-400 hover:text-red-300"
                      >
                        ✕
                      </button>
                    </div>
                  ))}
                </div>

                <button
                  onClick={addPosition}
                  className="mt-2 text-accent hover:underline"
                >
                  + Add Position
                </button>
              </div>

              <button
                onClick={handleManualUpload}
                className="w-full bg-accent text-black py-2 px-4 rounded hover:opacity-90"
              >
                ✍️ Create Portfolio
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
