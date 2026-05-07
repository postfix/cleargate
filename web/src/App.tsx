import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import { Terminal } from 'lucide-react'
import CatalogPage from './pages/CatalogPage'
import ExecutionPage from './pages/ExecutionPage'
import './index.css'
import './App.css'

function App() {
  return (
    <BrowserRouter>
      <div className="app-container">
        <header className="app-header">
          <Link to="/" className="logo-link">
            <Terminal size={24} color="var(--color-accent)" />
            <span className="display" style={{ fontSize: '20px' }}>ClearGate</span>
          </Link>
          <nav className="nav-tabs">
            <Link to="/" className="nav-tab active">Catalog</Link>
          </nav>
        </header>
        
        <main className="app-content">
          <Routes>
            <Route path="/" element={<CatalogPage />} />
            <Route path="/tool/:id" element={<ExecutionPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  )
}

export default App
