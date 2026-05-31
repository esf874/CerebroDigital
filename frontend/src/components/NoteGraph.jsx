import * as d3 from 'd3'
import { useCallback, useEffect, useRef, useState } from 'react'

const COLORS = {
  current: '#d91cc0',
  near: '#60f0c0',
  far: '#f2b2d5',
  link: '#4a4a5a',
  linkHover: '#60f0c0',
  text: '#e0e0e0'
}

// Componente para visualizar el subgrafo de conexiones entre notas.
function NoteGraph({ noteId, onNodeClick, depth = 2, limit = 20, graphData: externalData, loading: externalLoading }) {
  const [internalData, setInternalData] = useState({ nodes: [], edges: [] })
  const [internalLoading, setInternalLoading] = useState(false)
  const [hoveredNode, setHoveredNode] = useState(null)

  const graphData = externalData || internalData
  const loading = externalLoading ?? internalLoading

  const svgRef = useRef()
  const simulationRef = useRef(null)

  useEffect(() => {
    if (externalData) return

    const fetchGraph = async () => {
      try {
        setInternalLoading(true)
        const { api } = await import('../services/api')
        const data = await api.getGraph(noteId, depth, limit)
        setInternalData(data)
      } catch (err) {
        console.error('Error fetching graph:', err)
      } finally {
        setInternalLoading(false)
      }
    }

    if (noteId) fetchGraph()
  }, [noteId, depth, limit, externalData])

  const dragstarted = useCallback((event, d) => {
    if (!event.active && simulationRef.current) {
      simulationRef.current.alphaTarget(0.3).restart()
    }
    d.fx = d.x
    d.fy = d.y
  }, [])

  const dragged = useCallback((event, d) => {
    d.fx = event.x
    d.fy = event.y
  }, [])

  const dragended = useCallback((event, d) => {
    if (!event.active && simulationRef.current) {
      simulationRef.current.alphaTarget(0)
    }
    d.fx = null
    d.fy = null
  }, [])


  useEffect(() => {
    if (!graphData.nodes?.length || !svgRef.current) return

    const width = 600
    const height = 400

    const svg = d3.select(svgRef.current)
    svg.selectAll('*').remove()

    // Contenedor principal con zoom
    const g = svg
      .attr('viewBox', [-width / 2, -height / 2, width, height])
      .append('g')

    // Zoom behavior
    const zoom = d3.zoom()
      .scaleExtent([0.3, 3])
      .on('zoom', (event) => {
        g.attr('transform', event.transform)
      })

    svg.call(zoom)

    // Definir gradiente para enlaces
    const defs = svg.append('defs')

    const filter = defs.append('filter')
      .attr('id', 'glow')
      .attr('x', '-50%')
      .attr('y', '-50%')
      .attr('width', '200%')
      .attr('height', '200%')

    filter.append('feGaussianBlur')
      .attr('stdDeviation', '3')
      .attr('result', 'coloredBlur')

    const feMerge = filter.append('feMerge')
    feMerge.append('feMergeNode').attr('in', 'coloredBlur')
    feMerge.append('feMergeNode').attr('in', 'SourceGraphic')

    // Enlaces curvos 
    const link = g
      .append('g')
      .attr('class', 'links')
      .selectAll('path')
      .data(graphData.edges)
      .join('path')
      .attr('fill', 'none')
      .attr('stroke', COLORS.link)
      .attr('stroke-width', 1.5)
      .attr('stroke-opacity', 0.6)
      .style('transition', 'stroke 0.2s, stroke-opacity 0.2s')

    // Nodos
    const node = g
      .append('g')
      .attr('class', 'nodes')
      .selectAll('g')
      .data(graphData.nodes)
      .join('g')
      .attr('class', 'node-group')
      .style('cursor', 'pointer')
      .call(
        d3.drag()
          .on('start', dragstarted)
          .on('drag', dragged)
          .on('end', dragended)
      )

    node.append('circle')
      .attr('class', 'node-ring')
      .attr('r', d => (d.isCurrent ? 16 : 12))
      .attr('fill', 'transparent')
      .attr('stroke', d => getNodeColor(d))
      .attr('stroke-width', 2)
      .attr('stroke-opacity', 0.3)

    // Círculo principal
    node.append('circle')
      .attr('class', 'node-circle')
      .attr('r', d => (d.isCurrent ? 10 : 7))
      .attr('fill', d => getNodeColor(d))
      .attr('filter', d => d.isCurrent ? 'url(#glow)' : null)
      .style('transition', 'r 0.2s, fill 0.2s')

    const label = g
      .append('g')
      .attr('class', 'labels')
      .selectAll('text')
      .data(graphData.nodes)
      .join('text')
      .text(d => truncateText(d.title, 20))
      .attr('font-size', d => d.isCurrent ? 11 : 9)
      .attr('font-weight', d => d.isCurrent ? 600 : 400)
      .attr('fill', COLORS.text)
      .attr('opacity', 0.85)
      .attr('pointer-events', 'none')
      .style('text-shadow', '0 1px 3px rgba(0,0,0,0.8)')

    const tooltip = d3.select('body')
      .append('div')
      .attr('class', 'graph-tooltip')
      .style('position', 'fixed')
      .style('visibility', 'hidden')
      .style('background', 'rgba(20, 20, 25, 0.95)')
      .style('border', '1px solid rgba(255,255,255,0.1)')
      .style('border-radius', '8px')
      .style('padding', '8px 12px')
      .style('font-size', '12px')
      .style('color', '#fff')
      .style('pointer-events', 'none')
      .style('z-index', '1000')
      .style('max-width', '200px')
      .style('backdrop-filter', 'blur(8px)')

    // Eventos de hover
    node
      .on('mouseenter', function (event, d) {
        // Highlight nodo actual
        d3.select(this).select('.node-circle')
          .transition()
          .duration(150)
          .attr('r', d.isCurrent ? 13 : 10)

        d3.select(this).select('.node-ring')
          .transition()
          .duration(150)
          .attr('stroke-opacity', 0.6)

        link
          .attr('stroke', l =>
            (l.source.id === d.id || l.target.id === d.id)
              ? COLORS.linkHover
              : COLORS.link
          )
          .attr('stroke-opacity', l =>
            (l.source.id === d.id || l.target.id === d.id) ? 1 : 0.3
          )
          .attr('stroke-width', l =>
            (l.source.id === d.id || l.target.id === d.id) ? 2 : 1.5
          )

        // Mostrar tooltip de cada nodo
        tooltip
          .style('visibility', 'visible')
          .html(`
            <div style="font-weight: 600; margin-bottom: 4px;">${d.title}</div>
            <div style="font-size: 10px; color: #888;">
              ${d.isCurrent ? 'Nota actual' : `Profundidad: ${d.depth || '?'}`}
            </div>
          `)
          .style('left', (event.clientX + 15) + 'px')
          .style('top', (event.clientY - 10) + 'px')
      })
      .on('mousemove', function (event) {
        tooltip
          .style('left', (event.clientX + 15) + 'px')
          .style('top', (event.clientY - 10) + 'px')
      })
      .on('mouseleave', function (event, d) {
        d3.select(this).select('.node-circle')
          .transition()
          .duration(150)
          .attr('r', d.isCurrent ? 10 : 7)

        d3.select(this).select('.node-ring')
          .transition()
          .duration(150)
          .attr('stroke-opacity', 0.3)

        link
          .attr('stroke', COLORS.link)
          .attr('stroke-opacity', 0.6)
          .attr('stroke-width', 1.5)

        tooltip.style('visibility', 'hidden')
      })
      .on('click', (event, d) => {
        event.stopPropagation()
        onNodeClick?.(d.id)
      })

    const simulation = d3
      .forceSimulation(graphData.nodes)
      .force('link', d3.forceLink(graphData.edges)
        .id(d => d.id)
        .distance(100)
      )
      .force('charge', d3.forceManyBody().strength(-300))
      .force('center', d3.forceCenter(0, 0))
      .force('collision', d3.forceCollide().radius(40))
      .on('tick', () => {
        link.attr('d', d => {
          const dx = d.target.x - d.source.x
          const dy = d.target.y - d.source.y
          const dr = Math.sqrt(dx * dx + dy * dy) * 1.5
          return `M${d.source.x},${d.source.y}A${dr},${dr} 0 0,1 ${d.target.x},${d.target.y}`
        })

        node.attr('transform', d => `translate(${d.x},${d.y})`)

        label
          .attr('x', d => d.x + 14)
          .attr('y', d => d.y + 4)
      })

    simulationRef.current = simulation

    return () => {
      simulation.stop()
      tooltip.remove()
    }
  }, [graphData, onNodeClick, dragstarted, dragged, dragended])

  // Helpers
  function getNodeColor(d) {
    if (d.isCurrent) return COLORS.current
    if (d.depth === 1) return COLORS.near
    return COLORS.far
  }

  function truncateText(text, maxLength) {
    if (!text) return ''
    return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
  }

  // Estados de carga y vacío
  if (loading) {
    return (
      <div style={{
        height: '350px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: 'var(--bg-secondary)',
        borderRadius: '16px',
        color: 'var(--text-muted)'
      }}>
        <span style={{ marginRight: '8px' }}>🔄</span>
        <span>Cargando conexiones...</span>
      </div>
    )
  }

  if (!graphData.nodes?.length) {
    return (
      <div style={{
        height: '350px',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: 'var(--bg-secondary)',
        borderRadius: '16px',
        color: 'var(--text-muted)'
      }}>
        <span style={{ fontSize: '2rem', marginBottom: '0.5rem' }}>📭</span>
        <p style={{ fontSize: '0.85rem' }}>Esta nota aún no tiene conexiones</p>
        <p style={{ fontSize: '0.7rem', opacity: 0.6 }}>Las conexiones se crean al compartir etiquetas</p>
      </div>
    )
  }

  return (
    <div style={{ position: 'relative' }}>
      {/* Grafo SVG */}
      <div
        style={{
          borderRadius: '16px',
          overflow: 'hidden',
          background: 'linear-gradient(135deg, rgba(20,20,25,0.8) 0%, rgba(30,30,40,0.6) 100%)',
          border: '1px solid rgba(255,255,255,0.06)'
        }}
      >
        <svg
          ref={svgRef}
          style={{ width: '100%', height: '380px', display: 'block' }}
        />
      </div>

      {/* Leyenda */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: '1.5rem',
        marginTop: '1rem',
        fontSize: '0.7rem',
        color: 'var(--text-muted)'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
          <span style={{ width: '10px', height: '10px', borderRadius: '50%', backgroundColor: COLORS.current, boxShadow: `0 0 6px ${COLORS.current}50` }}></span>
          <span>Nota actual</span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
          <span style={{ width: '10px', height: '10px', borderRadius: '50%', backgroundColor: COLORS.near, boxShadow: `0 0 6px ${COLORS.near}50` }}></span>
          <span>Cercanas</span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
          <span style={{ width: '10px', height: '10px', borderRadius: '50%', backgroundColor: COLORS.far, boxShadow: `0 0 6px ${COLORS.far}50` }}></span>
          <span>Lejanas</span>
        </div>
      </div>

      <p style={{
        textAlign: 'center',
        fontSize: '0.65rem',
        color: 'var(--text-muted)',
        opacity: 0.5,
        marginTop: '0.5rem'
      }}>
        Scroll para zoom · Arrastra para mover
      </p>
    </div>
  )
}

export default NoteGraph
