import { useState } from 'react'
import { FiHelpCircle } from 'react-icons/fi'
import AnswerCard from '../components/ui/AnswerCard'
import PrimaryButton from '../components/ui/PrimaryButton'
import { api } from '../services/api'

// Página para interactuar con el asistente inteligente.
function AskPage() {
  const [question, setQuestion] = useState('')
  const [answer, setAnswer] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!question.trim()) return

    setLoading(true)
    setAnswer('')

    try {
      const data = await api.ask(question)
      setAnswer(data.answer)
    } catch (error) {
      console.error('Error:', error)
      setAnswer('Error al conectar con el servidor')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <div className="page-header" style={{ textAlign: 'center', marginBottom: '2rem' }}>
        <div style={{
          display: 'inline-flex',
          alignItems: 'center',
          justifyContent: 'center',
          width: '64px',
          height: '64px',
          borderRadius: '32px',
          backgroundColor: 'rgba(217, 28, 192, 0.15)',
          marginBottom: '1rem'
        }}>
          <FiHelpCircle size={32} strokeWidth={1.5} style={{ color: '#d91cc0' }} />
        </div>
        <h1 className="page-title" style={{ fontSize: '2rem', fontWeight: 600, marginBottom: '0.5rem' }}>
          Asistente Inteligente
        </h1>
        <p className="page-subtitle" style={{ color: 'var(--text-muted)', fontSize: '0.9rem' }}>
          Estoy aquí para ayudarte, pregúntame lo que necesites en relación con tus notas.
        </p>
      </div>

      <div className="card" style={{
        backgroundColor: 'var(--bg-secondary)',
        border: '1px solid var(--border)',
        borderRadius: '20px',
        padding: '1.5rem'
      }}>

        <form onSubmit={handleSubmit}>
          <input
            type="text"
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleSubmit(e)}
            placeholder="Escribe aquí tu pregunta..."
            style={{
              width: '100%',
              marginBottom: '1.5rem',
              fontSize: '1rem',
              backgroundColor: 'var(--bg-tertiary)',
              color: 'var(--text-primary)',
              border: '1px solid var(--border)',
              borderRadius: '12px',
              padding: '12px 16px',
              fontFamily: 'inherit'
            }}
          />

          <div style={{ display: 'flex', justifyContent: 'center' }}>
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
        </form>
      </div>

      {/* Respuesta */}
      <AnswerCard answer={answer} />
    </div>
  )
}

export default AskPage