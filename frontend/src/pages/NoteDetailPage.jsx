import { Share2 } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import MarkdownEditor from '../components/MarkdownEditor'
import MarkdownWithLinks from '../components/MarkdownWithLinks'
import NoteGraph from '../components/NoteGraph'
import { NotePriorityBadge, NotePrioritySelector, NoteStatusBadge, NoteStatusSelector } from '../components/NoteStatusBadge'
import AnswerCard from '../components/ui/AnswerCard'
import PrimaryButton from '../components/ui/PrimaryButton'
import { api } from '../services/api'

// Página de detalle de una nota, con edición, eliminación, grafo y sección de pregunta contextual.
function NoteDetailPage({ onRefreshNotes }) {
  const { id } = useParams()
  const navigate = useNavigate()
  const [note, setNote] = useState(null)
  const [loading, setLoading] = useState(true)
  const [isEditing, setIsEditing] = useState(false)
  const [editTitle, setEditTitle] = useState('')
  const [editContent, setEditContent] = useState('')
  const [editTags, setEditTags] = useState('')
  const [question, setQuestion] = useState('')
  const [answer, setAnswer] = useState('')
  const [asking, setAsking] = useState(false)

  const loadNote = async () => {
    try {
      const data = await api.getNote(id)
      setNote(data)
      setEditTitle(data.title)
      setEditContent(data.content)
      setEditTags(data.tags?.join(', ') || '')
    } catch (error) {
      console.error('Error al cargar nota:', error)
      navigate('/') // Si no existe, volver al inicio
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadNote()
  }, [id])

  const handleUpdate = async () => {
    try {
      await api.updateNote(id, {
        title: editTitle,
        content: editContent,
        tags: editTags.split(',').map(t => t.trim()).filter(t => t),
        status: note.status,
        priority: note.priority
      })

      // Recargar la nota actual (para mostrar cambios)
      await loadNote()
      setIsEditing(false)

      if (onRefreshNotes) await onRefreshNotes();

    } catch (error) {
      console.error('Error al actualizar:', error)
    }
  }

  const handleDelete = async () => {
    if (window.confirm('¿Eliminar esta nota permanentemente?')) {
      try {
        await api.deleteNote(id)

        // Refrescar la lista 
        await onRefreshNotes()

        // Volver a lista de notas
        setTimeout(() => {
          navigate('/')
        }, 100)

      } catch (error) {
        console.error('Error al eliminar:', error)
      }
    }
  }

  const handleAsk = async () => {
    if (!question.trim()) return
    setAsking(true)
    setAnswer('')
    try {
      const data = await api.ask(question, id)
      setAnswer(data.answer)
    } catch (error) {
      setAnswer('Error al consultar al asistente')
    } finally {
      setAsking(false)
    }
  }

  if (loading) return <div className="empty-state">Cargando nota...</div>
  if (!note) return <div className="empty-state">Nota no encontrada</div>

  return (
    <div>
      {isEditing ? (
        <div className="edit-section card">
          {/* Campo de título */}
          <input
            value={editTitle}
            onChange={(e) => setEditTitle(e.target.value)}
            className="title-edit-input"
            style={{ width: '100%', fontSize: '1.5rem', fontWeight: 600, marginBottom: '1rem', backgroundColor: 'var(--bg-tertiary)', color: 'var(--text-primary)', border: '1px solid var(--border)', padding: '8px', borderRadius: '8px' }}
          />

          {/* Selectores de estado y prioridad */}
          <div style={{ marginBottom: '1rem' }}>
            <label style={{ display: 'block', marginBottom: '0.5rem', color: 'var(--text-secondary)', fontSize: '0.75rem', textTransform: 'uppercase', letterSpacing: '1px' }}>
              Estado
            </label>
            <NoteStatusSelector
              value={note.status}
              onChange={(val) => setNote({ ...note, status: val })}
            />
          </div>

          <div style={{ marginBottom: '1rem' }}>
            <label style={{ display: 'block', marginBottom: '0.5rem', color: 'var(--text-secondary)', fontSize: '0.75rem', textTransform: 'uppercase', letterSpacing: '1px' }}>
              Prioridad
            </label>
            <NotePrioritySelector
              value={note.priority}
              onChange={(val) => setNote({ ...note, priority: val })}
            />
          </div>

          {/* Editor con detección de #etiquetas y [[links]] */}
          <MarkdownEditor
            value={editContent}
            onChange={setEditContent}
            rows={12}
            onHashtagDetected={(hashtags) => {
              const currentTags = editTags.split(',').map(t => t.trim()).filter(t => t)
              const combined = [...new Set([...currentTags, ...hashtags])]
              setEditTags(combined.join(', '))
            }}
          />

          {/* Edicion etiquetas manual */}
          <div style={{ marginTop: '1rem' }}>
            <label style={{ display: 'block', marginBottom: '0.5rem', color: 'var(--text-secondary)', fontSize: '0.75rem' }}>
              Etiquetas (separadas por comas)
            </label>
            <input
              value={editTags}
              onChange={(e) => setEditTags(e.target.value)}
              placeholder="ej: receta, importante, trabajo"
              style={{
                width: '100%',
                marginBottom: '1rem',
                backgroundColor: 'var(--bg-tertiary)',
                color: 'var(--text-primary)',
                padding: '10px',
                borderRadius: '8px',
                border: '1px solid var(--border)'
              }}
            />
          </div>

          <div style={{ display: 'flex', gap: '1rem' }}>
            <PrimaryButton
              onClick={handleUpdate}
              style={{
                padding: '10px 20px'
              }}
            >
              Guardar Cambios
            </PrimaryButton>

            <button className="btn btn-secondary" onClick={() => setIsEditing(false)}>Cancelar</button>
          </div>
        </div>
      ) : (
        <>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '1.5rem' }}>
            <h1 style={{ fontSize: '2rem', fontWeight: 600 }}>{note.title}</h1>

            <div style={{ display: 'flex', gap: '0.75rem' }}>
              <button
                className="btn-edit"
                onClick={() => setIsEditing(true)}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '8px',
                  background: 'transparent',
                  border: '1px solid rgba(255,255,255,0.2)',
                  color: 'white',
                  padding: '8px 16px',
                  borderRadius: '8px',
                  cursor: 'pointer',
                  transition: 'all 0.2s'
                }}
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                </svg>
                <span>Editar</span>
              </button>

              <PrimaryButton
                onClick={handleDelete}
                style={{
                  backgroundColor: '#d91cc0',
                  padding: '8px 16px'
                }}
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <polyline points="3 6 5 6 21 6"></polyline>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  <line x1="10" y1="11" x2="10" y2="17"></line>
                  <line x1="14" y1="11" x2="14" y2="17"></line>
                </svg>
                <span>Eliminar</span>
              </PrimaryButton>
            </div>
          </div>

          {/* Badges de estado y prioridad */}
          <div style={{ display: 'flex', gap: '0.75rem', marginBottom: '1rem', flexWrap: 'wrap' }}>
            <NoteStatusBadge status={note.status} />
            <NotePriorityBadge priority={note.priority} />
          </div>

          {/* Fecha y Etiquetas */}
          <div style={{ marginBottom: '1.5rem', fontSize: '0.8rem', color: 'var(--text-muted)' }}>
            📅 Última actualización: {new Date(note.updatedAt).toLocaleString()}
          </div>


          {note.tags && note.tags.length > 0 && (
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem', marginBottom: '1.5rem' }}>
              {note.tags.map(tag => (
                <span key={tag} className="tag">#{tag}</span>
              ))}
            </div>
          )}

          {/* Renderiza contenido con soporte para links */}
          <div className="card markdown-content" style={{ lineHeight: 1.7, marginBottom: '2rem' }}>
            <MarkdownWithLinks
              content={note.content}
              onNavigateToNote={(targetId) =>
                navigate(`/note/${targetId}`)}
            />
          </div>

          {/* Subgrafo de conexiones*/}
          <div className="card" style={{ marginBottom: '2rem' }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '12px',
              marginBottom: '1.5rem'
            }}>
              <div style={{
                backgroundColor: 'rgba(255, 21, 193, 0.1)',
                padding: '8px',
                borderRadius: '50%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
              }}>
                <Share2 size={20} color="#ffffff" strokeWidth={1.5} />
              </div>
              <h3 style={{ margin: 0, fontSize: '1.2rem', color: '#ffffff' }}>
                Conexiones de esta nota
              </h3>
            </div>
            <NoteGraph
              noteId={id}
              onNodeClick={(clickedId) => {
                if (clickedId !== id) {
                  navigate(`/note/${clickedId}`)
                }
              }}
              depth={2}
              limit={15}
            />
          </div>

          {/* Sección de pregunta contextual */}
          <div className="card" style={{ backgroundColor: 'rgba(99, 102, 241, 0.05)' }}>
            <h3 style={{ marginBottom: '1rem' }}> Pregunta sobre esta nota</h3>
            <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
              <input
                value={question}
                onChange={(e) => setQuestion(e.target.value)}
                placeholder="Pregunta en base a la nota. "
                style={{ flex: 1, minWidth: '200px' }}
                onKeyPress={(e) => e.key === 'Enter' && handleAsk()}
              />

              <PrimaryButton
                type="submit"
                disabled={loading}
                style={{
                  padding: '8px 20px'
                }}
              >
                {loading ? 'Pensando...' : 'Preguntar'}
              </PrimaryButton>

            </div>
            {answer && <AnswerCard answer={answer} />}
          </div>
        </>
      )}
    </div>
  )
}

export default NoteDetailPage
