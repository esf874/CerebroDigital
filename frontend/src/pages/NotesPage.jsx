import { FiInbox } from 'react-icons/fi'
import { useNavigate } from 'react-router-dom'
import MarkdownWithLinks from '../components/MarkdownWithLinks'
import NoteForm from '../components/NoteForm'

// Página principal que lista todas las notas del usuario.
function NotesPage({ notes, loading, onRefreshNotes }) {
  const navigate = useNavigate()

  const handleNoteClick = (noteId) => {
    navigate(`/note/${noteId}`)
  }

  return (
    <div>
      <div className="page-header" style={{ textAlign: 'center' }}>
        <h1 className="page-title" style={{ textAlign: 'center' }}>Mis Notas</h1>
        <p className="page-subtitle" style={{ textAlign: 'center' }}>Rellena tu base de conocimiento</p>
      </div>

      <div style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: '2rem' }}>
        <NoteForm onNoteCreated={onRefreshNotes} />
      </div>

      {loading ? (
        <div className="empty-state">
          <div className="empty-state-icon">⏳</div>
          <div className="empty-state-title">Cargando notas...</div>
        </div>
      ) : notes.length === 0 ? (
        <div className="empty-state">
          <FiInbox size={48} className="empty-state-icon" strokeWidth={1} />
          <div className="empty-state-title">No hay notas aún</div>
          <div className="empty-state-text">Crea tu primera nota para empezar</div>
        </div>
      ) : (
        <div className="notes-grid">
          {notes.map(note => (
            <div key={note.id} className="note-card" onClick={() => handleNoteClick(note.id)}>

              <div className="note-card-title">{note.title}</div>

              <div style={{ display: 'flex', gap: '6px', alignItems: 'center' }}>
                {note.status === 'finished' && <span style={{ color: '#0dae78' }}>●</span>}
                {note.status === 'in_progress' && <span style={{ color: '#60f0c0' }}>◐</span>}
                {note.status === 'pending' && <span style={{ color: '#a2d9c7' }}>○</span>}

                {note.priority === 'low' && <span style={{ color: "#f2b2d5" }}>●</span>}
                {note.priority === 'medium' && <span style={{ color: "#df82d2" }}>●●</span>}
                {note.priority === 'high' && <span style={{ color: "#de3ac8" }}>●●●</span>}
              </div>


              <div className="note-card-content">
                <MarkdownWithLinks
                  content={
                    note.content && note.content.length > 120
                      ? note.content.substring(0, 120) + '...'
                      : note.content || 'Sin contenido'
                  }
                  onNavigateToNote={handleNoteClick}
                />
              </div>

              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '8px' }}>
                <div className="note-card-tags">
                  {note.tags && note.tags.slice(0, 2).map(tag => (
                    <span key={tag} className="tag">#{tag}</span>
                  ))}
                </div>
                <span style={{ fontSize: '10px', color: 'var(--text-muted)' }}>
                  {note.updatedAt ? new Date(note.updatedAt).toLocaleDateString() : 'Fecha desconocida'}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default NotesPage