//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"syscall/js"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func renderMarkdownWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "Error: Se esperaba 1 argumento"
	}
	
	input := args[0].String()

	// Usamos la misma configuración que tu sitio web (goldmark)
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(), // Permitir HTML crudo si tu sitio lo permite
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		return "Error renderizando Markdown: " + err.Error()
	}

	return js.ValueOf(buf.String())
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("renderMarkdown", js.FuncOf(renderMarkdownWrapper))
	<-c
}