import { useState } from 'react'

// Componente de botón reutilizable.
function PrimaryButton({ children, onClick, disabled, style = {}, type = 'button' }) {
  const [hover, setHover] = useState(false)

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: '8px',
        backgroundColor: '#d91cc0',
        color: 'white',
        border: 'none',
        padding: '10px 20px',
        borderRadius: '8px',
        fontWeight: '600',
        fontSize: '0.9rem',
        cursor: disabled ? 'not-allowed' : 'pointer',
        opacity: disabled ? 0.6 : 1,
        filter: hover && !disabled ? 'brightness(1.1)' : 'none',
        transition: 'all 0.2s',
        ...style
      }}
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
    >
      {children}
    </button>
  )
}

export default PrimaryButton