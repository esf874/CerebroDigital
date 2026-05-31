// Componente para mostrar la respuesta del modelo de lenguaje.
function AnswerCard({ answer }) {
  if (!answer) return null

  return (
    <div style={{
      marginTop: '2rem',
      backgroundColor: 'rgba(217, 28, 192, 0.05)',
      border: '1px solid rgba(217, 28, 192, 0.2)',
      borderRadius: '20px',
      padding: '1.5rem'
    }}>
      <div style={{
        display: 'flex',
        alignItems: 'center',
        gap: '8px',
        marginBottom: '1rem'
      }}>
        <div style={{
          width: '6px',
          height: '24px',
          backgroundColor: '#d91cc0',
          borderRadius: '3px'
        }} />
        <h3 style={{ margin: 0, fontSize: '1rem', fontWeight: 600 }}>
          Respuesta
        </h3>
      </div>

      <p style={{
        whiteSpace: 'pre-wrap',
        lineHeight: 1.7,
        margin: 0,
        fontSize: '0.95rem'
      }}>
        {answer}
      </p>
    </div>
  )
}

export default AnswerCard