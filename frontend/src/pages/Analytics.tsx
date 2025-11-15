import { useEffect, useState } from 'react'
import api from '../lib/api'
import Heatmap from '../components/Heatmap'
import RollingVaRChart from '../components/RollingVaRChart'
import PnLHistogram from '../components/PnLHistogram'
import PCAScreePlot from '../components/PCAScreePlot'

export default function AnalyticsPage() {
  
  const [correlation, setCorrelation] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  
  const symbols = ['BTC', 'ETH', 'GLD', 'TLT', 'SPY']
  
  const loadCorrelations = async () => {
    setLoading(true)
    try {
      const res = await api.post('/risk/correlation', {
        symbols,
        window_days: 10, // У нас только 5 дней данных, используем минимальное окно
      })
      setCorrelation(res.data)
      console.log('✅ Correlation loaded:', res.data)
    } catch (err: any) {
      console.error('Correlation error:', err.response?.data || err.message)
      // Show user-friendly error
      alert('Not enough price data for correlation analysis. Need at least 10 days of historical prices.')
    }
    setLoading(false)
  }
  
  useEffect(() => {
    loadCorrelations()
  }, [])
  
  const rollingVaRData = Array.from({ length: 30 }, (_, i) => ({
    date: `Day ${i + 1}`,
    var: 100000 + Math.random() * 50000,
    cvar: 150000 + Math.random() * 60000,
  }))
  
  const pnlHistData = Array.from({ length: 20 }, (_, i) => ({
    bucket: `${(i - 10) * 10}k`,
    count: Math.floor(Math.random() * 100),
  }))
  
  const pcaData = Array.from({ length: 5 }, (_, i) => ({
    component: `PC${i + 1}`,
    variance: Math.max(0.5 - i * 0.1, 0.05),
    cumulative: Math.min(0.5 + i * 0.15, 0.95),
  }))
  
  return (
    <div className="min-h-screen bg-bg">
      <nav className="sa-card mx-4 my-4 p-4 flex justify-between items-center">
        <h1 className="text-2xl font-bold" style={{ color: 'var(--accent)' }}>
          SAA Risk Analyzer — Analytics
        </h1>
        <a href="/" className="px-4 py-2 border border-accent text-accent rounded hover:bg-accent hover:text-black">
          ← Dashboard
        </a>
      </nav>

      <div className="max-w-7xl mx-auto p-6 space-y-6">
        <div className="sa-card p-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">Correlation Matrix</h2>
            <button onClick={loadCorrelations} className="sa-btn" disabled={loading}>
              {loading ? 'Loading...' : '🔄 Refresh'}
            </button>
          </div>
          
          {correlation?.matrix ? (
            <Heatmap data={correlation.matrix} labels={symbols} />
          ) : (
            <Heatmap 
              data={[
                [1.0, 0.65, 0.23, 0.15, -0.12],
                [0.65, 1.0, 0.34, 0.21, -0.08],
                [0.23, 0.34, 1.0, 0.12, 0.05],
                [0.15, 0.21, 0.12, 1.0, 0.18],
                [-0.12, -0.08, 0.05, 0.18, 1.0],
              ]}
              labels={symbols}
            />
          )}
        </div>

        <div className="sa-card p-6">
          <h2 className="text-xl font-semibold mb-4">Rolling VaR & CVaR (30 days)</h2>
          <RollingVaRChart data={rollingVaRData} />
        </div>

        <div className="sa-card p-6">
          <h2 className="text-xl font-semibold mb-4">P&L Distribution</h2>
          <PnLHistogram data={pnlHistData} var={-80000} cvar={-120000} />
        </div>

        <div className="sa-card p-6">
          <h2 className="text-xl font-semibold mb-4">PCA — Principal Components</h2>
          <PCAScreePlot data={pcaData} />
        </div>
      </div>
    </div>
  )
}
