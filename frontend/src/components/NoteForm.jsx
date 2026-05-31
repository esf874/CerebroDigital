import { useState } from 'react'
import { api } from '../services/api'
import MarkdownEditor from './MarkdownEditor'
import Modal from './Modal'
import { NotePrioritySelector, NoteStatusSelector } from './NoteStatusBadge'
import PrimaryButton from './ui/PrimaryButton'

// Manejo creación nuevas notas.
function NoteForm({ onNoteCreated }) {
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [tagsInput, setTagsInput] = useState('')
  const [status, setStatus] = useState('pending')
  const [priority, setPriority] = useState('medium')
  const [loading, setLoading] = useState(false)
  const [isOpen, setIsOpen] = useState(false)
  const [titleError, setTitleError] = useState('')
  const [detectedHashtags, setDetectedHashtags] = useState([])

  // Combinación tags manuales y detectados
  const getCombinedTags = (manualTagsStr, detectedTags) => {
    const manualTags = manualTagsStr.split(',').map(t => t.trim()).filter(t => t)
    const allTags = [...new Set([...manualTags, ...detectedTags])]
    return allTags.join(', ')
  }

  const handleTagsDetected = (hashtags) => {
    setDetectedHashtags(hashtags)
    const combined = getCombinedTags(tagsInput, hashtags)
    setTagsInput(combined)
  }

  const handleTagsChange = (e) => {
    const newValue = e.target.value
    setTagsInput(newValue)
  }


  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!title.trim()) {
      setTitleError(' El título es obligatorio')
      return
    }

    setLoading(true)
    try {
      // Combinar todas las etiquetas
      const manualTags = tagsInput.split(',').map(t => t.trim()).filter(t => t)
      const allTags = [...new Set([...manualTags, ...detectedHashtags])]

      await api.createNote({
        title: title.trim(),
        content: content,
        tags: allTags,
        status: status,
        priority: priority
      })

      handleClose() // Función de limpiar 
      if (onNoteCreated) onNoteCreated()
    } catch (error) {
      console.error('Error al crear nota:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleClose = () => {
    setIsOpen(false)
    setTitle('')
    setContent('')
    setTagsInput('')
    setDetectedHashtags([])
    setStatus('pending')
    setPriority('medium')
  }

  return (
    <>
      <PrimaryButton
        onClick={() => setIsOpen(true)}
        style={{
          padding: '8px 16px'
        }}
      >
        + Nueva Nota
      </PrimaryButton>

      <Modal isOpen={isOpen} onClose={handleClose} title="Crear Nueva Nota">
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="Título *"
            value={title}
            onChange={(e) => {
              setTitle(e.target.value)
              if (titleError) setTitleError('')
            }}
          />

          {titleError && (
            <div style={{
              color: '#e66bd5',
              fontSize: '0.75rem',
              marginTop: '0.2rem',
              marginBottom: '0.8rem'
            }}>
              {titleError}
            </div>
          )}

          <MarkdownEditor
            value={content}
            onChange={setContent}
            placeholder="Contenido en Markdown... (usa #etiqueta para crear etiquetas)"
            rows={8}
            onHashtagDetected={handleTagsDetected}
          />

          {/* Mostrar hashtags detectados */}
          {detectedHashtags.length > 0 && (
            <div style={{ marginTop: '0.5rem', marginBottom: '0.5rem' }}>
              <span style={{ fontSize: '0.7rem', color: 'var(--text-muted)' }}>🔍 Hashtags detectados: </span>
              {detectedHashtags.map(tag => (
                <span key={tag} className="tag" style={{ marginLeft: '4px' }}>#{tag}</span>
              ))}
            </div>
          )}

          {/* Selectores de estado y prioridad */}
          <div style={{ marginTop: '1rem', marginBottom: '1rem' }}>
            <label style={{ display: 'block', marginBottom: '0.5rem', color: 'var(--text-secondary)', fontSize: '0.75rem' }}>
              Estado
            </label>
            <NoteStatusSelector value={status} onChange={setStatus} />
          </div>

          <div style={{ marginBottom: '1rem' }}>
            <label style={{ display: 'block', marginBottom: '0.5rem', color: 'var(--text-secondary)', fontSize: '0.75rem' }}>
              Prioridad
            </label>
            <NotePrioritySelector value={priority} onChange={setPriority} />
          </div>

          {/* Campo de etiquetas (editable) */}
          <div>
            <label style={{ display: 'block', marginBottom: '0.5rem', color: 'var(--text-secondary)', fontSize: '0.75rem' }}>
              Etiquetas (separadas por comas)
            </label>
            <input
              type="text"
              placeholder="ej: receta, importante, trabajo"
              value={tagsInput}
              onChange={handleTagsChange}
              style={{
                width: '100%',
                marginBottom: '0.5rem',
                backgroundColor: 'var(--bg-tertiary)',
                color: 'var(--text-primary)',
                padding: '10px',
                borderRadius: '8px',
                border: '1px solid var(--border)'
              }}
            />
            <small style={{ color: 'var(--text-muted)', fontSize: '0.7rem' }}>
              Ejemplo: personal, trabajo, ideas (usa comas para separar)
            </small>
          </div>

          <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
            <button type="button" className="btn btn-secondary" onClick={handleClose}>
              Cancelar
            </button>

            <PrimaryButton type="submit" disabled={loading}>
              {loading ? 'Guardando...' : 'Crear Nota'}
            </PrimaryButton>
          </div>
        </form>
      </Modal>
    </>
  )
}

export default NoteForm
