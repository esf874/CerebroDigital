import { Navigate, Route, Routes } from 'react-router-dom'
import AskPage from './pages/AskPage'
import NoteDetailPage from './pages/NoteDetailPage'
import NotesPage from './pages/NotesPage'

// Define las rutas de la aplicación y qué componente renderizar para cada una.
function AppRoutes({ notes, loadingNotes, onRefreshNotes }) {
  return (
    <Routes>
      {/* Pagina principal/home: lista de notas */}
      <Route 
        path="/" 
        element={
          <NotesPage 
            notes={notes}
            loading={loadingNotes}
            onRefreshNotes={onRefreshNotes}
          />
        } 
      />
      
      {/* Detalle de una nota */}
      <Route 
        path="/note/:id" 
        element={
          <NoteDetailPage 
            onRefreshNotes={onRefreshNotes}
          />
        } 
      />
      
      {/* asistente RAG*/}
      <Route path="/ask" element={<AskPage />} />
      
      {/* Redireccion por defecto */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default AppRoutes