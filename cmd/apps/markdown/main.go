//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"syscall/js"

	"github.com/IfisUASD/ifisuasd.github.io/internal/parsers"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

// renderMarkdownWrapper es la función que JavaScript llamará.
// Recibe un string (Markdown) y devuelve un string (HTML).
func renderMarkdownWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "Error: Se esperaba 1 argumento (texto markdown)"
	}

	input := args[0].String()

	// Configuración de Goldmark
	// Incluimos 'frontmatter' para que si el usuario pega un archivo con cabecera YAML,
	// esta se procese como metadata y no ensucie la visualización del contenido.
	// Incluimos 'GFM' para soportar tablas, tareas, etc.
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			parsers.LatexMath,
			&frontmatter.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(), // Necesario si permites HTML incrustado en tus posts
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		return "<div class='alert alert-error'>Error renderizando Markdown: " + err.Error() + "</div>"
	}

	return js.ValueOf(buf.String())
}

func main() {
	// Canal para mantener el programa corriendo
	c := make(chan struct{}, 0)

	// Exponer la función 'renderMarkdown' al ámbito global de JavaScript
	js.Global().Set("renderMarkdown", js.FuncOf(renderMarkdownWrapper))

	// Mantener vivo el proceso WASM
	<-c
}