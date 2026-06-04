import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { api } from '../services/api'

// Componente para renderizar Markdown con enlaces internos a notas.
function MarkdownWithLinks({ content, onNavigateToNote }) {
  const processLinks = (text) => {
    if (!text) return ''

    let processed = text.replace(/\[\[([^\[\]]+)\]\]/g, (match, title) => {
      // Resolver el título a un ID de nota real
      return `[${title}](/note/search?title=${encodeURIComponent(title)})`
    })

    processed = processed.replace( /\[([^\]]+)\]\((?!https?:\/\/|\/)([^)]+)\)/g, (match, label, id) => {
      return `[${label}](/note/${id})`
    })

    return (
    <div>
      {processedContent}
    </div>)
  }

  const handleLinkClick = async (href, event) => {
    event.preventDefault()
    event.stopPropagation()

    if (href.includes('?title=')) {
      const titleToFind = decodeURIComponent(href.split('?title=')[1]);

      try {
        const targetNote = await api.getNoteByTitle(titleToFind);

        if (targetNote && targetNote.id) {
          onNavigateToNote(targetNote.id);
        }
      } catch (error) {
        console.error("Error resolviendo título:", titleToFind, error);
        alert("No se encontró una nota con el título: " + titleToFind);
      }
    } else if (href.startsWith('/note/')) {
      const noteId = href.replace('/note/', '')
      onNavigateToNote(noteId)
    }
  }

  return (
    <ReactMarkdown
      remarkPlugins={[remarkGfm]}
      components={{
        a: ({ href, children, ...props }) => (
          <a
            href={href}
            onClick={(e) => handleLinkClick(href, e)}
            style={{
              color: '#a78bfa',
              fontWeight: '500',
              textDecoration: 'none',
              borderBottom: '1px solid rgba(167, 139, 250, 0.3)',
              transition: 'all 0.2s ease'
            }}
            onMouseOver={(e) => e.target.style.borderBottom = '1px solid #a78bfa'}
            onMouseOut={(e) => e.target.style.borderBottom = '1px solid rgba(167, 139, 250, 0.3)'}
            {...props}
          >
            {children}
          </a>
        )
      }}
    >
      {processLinks(content)}
    </ReactMarkdown>
  )
}

export default MarkdownWithLinks 
