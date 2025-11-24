package parsers

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type latexMath struct{}

var LatexMath = &latexMath{}

func (e *latexMath) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(newLatexParser(), 2000),
	))
}

type latexParser struct{}

func newLatexParser() parser.InlineParser {
	return &latexParser{}
}

func (s *latexParser) Trigger() []byte {
	return []byte{'\\', '$'}
}

func (s *latexParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) == 0 {
		return nil
	}

	char := line[0]
	
	// Manejo de \( y \[
	if char == '\\' {
		if len(line) < 2 {
			return nil
		}
		nextChar := line[1]
		if nextChar != '(' && nextChar != '[' {
			return nil
		}

		isDisplay := nextChar == '['
		closeDelim := []byte{'\\', ')'}
		if isDisplay {
			closeDelim = []byte{'\\', ']'}
		}

		return s.parseDelimited(block, closeDelim, 2)
	}

	// Manejo de $ y $$
	if char == '$' {
		isDisplay := false
		delimLen := 1
		if len(line) > 1 && line[1] == '$' {
			isDisplay = true
			delimLen = 2
		}

		// Si es $ inline, verificar que no haya espacio después (regla común)
		// y que no esté escapado (Goldmark maneja escapes antes si prioridad es baja, pero aquí somos prioridad alta)
		// Espera, si somos prioridad 2000, ganamos al escape.
		// Pero '\$' debería ser texto literal $.
		// Si Trigger incluye '\', ya lo manejamos arriba. Si es '\$', nextChar es '$', no '(' ni '['. Retornamos nil.
		// Entonces el escape parser de Goldmark (prioridad baja) lo tomará?
		// Sí, si retornamos nil, Goldmark sigue probando.
		
		closeDelim := []byte{'$'}
		if isDisplay {
			closeDelim = []byte{'$', '$'}
		}
		
		return s.parseDelimited(block, closeDelim, delimLen)
	}

	return nil
}

func (s *latexParser) parseDelimited(block text.Reader, closeDelim []byte, openLen int) ast.Node {
	node := ast.NewRawHTML()
	
	// Capturar primer segmento
	line, segment := block.PeekLine()
	
	// Buscar cierre
	found := false
	
	// Avanzar pasado la apertura en el primer segmento
	// Pero cuidado, si el cierre está en la misma línea, index buscará desde el principio
	
	// Estrategia: Buscar cierre desde openLen
	idx := bytes.Index(line[openLen:], closeDelim)
	if idx >= 0 {
		// Encontrado en la misma línea
		// idx es relativo a line[openLen:], así que posición real es openLen + idx
		totalLen := openLen + idx + len(closeDelim)
		
		seg := text.NewSegment(segment.Start, segment.Start+totalLen)
		node.Segments.Append(seg)
		block.Advance(totalLen)
		return node
	}
	
	// No encontrado en primera línea, añadirla completa
	node.Segments.Append(segment)
	block.Advance(len(line))
	
	// Buscar en líneas siguientes
	for {
		line, segment = block.PeekLine()
		if line == nil {
			break
		}

		idx = bytes.Index(line, closeDelim)
		if idx >= 0 {
			totalLen := idx + len(closeDelim)
			seg := text.NewSegment(segment.Start, segment.Start+totalLen)
			node.Segments.Append(seg)
			block.Advance(totalLen)
			found = true
			break
		}

		node.Segments.Append(segment)
		block.Advance(len(line))
	}

	if !found {
		return nil
	}

	return node
}
