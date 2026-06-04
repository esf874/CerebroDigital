import PropTypes from 'prop-types'
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

  const renderFormattedText = (text) => {
    if (!text) return '*Sin contenido*'

    const escaped = text
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')

    const parts = escaped.split(/(\n|#\w+(?=\s|$))/g)

    return parts.map((part, i) => {
      if (part === '\n') return <br key={i} />

      if (part.startsWith('#')) {
        const tag = part.slice(1)
        return (
          <span key={i} className="hashtag-inline" data-tag={tag}>
            #{tag}
          </span>
        )
      }

      return part
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
          <div>
            {renderFormattedText(value)}
          </div>
        </div>
      )}
    </div>
  )
}

MarkdownEditor.propTypes = {
  value: PropTypes.string,
  onChange: PropTypes.func.isRequired,
  placeholder: PropTypes.string,
  rows: PropTypes.number,
  onHashtagDetected: PropTypes.func.isRequired
}

export default MarkdownEditor
