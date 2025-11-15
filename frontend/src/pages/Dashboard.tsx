import { useEffect, useState, useRef } from 'react'
import { useLocation } from 'react-router-dom'
import api from '../lib/api'

// Dashboard v2.0 - Restored to original design
export default function DashboardPage() {
  const location = useLocation()
  const [data, setData] = useState<any>(null)
  const [portfolios, setPortfolios] = useState<any[]>([])
  const [selectedPortfolio, setSelectedPortfolio] = useState<string>('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  
  const [editingPortfolio, setEditingPortfolio] = useState<any>(null)
  const [editName, setEditName] = useState('')
  const [editDescription, setEditDescription] = useState('')
  const [editingPosition, setEditingPosition] = useState<any>(null)
  const [editPositionSymbol, setEditPositionSymbol] = useState('')
  const [editPositionQuantity, setEditPositionQuantity] = useState(0)
  const [editPositionPrice, setEditPositionPrice] = useState(0)
  const [positionPrices, setPositionPrices] = useState<Record<string, number>>({})
  const [refreshingPrices, setRefreshingPrices] = useState(false)

  // Define loadDashboard first to avoid reference errors
  const loadDashboard = async (portfolioId?: string) => {
    try {
      console.log('📊 Loading dashboard for portfolio:', portfolioId || selectedPortfolio || 'none')
      // Используем portfolio_id если есть
      const pid = portfolioId || selectedPortfolio
      if (pid) {
        const res = await api.get(`/risk/dashboard?portfolio_id=${pid}`)
        setData(res.data)
        setError(null)
        console.log('✅ Dashboard data loaded for portfolio:', pid)
      } else {
        // Fallback на mock данные если нет портфолио
        const res = await api.get('/dashboard')
        setData(res.data)
        setError(null)
        console.log('✅ Dashboard data loaded (fallback)')
      }
    } catch (err: any) {
      console.error('❌ Error loading dashboard:', err)
      setError(err.message || 'Failed to load dashboard')
      // Fallback на mock при ошибке
      try {
        const res = await api.get('/dashboard')
        setData(res.data)
        setError(null)
        console.log('✅ Dashboard data loaded (fallback after error)')
      } catch (fallbackErr) {
        console.error('❌ Fallback failed:', fallbackErr)
        // Set empty data to prevent black screen
        setData({
          var_1d: 0,
          cvar_1d: 0,
          vol: 0,
          contributors: []
        })
      }
    }
  }

  // Define loadPositionPrices
  const loadPositionPrices = async (portfolioId: string, showLoading = false) => {
    if (showLoading) {
      setRefreshingPrices(true)
    }
    
    try {
      // Reload portfolio to get latest positions
      const portfolioRes = await api.get(`/portfolios/${portfolioId}`)
      const portfolio = portfolioRes.data

      if (!portfolio || !portfolio.positions || portfolio.positions.length === 0) {
        if (showLoading) setRefreshingPrices(false)
        return
      }

      const prices: Record<string, number> = {}
      
      // Fetch prices for all positions in parallel
      const pricePromises = portfolio.positions.map(async (pos: any) => {
        const symbol = pos.asset?.symbol
        if (!symbol) return
        
        try {
          const res = await api.get(`/market/price/${symbol}`)
          if (res.data && res.data.price) {
            prices[symbol] = res.data.price
            console.log(`✅ Updated price for ${symbol}: $${res.data.price}`)
          }
        } catch (err) {
          console.error(`Failed to fetch price for ${symbol}:`, err)
          // Use avg price as fallback
          if (pos.avg_price) {
            prices[symbol] = pos.avg_price
          }
        }
      })

      await Promise.all(pricePromises)
      setPositionPrices(prices)
      
      // Update portfolios state to reflect latest data
      const updatedPortfolios = portfolios.map(p => 
        p.id === portfolioId ? portfolio : p
      )
      setPortfolios(updatedPortfolios)
      
      if (showLoading) setRefreshingPrices(false)
    } catch (err) {
      console.error('Failed to load position prices:', err)
      if (showLoading) setRefreshingPrices(false)
    }
  }

  // Load portfolios function
  const loadPortfolios = async (preserveSelection = false) => {
    try {
      setLoading(true)
      setError(null)
      console.log('🔄 Loading portfolios...')
      const res = await api.get('/portfolios')
      const portfoliosData = res.data || []
      console.log('📊 Loaded portfolios:', portfoliosData)
      console.log('📊 Number of portfolios:', portfoliosData.length)
      
      if (portfoliosData.length === 0) {
        console.log('⚠️ No portfolios found')
        setPortfolios([])
        setSelectedPortfolio('')
        setData(null)
        setLoading(false)
        return
      }
      
      // Log positions with their IDs
      portfoliosData.forEach((p: any) => {
        console.log(`📁 Portfolio: ${p.name} (ID: ${p.id})`)
        if (p.positions && p.positions.length > 0) {
          console.log(`  Positions (${p.positions.length}):`, p.positions.map((pos: any) => ({
            id: pos.id,
            symbol: pos.asset?.symbol,
            quantity: pos.quantity
          })))
        } else {
          console.log(`  No positions`)
        }
      })
      
      // Update portfolios state first
      setPortfolios(portfoliosData)
      
      // Handle selection
      if (preserveSelection && selectedPortfolio) {
        const stillExists = portfoliosData.find((p: any) => p.id === selectedPortfolio)
        if (stillExists) {
          console.log('✅ Keeping current selection:', selectedPortfolio)
          // Just reload data for current portfolio
          await loadDashboard(selectedPortfolio)
          setTimeout(() => loadPositionPrices(selectedPortfolio), 500)
          setLoading(false)
          return
        }
      }
      
      // Select first portfolio (or newly created one)
      const firstId = portfoliosData[0].id
      console.log('✅ Auto-selecting first portfolio:', firstId)
      
      // Always update selection to ensure new portfolios are shown
      setSelectedPortfolio(firstId)
      
      // Загружаем данные для первого портфолио
      await loadDashboard(firstId)
      // Загружаем цены для первого портфолио
      setTimeout(() => loadPositionPrices(firstId), 500)
      setLoading(false)
      
    } catch (err: any) {
      console.error('❌ Error loading portfolios:', err)
      setError(err.message || 'Failed to load portfolios')
      setPortfolios([])
      setSelectedPortfolio('')
      setData(null)
      setLoading(false)
    }
  }

  // Load portfolios on mount
  useEffect(() => {
    console.log('🚀 Dashboard mounted, loading portfolios...')
    loadPortfolios().catch((err) => {
      console.error('❌ Failed to load portfolios:', err)
      setLoading(false)
      setError('Failed to load portfolios')
    })
  }, [])

  // Reload when location changes (user navigates back from import page)
  const prevPathnameRef = useRef<string>(location.pathname)
  const hasReloadedRef = useRef<boolean>(false)
  
  useEffect(() => {
    // Only reload if:
    // 1. We're on dashboard
    // 2. We came from another page (pathname changed)
    // 3. reload flag is set
    // 4. We haven't already reloaded for this navigation
    if (location.pathname === '/' && 
        prevPathnameRef.current !== '/' && 
        location.state?.reload &&
        !hasReloadedRef.current) {
      console.log('🔄 Navigated to dashboard, reloading portfolios...')
      hasReloadedRef.current = true
      prevPathnameRef.current = location.pathname
      setTimeout(() => {
        loadPortfolios(false)
        // Reset reload flag after a delay
        setTimeout(() => {
          hasReloadedRef.current = false
        }, 1000)
      }, 200)
    } else if (location.pathname !== prevPathnameRef.current) {
      // Update pathname but don't reload
      prevPathnameRef.current = location.pathname
      hasReloadedRef.current = false
    }
  }, [location.pathname, location.state])

  useEffect(() => {
    if (!selectedPortfolio) return
    
    console.log('📊 Loading dashboard for selected portfolio:', selectedPortfolio)
    loadDashboard(selectedPortfolio).catch(err => {
      console.error('Error loading dashboard:', err)
    })
    loadPositionPrices(selectedPortfolio).catch(err => {
      console.error('Error loading position prices:', err)
    })
  }, [selectedPortfolio])

  // Auto-refresh prices every 30 seconds
  useEffect(() => {
    if (!selectedPortfolio) return
    
    const interval = setInterval(() => {
      const portfolio = portfolios.find(p => p.id === selectedPortfolio)
      if (portfolio && portfolio.positions && portfolio.positions.length > 0) {
        loadPositionPrices(selectedPortfolio)
      }
    }, 30000) // 30 seconds
    
    return () => clearInterval(interval)
  }, [selectedPortfolio, portfolios])

  const recalculateRisk = async () => {
    if (!selectedPortfolio) {
      alert('Please create a portfolio first!')
      return
    }
    
    setLoading(true)
    try {
      console.log('Recalculating risk for portfolio:', selectedPortfolio)
      const res = await api.get(`/risk/dashboard?portfolio_id=${selectedPortfolio}`)
      console.log('Risk calculation result:', res.data)
      setData(res.data)
      alert('✅ Risk recalculated successfully!')
    } catch (err: any) {
      console.error('Risk calculation error:', err)
      alert('Calculation error: ' + (err.response?.data?.error || err.message))
    }
    setLoading(false)
  }

  const handleEditPortfolio = (portfolio: any) => {
    setEditingPortfolio(portfolio)
    setEditName(portfolio.name)
    setEditDescription(portfolio.description || '')
  }

  const handleSavePortfolio = async () => {
    if (!editingPortfolio) return
    
    try {
      await api.put(`/portfolios/${editingPortfolio.id}`, {
        name: editName,
        description: editDescription,
      })
      await loadPortfolios()
      setEditingPortfolio(null)
      alert('✅ Portfolio updated successfully!')
    } catch (err: any) {
      alert('Error: ' + (err.response?.data?.error || err.message))
    }
  }

  const handleDeletePortfolio = async (portfolioId: string) => {
    if (!confirm('Are you sure you want to delete this portfolio? All positions will be deleted.')) {
      return
    }
    
    try {
      await api.delete(`/portfolios/${portfolioId}`)
      await loadPortfolios()
      if (selectedPortfolio === portfolioId) {
        setSelectedPortfolio('')
        setData(null)
      }
      alert('✅ Portfolio deleted successfully!')
    } catch (err: any) {
      alert('Error: ' + (err.response?.data?.error || err.message))
    }
  }

  const handleEditPosition = async (position: any) => {
    // Fetch current price for the symbol
    try {
      const priceRes = await api.get(`/market/price/${position.asset.symbol}`)
      setEditPositionPrice(priceRes.data.price)
    } catch (err) {
      setEditPositionPrice(position.avg_price)
    }
    
    setEditingPosition(position)
    setEditPositionSymbol(position.asset.symbol)
    setEditPositionQuantity(position.quantity)
  }

  const handleSavePosition = async () => {
    if (!editingPosition || !selectedPortfolio) return
    
    try {
      await api.put(`/portfolios/${selectedPortfolio}/positions/${editingPosition.id}`, {
        symbol: editPositionSymbol,
        quantity: editPositionQuantity,
        avg_price: editPositionPrice,
      })
      await loadPortfolios()
      await loadDashboard(selectedPortfolio)
      await loadPositionPrices(selectedPortfolio)
      setEditingPosition(null)
      alert('✅ Position updated successfully!')
    } catch (err: any) {
      alert('Error: ' + (err.response?.data?.error || err.message))
    }
  }

  const handleDeletePosition = async (positionId: string) => {
    if (!selectedPortfolio) {
      alert('No portfolio selected')
      return
    }
    if (!positionId) {
      alert('Position ID is missing')
      return
    }
    if (!confirm('Are you sure you want to delete this position?')) {
      return
    }
    
    // Validate UUID format (basic check)
    const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i
    if (!uuidRegex.test(positionId)) {
      console.error('Invalid UUID format:', positionId)
      alert('Invalid position ID format')
      return
    }
    
    const url = `/portfolios/${selectedPortfolio}/positions/${positionId}`
    console.log('Deleting position:', {
      positionId,
      portfolioId: selectedPortfolio,
      url,
      fullUrl: `${api.defaults.baseURL}${url}`
    })
    
    try {
      const response = await api.delete(url)
      console.log('Delete response:', response.data)
      await loadPortfolios()
      await loadDashboard(selectedPortfolio)
      await loadPositionPrices(selectedPortfolio)
      alert('✅ Position deleted successfully!')
    } catch (err: any) {
      console.error('Delete position error:', err)
      console.error('Error details:', {
        status: err.response?.status,
        statusText: err.response?.statusText,
        data: err.response?.data,
        message: err.message,
        url: err.config?.url,
        method: err.config?.method
      })
      const errorMsg = err.response?.data?.error || err.message || 'Unknown error'
      alert(`Error deleting position: ${errorMsg}`)
    }
  }

  // Debug: log current state
  console.log('🔍 Dashboard render state:', {
    loading,
    error,
    portfoliosCount: portfolios.length,
    selectedPortfolio,
    hasData: !!data
  })
  
  // Early returns - check loading and error states first
  if (loading) {
    console.log('⏳ Showing loading screen')
    return (
      <div className="min-h-screen bg-bg flex items-center justify-center">
        <div className="text-center">
          <p className="text-2xl text-muted mb-4">Loading...</p>
          <p className="text-sm text-muted">Please wait while we load your portfolios</p>
        </div>
      </div>
    )
  }
  
  // Safety check: if error, show error message
  if (error) {
    return (
      <div className="min-h-screen bg-bg flex items-center justify-center">
        <div className="text-center">
          <p className="text-2xl text-red-400 mb-4">Error: {error}</p>
          <button 
            onClick={() => {
              setError(null)
              setLoading(true)
              loadPortfolios().catch(err => {
                console.error('Error retrying:', err)
                setError(err.message || 'Failed to load')
                setLoading(false)
              })
            }}
            className="px-4 py-2 border border-accent text-accent rounded hover:bg-accent hover:text-black"
          >
            Retry
          </button>
        </div>
      </div>
    )
  }
  
  // Safety check: if no data and no portfolios, show message
  if (!loading && !data && portfolios.length === 0) {
    return (
      <div className="min-h-screen bg-bg flex items-center justify-center">
        <div className="text-center">
          <p className="text-2xl text-muted mb-4">No portfolios found</p>
          <a href="/import" className="px-4 py-2 border border-accent text-accent rounded hover:bg-accent hover:text-black">
            Create your first portfolio
          </a>
        </div>
      </div>
    )
  }

  // Calculate current portfolio safely
  const currentPortfolio = portfolios.find(p => p.id === selectedPortfolio) || null
  const totalValue = currentPortfolio?.positions?.reduce((sum: number, pos: any) => {
    try {
      const symbol = pos.asset?.symbol || ''
      const currentPrice = positionPrices[symbol] || pos.avg_price || 0
      return sum + ((pos.quantity || 0) * currentPrice)
    } catch (err) {
      console.error('Error calculating position value:', err)
      return sum
    }
  }, 0) || 0

  return (
    <div className="min-h-screen bg-bg">
      <nav className="sa-card mx-4 my-4 p-4 flex justify-between items-center">
        <h1 className="text-2xl font-bold" style={{ color: 'var(--accent)' }}>
          SAA Risk Analyzer
        </h1>
        <div className="flex gap-3">
          <a href="/import" className="px-4 py-2 border border-accent text-accent rounded hover:bg-accent hover:text-black">
            📂 Import Data
          </a>
          <a href="/analytics" className="px-4 py-2 border border-accent text-accent rounded hover:bg-accent hover:text-black">
            📊 Analytics
          </a>
          <div className="flex items-center gap-2">
            <button 
              onClick={recalculateRisk} 
              className="sa-btn px-4 py-2" 
              disabled={loading}
              title="Recalculate risk metrics based on current portfolio"
            >
              {loading ? 'Calculating...' : '🔄 Recalculate Risk'}
            </button>
            {data && (
              <span className="text-xs text-muted" title="Last calculated">
                Last: {new Date().toLocaleTimeString()}
              </span>
            )}
          </div>
        </div>
      </nav>

      {portfolios.length > 0 && (
        <div className="max-w-7xl mx-auto px-6 mb-4 flex items-center gap-4">
          <label className="text-sm text-muted">Portfolio:</label>
          <select 
            value={selectedPortfolio}
            onChange={(e) => setSelectedPortfolio(e.target.value)}
            className="px-4 py-2 bg-card border border-border rounded"
          >
            {portfolios.map(p => (
              <option key={p.id} value={p.id}>{p.name}</option>
            ))}
          </select>
          {currentPortfolio && (
            <div className="flex items-center gap-2">
              <button
                onClick={() => handleEditPortfolio(currentPortfolio)}
                className="px-3 py-1 text-sm border border-accent text-accent rounded hover:bg-accent hover:text-black transition-colors"
                title="Edit portfolio"
              >
                ✏️ Edit
              </button>
              <button
                onClick={() => {
                  const portfolioName = currentPortfolio.name || 'this portfolio'
                  if (confirm(`Are you sure you want to delete portfolio "${portfolioName}"?\n\nThis will permanently delete the portfolio and all its positions. This action cannot be undone.`)) {
                    handleDeletePortfolio(currentPortfolio.id)
                  }
                }}
                className="px-3 py-1 text-sm border border-red-500/50 text-red-400 rounded hover:bg-red-500/20 hover:border-red-500 transition-colors opacity-70 hover:opacity-100"
                title="Delete portfolio (destructive action)"
              >
                🗑️
              </button>
            </div>
          )}
        </div>
      )}

      <div className="max-w-7xl mx-auto p-6">
        {/* Portfolio Summary - показываем всегда если есть портфолио */}
        {currentPortfolio && (() => {
          const positions = currentPortfolio.positions || []
          if (positions.length === 0) {
            return (
              <div className="sa-card p-6 mb-6">
                <h2 className="text-xl font-semibold mb-4">Portfolio Summary</h2>
                <p className="text-muted">No positions in this portfolio. Add positions via Import Data.</p>
              </div>
            )
          }
          const totalPurchaseValue = currentPortfolio.positions.reduce((sum: number, pos: any) => 
            sum + (pos.quantity * pos.avg_price), 0)
          const totalCurrentValue = currentPortfolio.positions.reduce((sum: number, pos: any) => {
            const symbol = pos.asset?.symbol || ''
            const currentPrice = positionPrices[symbol] || pos.avg_price
            return sum + (pos.quantity * currentPrice)
          }, 0)
          const totalPnl = totalCurrentValue - totalPurchaseValue
          const totalPnlPercent = totalPurchaseValue > 0 ? (totalPnl / totalPurchaseValue) * 100 : 0
          
          return (
            <div className="sa-card p-6 mb-6">
              <h2 className="text-xl font-semibold mb-4">Portfolio Summary</h2>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <div className="text-sm text-muted mb-1">Total Purchase Value</div>
                  <div className="text-2xl font-bold">${totalPurchaseValue.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</div>
                </div>
                <div>
                  <div className="text-sm text-muted mb-1">Total Current Value</div>
                  <div className="text-2xl font-bold">${totalCurrentValue.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</div>
                </div>
                <div>
                  <div className="text-sm text-muted mb-1">Total P&L</div>
                  <div className={`text-2xl font-bold ${totalPnl >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                    {totalPnl >= 0 ? '+' : ''}${totalPnl.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    <span className="text-lg ml-2">({totalPnlPercent >= 0 ? '+' : ''}{totalPnlPercent.toFixed(2)}%)</span>
                  </div>
                </div>
              </div>
            </div>
          )
        })()}

        {/* Risk Metrics - показываем всегда */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <MetricCard 
            title="1-day VaR (99%)"
            tooltip="Value at Risk: Maximum potential loss over 1 day with 99% confidence"
            value={`$${Math.abs(data?.var_1d || 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`} 
            change={data?.var_1d && totalValue ? ((Math.abs(data.var_1d) / totalValue) * 100).toFixed(2) + '% of portfolio' : 'N/A'}
          />
          <MetricCard 
            title="1-day CVaR / ES (99%)"
            tooltip="Conditional Value at Risk (Expected Shortfall): Average loss beyond VaR threshold"
            value={`$${Math.abs(data?.cvar_1d || 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`} 
            change={data?.cvar_1d && totalValue ? ((Math.abs(data.cvar_1d) / totalValue) * 100).toFixed(2) + '% of portfolio' : 'N/A'}
          />
          <MetricCard 
            title="Portfolio Vol"
            tooltip="Annualized portfolio volatility based on historical returns"
            value={`${((data?.vol || 0) * 100).toFixed(1)}%`} 
            change="Annual"
          />
        </div>

        {/* Risk Contributors - показываем всегда */}
        <div className="sa-card p-6 mb-6">
          <h2 className="text-xl font-semibold mb-2">Risk Contributors (to portfolio value)</h2>
          <p className="text-sm text-muted mb-4">Contribution of each asset to total portfolio value</p>
          {data?.contributors && data.contributors.length > 0 ? (
            <div className="space-y-3">
              {data.contributors.map((item: any, idx: number) => (
                <RiskItem 
                  key={idx}
                  symbol={item.symbol} 
                  contribution={`${(item.contribution * 100).toFixed(1)}%`} 
                />
              ))}
            </div>
          ) : (
            <p className="text-muted text-sm">No portfolio data available. Create a portfolio and add positions to see weights.</p>
          )}
        </div>

        {/* Portfolio Positions - показываем всегда */}
        <div className="sa-card p-6 mt-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">Portfolio Positions</h2>
            {currentPortfolio && currentPortfolio.positions && currentPortfolio.positions.length > 0 && (
              <button
                onClick={() => loadPositionPrices(selectedPortfolio, true)}
                disabled={refreshingPrices}
                className="text-sm px-3 py-1 border border-accent text-accent rounded hover:bg-accent hover:text-black disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                title="Refresh current market prices for all positions"
              >
                {refreshingPrices ? (
                  <>
                    <span className="animate-spin">⏳</span>
                    <span>Refreshing...</span>
                  </>
                ) : (
                  <>
                    <span>🔄</span>
                    <span>Refresh Prices</span>
                  </>
                )}
              </button>
            )}
          </div>
          
          {portfolios.length === 0 ? (
            <div className="space-y-2">
              <p className="text-muted">No portfolios found.</p>
              <a href="/import" className="inline-block px-4 py-2 bg-accent text-black rounded hover:opacity-90">
                📂 Create Portfolio
              </a>
            </div>
          ) : !selectedPortfolio ? (
            <div className="space-y-2">
              <p className="text-muted">Please select a portfolio from the dropdown above.</p>
            </div>
          ) : !currentPortfolio ? (
            <p className="text-muted">Loading portfolio data...</p>
          ) : !currentPortfolio.positions || currentPortfolio.positions.length === 0 ? (
            <div className="space-y-2">
              <p className="text-muted">No positions in portfolio "{currentPortfolio.name}".</p>
              <a href="/import" className="inline-block px-4 py-2 bg-accent text-black rounded hover:opacity-90">
                📂 Add Positions
              </a>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-border">
                    <th className="text-left py-2 px-4">Name</th>
                    <th className="text-right py-2 px-4">Quantity</th>
                    <th className="text-right py-2 px-4">Purchase Price</th>
                    <th className="text-right py-2 px-4">Purchase Value</th>
                    <th className="text-right py-2 px-4">Current Price</th>
                    <th className="text-right py-2 px-4">Current Value</th>
                    <th className="text-right py-2 px-4">P&L</th>
                    <th className="text-right py-2 px-4">P&L %</th>
                    <th className="text-center py-2 px-4">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {currentPortfolio.positions.map((pos: any) => {
                    const symbol = pos.asset?.symbol || 'N/A'
                    const name = pos.asset?.name || symbol
                    const quantity = pos.quantity
                    const purchasePrice = pos.avg_price
                    const purchaseValue = quantity * purchasePrice
                    const currentPrice = positionPrices[symbol] || purchasePrice
                    const currentValue = quantity * currentPrice
                    const pnl = currentValue - purchaseValue
                    const pnlPercent = purchaseValue > 0 ? (pnl / purchaseValue) * 100 : 0
                    const isProfit = pnl >= 0
                    
                    return (
                      <tr key={pos.id} className="border-b border-border hover:bg-card/50">
                        <td className="py-2 px-4 font-semibold">{name}</td>
                        <td className="py-2 px-4 text-right">{quantity.toLocaleString(undefined, { maximumFractionDigits: 8 })}</td>
                        <td className="py-2 px-4 text-right">${purchasePrice.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                        <td className="py-2 px-4 text-right">${purchaseValue.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                        <td className="py-2 px-4 text-right">
                          <span className={positionPrices[symbol] ? 'text-accent font-semibold' : 'text-muted'}>
                            ${currentPrice.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                            {positionPrices[symbol] && <span className="ml-1 text-xs">🟢</span>}
                          </span>
                        </td>
                        <td className="py-2 px-4 text-right">${currentValue.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                        <td className={`py-2 px-4 text-right font-semibold ${isProfit ? 'text-green-400' : 'text-red-400'}`}>
                          {isProfit ? '+' : ''}${pnl.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                        </td>
                        <td className={`py-2 px-4 text-right font-semibold ${isProfit ? 'text-green-400' : 'text-red-400'}`}>
                          {isProfit ? '+' : ''}{pnlPercent.toFixed(2)}%
                        </td>
                        <td className="py-2 px-4 text-center">
                          <button
                            onClick={() => handleEditPosition(pos)}
                            className="text-accent hover:underline mr-3"
                            title="Edit position"
                          >
                            ✏️
                          </button>
                          <button
                            onClick={() => {
                              if (!pos.id) {
                                console.error('Position ID is missing:', pos)
                                alert('Error: Position ID is missing')
                                return
                              }
                              handleDeletePosition(pos.id)
                            }}
                            className="text-red-400 hover:underline"
                            title="Delete position"
                          >
                            🗑️
                          </button>
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
                <tfoot>
                  <tr className="border-t-2 border-accent font-bold bg-card/50">
                    <td className="py-3 px-4 font-bold">Total</td>
                    <td className="py-3 px-4 text-right">
                      {currentPortfolio.positions.reduce((sum: number, pos: any) => sum + pos.quantity, 0).toLocaleString('en-US', { maximumFractionDigits: 8 })}
                    </td>
                    <td className="py-3 px-4 text-right">—</td>
                    <td className="py-3 px-4 text-right">
                      ${currentPortfolio.positions.reduce((sum: number, pos: any) => sum + (pos.quantity * pos.avg_price), 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </td>
                    <td className="py-3 px-4 text-right">—</td>
                    <td className="py-3 px-4 text-right">
                      ${currentPortfolio.positions.reduce((sum: number, pos: any) => {
                        const symbol = pos.asset?.symbol || ''
                        const currentPrice = positionPrices[symbol] || pos.avg_price
                        return sum + (pos.quantity * currentPrice)
                      }, 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </td>
                    <td className={`py-3 px-4 text-right ${
                      (() => {
                        const totalPnl = currentPortfolio.positions.reduce((sum: number, pos: any) => {
                          const symbol = pos.asset?.symbol || ''
                          const currentPrice = positionPrices[symbol] || pos.avg_price
                          return sum + ((currentPrice - pos.avg_price) * pos.quantity)
                        }, 0)
                        return totalPnl >= 0 ? 'text-green-400' : 'text-red-400'
                      })()
                    }`}>
                      {(() => {
                        const totalPnl = currentPortfolio.positions.reduce((sum: number, pos: any) => {
                          const symbol = pos.asset?.symbol || ''
                          const currentPrice = positionPrices[symbol] || pos.avg_price
                          return sum + ((currentPrice - pos.avg_price) * pos.quantity)
                        }, 0)
                        return (totalPnl >= 0 ? '+' : '') + '$' + totalPnl.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
                      })()}
                    </td>
                    <td className={`py-3 px-4 text-right ${
                      (() => {
                        const totalPurchaseValue = currentPortfolio.positions.reduce((sum: number, pos: any) => sum + (pos.quantity * pos.avg_price), 0)
                        const totalCurrentValue = currentPortfolio.positions.reduce((sum: number, pos: any) => {
                          const symbol = pos.asset?.symbol || ''
                          const currentPrice = positionPrices[symbol] || pos.avg_price
                          return sum + (pos.quantity * currentPrice)
                        }, 0)
                        const totalPnlPercent = totalPurchaseValue > 0 ? ((totalCurrentValue - totalPurchaseValue) / totalPurchaseValue) * 100 : 0
                        return totalPnlPercent >= 0 ? 'text-green-400' : 'text-red-400'
                      })()
                    }`}>
                      {(() => {
                        const totalPurchaseValue = currentPortfolio.positions.reduce((sum: number, pos: any) => sum + (pos.quantity * pos.avg_price), 0)
                        const totalCurrentValue = currentPortfolio.positions.reduce((sum: number, pos: any) => {
                          const symbol = pos.asset?.symbol || ''
                          const currentPrice = positionPrices[symbol] || pos.avg_price
                          return sum + (pos.quantity * currentPrice)
                        }, 0)
                        const totalPnlPercent = totalPurchaseValue > 0 ? ((totalCurrentValue - totalPurchaseValue) / totalPurchaseValue) * 100 : 0
                        return (totalPnlPercent >= 0 ? '+' : '') + totalPnlPercent.toFixed(2) + '%'
                      })()}
                    </td>
                    <td className="py-3 px-4"></td>
                  </tr>
                </tfoot>
              </table>
            </div>
          )}
        </div>

        {/* Edit Portfolio Modal */}
        {editingPortfolio && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="sa-card p-6 max-w-md w-full mx-4">
              <h2 className="text-xl font-semibold mb-4">Edit Portfolio</h2>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Name</label>
                  <input
                    type="text"
                    value={editName}
                    onChange={(e) => setEditName(e.target.value)}
                    className="w-full px-3 py-2 bg-bg border border-border rounded"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Description</label>
                  <textarea
                    value={editDescription}
                    onChange={(e) => setEditDescription(e.target.value)}
                    className="w-full px-3 py-2 bg-bg border border-border rounded"
                    rows={3}
                  />
                </div>
                <div className="flex gap-3">
                  <button
                    onClick={handleSavePortfolio}
                    className="flex-1 bg-accent text-black py-2 px-4 rounded hover:opacity-90"
                  >
                    Save
                  </button>
                  <button
                    onClick={() => setEditingPortfolio(null)}
                    className="flex-1 bg-border text-fg py-2 px-4 rounded hover:opacity-90"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Edit Position Modal */}
        {editingPosition && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="sa-card p-6 max-w-md w-full mx-4">
              <h2 className="text-xl font-semibold mb-4">Edit Position</h2>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Symbol</label>
                  <input
                    type="text"
                    value={editPositionSymbol}
                    onChange={async (e) => {
                      const symbol = e.target.value.toUpperCase()
                      setEditPositionSymbol(symbol)
                      if (symbol.length > 0) {
                        try {
                          const res = await api.get(`/market/price/${symbol}`)
                          if (res.data && res.data.price) {
                            setEditPositionPrice(parseFloat(res.data.price.toFixed(2)))
                          }
                        } catch (err) {
                          console.log(`Could not fetch price for ${symbol}`)
                        }
                      }
                    }}
                    className="w-full px-3 py-2 bg-bg border border-border rounded"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Quantity</label>
                  <input
                    type="number"
                    step="0.001"
                    value={editPositionQuantity}
                    onChange={(e) => setEditPositionQuantity(parseFloat(e.target.value) || 0)}
                    className="w-full px-3 py-2 bg-bg border border-border rounded"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Avg Price</label>
                  <input
                    type="number"
                    step="0.01"
                    value={editPositionPrice}
                    onChange={(e) => setEditPositionPrice(parseFloat(e.target.value) || 0)}
                    className="w-full px-3 py-2 bg-bg border border-border rounded"
                  />
                </div>
                <div className="flex gap-3">
                  <button
                    onClick={handleSavePosition}
                    className="flex-1 bg-accent text-black py-2 px-4 rounded hover:opacity-90"
                  >
                    Save
                  </button>
                  <button
                    onClick={() => setEditingPosition(null)}
                    className="flex-1 bg-border text-fg py-2 px-4 rounded hover:opacity-90"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

function MetricCard({ title, value, change, tooltip }: { title: string; value: string; change: string; tooltip?: string }) {
  return (
    <div className="sa-card p-6 relative group">
      <h3 className="text-sm text-muted mb-2 flex items-center gap-1">
        {title}
        {tooltip && (
          <span className="text-xs cursor-help" title={tooltip}>ℹ️</span>
        )}
      </h3>
      <p className="text-3xl font-bold mb-2">{value}</p>
      <p className="text-sm text-muted">{change}</p>
    </div>
  )
}

function RiskItem({ symbol, contribution }: { symbol: string; contribution: string }) {
  const contributionNum = parseFloat(contribution.replace('%', ''))
  return (
    <div className="flex justify-between items-center">
      <div className="flex items-center gap-2 flex-1">
        <span className="font-semibold min-w-[60px]">{symbol}</span>
        <div className="flex-1 bg-border rounded-full h-3 relative">
          <div
            className="h-full rounded-full transition-all"
            style={{
              width: contribution,
              background: 'var(--accent)',
            }}
          />
          <span className="absolute inset-0 flex items-center justify-center text-xs font-medium text-fg">
            {contributionNum > 5 && contribution}
          </span>
        </div>
      </div>
      <span className="text-muted w-16 text-right font-semibold">{contribution}</span>
    </div>
  )
}
