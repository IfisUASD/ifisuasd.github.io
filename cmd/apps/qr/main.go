//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/base64"
	"fmt"
	"strings"
	"syscall/js"
)

// generateQRWrapper is the function exported to JavaScript. It receives the
// parameters as strings (or base64‑encoded data for the logo) and returns a
// base64‑encoded PNG image, or an empty string on error.
//
// Expected JS call signature:
//
//	const imgBase64 = generateQR(text, size, fgColor, bgColor, logoBase64, level);
//	// imgBase64 contains "data:image/png;base64,...."
func generateQRWrapper(this js.Value, args []js.Value) interface{} {
	// Argument validation – we expect exactly 6 arguments.
	if len(args) != 6 {
		js.Global().Get("console").Call("error", "generateQRWrapper: expected 6 arguments")
		return js.Null()
	}

	text := args[0].String()
	size := args[1].Int()
	fgHex := args[2].String()
	bgHex := args[3].String()
	logoBase64 := args[4].String()
	levelStr := args[5].String()

	var logoBytes []byte
	if logoBase64 != "" {
		// The logo is passed as a base64 data URL (or raw base64). Strip any prefix.
		if strings.HasPrefix(logoBase64, "data:image/") {
			if idx := strings.Index(logoBase64, ","); idx != -1 {
				logoBase64 = logoBase64[idx+1:]
			}
		}
		var err error
		logoBytes, err = base64.StdEncoding.DecodeString(logoBase64)
		if err != nil {
			js.Global().Get("console").Call("error", "Failed to decode logo base64:", err.Error())
			return js.Null()
		}
	}

	// Call the pure Go QR generation logic.
	pngBytes, err := GenerateQR(text, size, fgHex, bgHex, logoBytes, levelStr)
	if err != nil {
		js.Global().Get("console").Call("error", "GenerateQR error:", err.Error())
		return js.Null()
	}

	// Encode the PNG bytes as a base64 data URL for easy use in the browser.
	b64 := base64.StdEncoding.EncodeToString(pngBytes)
	dataURL := "data:image/png;base64," + b64
	return js.ValueOf(dataURL)
}

func registerCallbacks() {
	// Expose the function under the name "generateQR" in the global JS scope.
	js.Global().Set("generateQR", js.FuncOf(generateQRWrapper))
}

func main() {
	fmt.Println("WASM QR generator initializing...")
	registerCallbacks()
	// Prevent the Go program from exiting.
	select {}
}
