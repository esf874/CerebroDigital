import { useEffect, useState } from 'react'
import { FiBook, FiMessageCircle } from 'react-icons/fi'
import { Link, useLocation } from 'react-router-dom'
import './App.css'
import AppRoutes from './routes'
import { api } from './services/api'

// Componente principal de la aplicación, maneja la navegación.
function App() {
  const [notes, setNotes] = useState([])
  const [loadingNotes, setLoadingNotes] = useState(true)
  const location = useLocation()

  const loadNotes = async () => {
    try {
      const data = await api.getNotes()
      setNotes(data)
    } catch (error) {
      console.error('Error al cargar notas:', error)
    } finally {
      setLoadingNotes(false)
    }
  }

  useEffect(() => {
    loadNotes()
  }, [])

  const refreshNotes = async () => {
    await loadNotes()
  }

  // Determinar qué vista está activa según la URL
  const isNotesActive = location.pathname === '/' || location.pathname.startsWith('/note/')
  const isAskActive = location.pathname === '/ask'

  return (
    <div className="app-container">
      <nav className="navbar">
        <div className="navbar-content">
          {/* Logo clicable hacia página inicial */}
          <Link to="/" className="logo" style={{ textDecoration: 'none', cursor: 'pointer' }}>
            <svg
              className="logo-icon"
              width="28"
              height="28"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M12 21c-2.5 0-4.5-1.5-5.5-3-2-2.5-2-6 0-8.5C7.5 8 8.5 7.5 9 7c.5-2 2-3 3-3s2.5 1 3 3c.5.5 1.5 1 2.5 2.5 2 2.5 2 6 0 8.5-1 1.5-3 3-5.5 3z" />
              <path d="M10 9c-1 0-2 .5-2.5 1.5.5.5 1 1.5.5 2.5-.5.5-1.5.5-2 0" />
              <path d="M9.5 14c-1 0-1.5 1-1 2s1.5 1 2 0" />
              <path d="M14 9c1 0 2 .5 2.5 1.5-.5.5-1 1.5-.5 2.5.5.5 1.5.5 2 0" />
              <path d="M14.5 14c1 0 1.5 1 1 2s-1.5 1-2 0" />
              <path d="M12 18v3" />
              <path d="M11 21h2" />
            </svg>
            <span className="logo-text">Cerebro Digital</span>
          </Link>

          <div className="nav-links">
            <Link
              to="/"
              className={`nav-link ${isNotesActive ? 'active' : ''}`}
              style={{ textDecoration: 'none' }}
            >
              <FiBook size={18} />
              <span>Mis Notas</span>
            </Link>
            <Link
              to="/ask"
              className={`nav-link ${isAskActive ? 'active' : ''}`}
              style={{ textDecoration: 'none' }}
            >
              <FiMessageCircle size={18} />
              <span>Asistente</span>
            </Link>
          </div>
        </div>
      </nav>

      <main className="main-content page-enter">
        <AppRoutes
          notes={notes}
          loadingNotes={loadingNotes}
          onRefreshNotes={refreshNotes}
        />
      </main>
    </div>
  )
}

export default App