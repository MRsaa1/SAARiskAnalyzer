import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './styles/theme.css'
import './styles/global.css'

console.log('🔥 main.tsx loaded')

try {
  const rootElement = document.getElementById('root')
  if (!rootElement) {
    console.error('❌ Root element not found!')
    throw new Error('Root element not found')
  }
  
  console.log('✅ Root element found')
  const root = ReactDOM.createRoot(rootElement)
  console.log('✅ React root created')
  
  root.render(
    <React.StrictMode>
      <App />
    </React.StrictMode>
  )
  
  console.log('✅ React rendered')
} catch (error) {
  console.error('❌ Error rendering React:', error)
  document.body.innerHTML = `
    <div style="padding: 20px; color: red;">
      <h1>Error loading application</h1>
      <p>${error instanceof Error ? error.message : String(error)}</p>
      <pre>${error instanceof Error ? error.stack : String(error)}</pre>
    </div>
  `
}
