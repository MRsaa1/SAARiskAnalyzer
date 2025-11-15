import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { useEffect } from 'react'
import DashboardPage from './pages/Dashboard'
import PortfolioImport from './pages/PortfolioImport'
import AnalyticsPage from './pages/Analytics'

function App() {
  useEffect(() => {
    console.log('🚀 App mounted')
    console.log('📡 API Base URL:', import.meta.env.VITE_API_BASE_URL)
    console.log('🏗️ Mode:', import.meta.env.MODE)
  }, [])

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<DashboardPage />} />
        <Route path="/import" element={<PortfolioImport />} />
        <Route path="/analytics" element={<AnalyticsPage />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
