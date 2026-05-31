import { useState } from 'react'

// Componente de editor de texto con vista previa y detección de etiquetas.
function MarkdownEditor({ value, onChange, placeholder, rows = 10, onHashtagDetected }) {
  const [previewMode, setPreviewMode] = useState(false)

  const escapeHtml = (str) => {
    if (!str) return ''
    return str.replace(/[&<>]/g, function (m) {
      if (m === '&') return '&amp;'
      if (m === '<') return '&lt;'
      if (m === '>') return '&gt;'
      return m
    })
  }

  const extractCompleteHashtags = (text) => {
    const hashtagRegex = /#([\p{L}0-9_-]+)(?=\s)/gu
    const matches = []
    let match

    while ((match = hashtagRegex.exec(text)) !== null) {
      matches.push(match[1])
    }

    return [...new Set(matches)]
  }

  const handleChange = (e) => {
    const newValue = e.target.value
    onChange(newValue)

    if (/\s$/.test(newValue)) {
      const completeHashtags = extractCompleteHashtags(newValue)
      onHashtagDetected(completeHashtags)
    }
  }

  // Convertir hashtags a spans (solo visual)
  const convertHashtags = (text) => {
    if (!text) return '*Sin contenido*'

    let escaped = escapeHtml(text)

    escaped = escaped.replace(/#(\w+)(?=\s|$)/g, (match, tag) => {
      return `<span class="hashtag-inline" data-tag="${tag}">#${tag}</span>`
    })

    escaped = escaped.replace(/\n/g, '<br>')

    return escaped
  }

  return (
    <div style={{
      border: '1px solid var(--border)',
      borderRadius: '12px',
      overflow: 'hidden',
      backgroundColor: 'var(--bg-tertiary)'
    }}>
      {/* Barra de herramientas */}
      <div style={{
        display: 'flex',
        gap: '8px',
        padding: '8px 12px',
        borderBottom: '1px solid var(--border)',
        backgroundColor: 'var(--bg-secondary)'
      }}>
        <button
          type="button"
          onClick={() => setPreviewMode(false)}
          style={{
            padding: '4px 12px',
            borderRadius: '6px',
            backgroundColor: !previewMode ? 'var(--accent)' : 'transparent',
            color: !previewMode ? 'gray' : 'var(--text-secondary)',
            border: 'none',
            cursor: 'pointer'
          }}
        >
          Escribir
        </button>
        <button
          type="button"
          onClick={() => setPreviewMode(true)}
          style={{
            padding: '4px 12px',
            borderRadius: '6px',
            backgroundColor: previewMode ? 'var(--accent)' : 'transparent',
            color: previewMode ? 'gray' : 'var(--text-secondary)',
            border: 'none',
            cursor: 'pointer'
          }}
        >
          Vista previa
        </button>
      </div>

      {/* Editor */}
      {!previewMode ? (
        <textarea
          value={value}
          onChange={handleChange}
          placeholder={placeholder}
          rows={rows}
          style={{
            width: '100%',
            padding: '12px',
            backgroundColor: 'var(--bg-tertiary)',
            color: 'var(--text-primary)',
            border: 'none',
            fontFamily: 'monospace',
            fontSize: '14px',
            resize: 'vertical'
          }}
        />
      ) : (
        <div
          className="markdown-preview"
          style={{
            padding: '12px',
            minHeight: `${rows * 20}px`,
            backgroundColor: 'var(--bg-tertiary)',
            color: 'var(--text-primary)',
            overflow: 'auto',
            lineHeight: '1.6'
          }}
        >
          <div
            dangerouslySetInnerHTML={{
              __html: convertHashtags(value || '*Sin contenido*')
            }}
          />
        </div>
      )}
    </div>
  )
}

export default MarkdownEditor
