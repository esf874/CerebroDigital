// Selector para mostrar el estado de una nota.
function NoteStatusSelector({ value, onChange }) {
  const options = [
    { value: "pending", label: "Pendiente", emoji: "○", color: "#a2d9c7" },
    { value: "in_progress", label: "Progreso", emoji: "◐", color: "#60f0c0" },
    { value: "finished", label: "Completada", emoji: "●", color: "#0dae78" }
  ]

  return (
    <div style={{
      display: 'flex',
      gap: '4px',
      backgroundColor: 'rgba(255,255,255,0.05)',
      padding: '4px',
      borderRadius: '12px'
    }}>
      {options.map(opt => (
        <button
          key={opt.value}
          type="button"
          onClick={() => onChange(opt.value)}
          style={{
            flex: 1,
            padding: '8px 16px',
            borderRadius: '8px',
            border: 'none',
            backgroundColor: value === opt.value ? `${opt.color}20` : 'transparent',
            color: value === opt.value ? opt.color : 'var(--text-secondary)',
            cursor: 'pointer',
            fontSize: '0.85rem',
            fontWeight: 500,
            transition: 'all 0.2s ease',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            gap: '6px'
          }}
        >
          <span style={{ fontSize: '1rem' }}>{opt.emoji}</span>
          <span>{opt.label}</span>
        </button>
      ))}
    </div>
  )
}

// Selector para mostrar la prioridad de una nota.
function NotePrioritySelector({ value, onChange }) {
  const options = [
    { value: "low", label: "Baja", color: "#f2b2d5", intensity: 1 },
    { value: "medium", label: "Media", color: "#df82d2", intensity: 2 },
    { value: "high", label: "Alta", color: "#de3ac8", intensity: 3 }
  ]

  return (
    <div style={{
      display: 'flex',
      gap: '4px',
      backgroundColor: 'rgba(255,255,255,0.05)',
      padding: '4px',
      borderRadius: '12px'
    }}>
      {options.map(opt => (
        <button
          key={opt.value}
          type="button"
          onClick={() => onChange(opt.value)}
          style={{
            flex: 1,
            padding: '8px 16px',
            borderRadius: '8px',
            border: 'none',
            backgroundColor: value === opt.value ? `${opt.color}20` : 'transparent',
            color: value === opt.value ? opt.color : 'var(--text-secondary)',
            cursor: 'pointer',
            fontSize: '0.85rem',
            fontWeight: 500,
            transition: 'all 0.2s ease',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            gap: '6px'
          }}
        >
          <span style={{
            display: 'flex',
            gap: '2px',
            opacity: opt.intensity / 3
          }}>
            {'•'.repeat(opt.intensity)}
          </span>
          <span>{opt.label}</span>
        </button>
      ))}
    </div>
  )
}

function NoteStatusBadge({ status }) {
  const config = {
    pending: { emoji: "○", label: "Pendiente", color: "#a2d9c7" },
    in_progress: { emoji: "◐", label: "Progreso", color: "#60f0c0" },
    finished: { emoji: "●", label: "Completada", color: "#0dae78" }
  }

  const c = config[status] || config.pending

  return (
    <span style={{
      display: 'inline-flex',
      alignItems: 'center',
      gap: '6px',
      padding: '4px 12px',
      backgroundColor: `${c.color}15`,
      color: c.color,
      borderRadius: '20px',
      fontSize: '0.75rem',
      fontWeight: 500
    }}>
      <span>{c.emoji}</span>
      <span>{c.label}</span>
    </span>
  )
}

function NotePriorityBadge({ priority }) {
  const config = {
    low: { dots: "•", label: "Baja", color: "#f2b2d5", intensity: 1 },
    medium: { dots: "••", label: "Media", color: "#df82d2", intensity: 2 },
    high: { dots: "•••", label: "Alta", color: "#de3ac8", intensity: 3 }
  }

  const c = config[priority] || config.medium

  return (
    <span style={{
      display: 'inline-flex',
      alignItems: 'center',
      gap: '6px',
      padding: '4px 12px',
      backgroundColor: `${c.color}15`,
      color: c.color,
      borderRadius: '20px',
      fontSize: '0.75rem',
      fontWeight: 500
    }}>
      <span style={{ opacity: c.intensity / 3 }}>{c.dots}</span>
      <span>{c.label}</span>
    </span>
  )
}

export { NotePriorityBadge, NotePrioritySelector, NoteStatusBadge, NoteStatusSelector }
