const API_BASE = '/api'

// Funciones para centralizar llamadas al backend.
export const api = {
  getNotes: async () => {
    const res = await fetch(`${API_BASE}/notes`)
    if (!res.ok) throw new Error('Error al cargar notas')
    return res.json()
  },

  getNote: async (id) => {
    const res = await fetch(`${API_BASE}/notes/${id}`)
    if (!res.ok) throw new Error('Error al cargar la nota')
    return res.json()
  },

  createNote: async (note) => {
    const res = await fetch(`${API_BASE}/notes`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(note)
    })
    if (!res.ok) throw new Error('Error al crear nota')
    return res.json()
  },

  updateNote: async (id, note) => {
    const res = await fetch(`${API_BASE}/notes/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(note)
    })
    if (!res.ok) throw new Error('Error al actualizar nota')
    return res.json()
  },

  deleteNote: async (id) => {
    const res = await fetch(`${API_BASE}/notes/${id}`, {
      method: 'DELETE'
    })
    if (!res.ok) throw new Error('Error al eliminar nota')

    const text = await res.text()
    if (!text) return { success: true }

    try {
      return JSON.parse(text)
    } catch {
      return { success: true }
    }
  },

  getGraph: async (noteId, depth = 2, limit = 20) => {
    const res = await fetch(`${API_BASE}/notes/${noteId}/graph?depth=${depth}&limit=${limit}`)
    if (!res.ok) throw new Error('Error al cargar el grafo')
    return res.json()
  },

  createLink: async (originNoteId, targetNoteId, alias) => {
    const res = await fetch(`${API_BASE}/links`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        originNoteId: originNoteId,
        targetNoteId: targetNoteId,
        alias: alias
      })
    })
    if (!res.ok) throw new Error('Error al crear enlace')
    return res.json()
  },

  getNoteByTitle: async (title) => {
    const url = `${API_BASE}/notes/lookup?title=${encodeURIComponent(title.trim())}`;
    const res = await fetch(url);
    console.log("url busqueda:", url);

    if (!res.ok) throw new Error('Nota no encontrada');
    return res.json();
  },

  ask: async (question, currentNoteId = null) => {
    const res = await fetch(`${API_BASE}/ask`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ question, currentNoteId: currentNoteId || undefined })
    })
    if (!res.ok) throw new Error('Error al consultar al asistente')
    return res.json()
  }
}