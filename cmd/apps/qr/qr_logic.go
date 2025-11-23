package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	qrcode "github.com/skip2/go-qrcode"
	xdraw "golang.org/x/image/draw"
)

// GenerateQR creates a PNG image of a QR code.
//
// Parameters:
//
//	text   – the data to encode.
//	size   – width and height of the resulting PNG in pixels (must be > 0).
//	fgHex  – foreground colour in hex notation (e.g. "#000000").
//	         Currently unused – the default black colour is used.
//	bgHex  – background colour in hex notation (e.g. "#ffffff").
//	         Currently unused – the default white colour is used.
//	logo   – optional PNG bytes for a centre logo (currently ignored).
//	levelStr – recovery level ("Low", "Medium", "High", "Highest").
//
// Returns a slice containing PNG data, or an error.
func GenerateQR(text string, size int, fgHex string, bgHex string, logo []byte, levelStr string) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid size %d", size)
	}

	// Map string level to qrcode.RecoveryLevel
	var level qrcode.RecoveryLevel
	switch levelStr {
	case "Low":
		level = qrcode.Low
	case "Medium":
		level = qrcode.Medium
	case "High":
		level = qrcode.High
	case "Highest":
		level = qrcode.Highest
	default:
		level = qrcode.High
	}

	// Debug: show the parameters we received.
	fmt.Printf("[DEBUG] GenerateQR called with text=%q size=%d fg=%s bg=%s logoBytes=%d level=%s\n",
		text, size, fgHex, bgHex, len(logo), levelStr)

	// Create the QR code with the specified error‑recovery level.
	qr, err := qrcode.New(text, level)
	if err != nil {
		return nil, fmt.Errorf("qrcode.New failed: %w", err)
	}

	// NOTE: The go‑qrcode library does not expose colour setters for in‑memory PNG
	// generation. For the purposes of the current tests we rely on the default
	// black‑on‑white colours.

	// Encode the QR as PNG bytes.
	// We decode it back to an image.Image to manipulate colors if needed,
	// or if we need to overlay a logo.
	// Since we need to support custom colors, we'll always decode it.
	pngBytes, err := qr.PNG(size)
	if err != nil {
		return nil, fmt.Errorf("qr.PNG failed: %w", err)
	}

	img, err := png.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode generated QR: %w", err)
	}

	// Parse colors
	fgColor, err := parseHexColor(fgHex)
	if err != nil {
		// Fallback to black
		fgColor = color.RGBA{0, 0, 0, 255}
	}
	bgColor, err := parseHexColor(bgHex)
	if err != nil {
		// Fallback to white
		bgColor = color.RGBA{255, 255, 255, 255}
	}

	// Convert to RGBA and apply colors
	b := img.Bounds()
	m := image.NewRGBA(b)
	
	// Iterate over pixels to apply colors
	// go-qrcode generates black (0,0,0) for modules and white (255,255,255) for background
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			r, g, bVal, _ := c.RGBA()
			// Check if it's black (module) or white (background)
			// RGBA() returns 0-65535. Black is close to 0, White is close to 65535.
			// Threshold: < 32768 is "dark", >= 32768 is "light"
			if r < 0x8000 && g < 0x8000 && bVal < 0x8000 {
				m.Set(x, y, fgColor)
			} else {
				m.Set(x, y, bgColor)
			}
		}
	}

	// If a logo is supplied, overlay it.
	if len(logo) > 0 {
		// Decode the logo.
		logoImg, _, err := image.Decode(bytes.NewReader(logo))
		if err != nil {
			// If logo decoding fails, return error (or could fallback to plain QR).
			return nil, fmt.Errorf("failed to decode logo: %w", err)
		}

		// Calculate target logo size (e.g., 20% of QR size).
		logoSize := size / 5
		if logoSize < 1 {
			logoSize = 1
		}

		// Resize logo using high-quality scaling.
		dstRect := image.Rect(0, 0, logoSize, logoSize)
		dst := image.NewRGBA(dstRect)
		xdraw.CatmullRom.Scale(dst, dstRect, logoImg, logoImg.Bounds(), xdraw.Over, nil)

		// Calculate offset to center the logo.
		offset := image.Pt((size-logoSize)/2, (size-logoSize)/2)

		// Draw the logo over the QR code.
		draw.Draw(m, dst.Bounds().Add(offset), dst, image.Point{}, draw.Over)
	}

	// Encode the final image back to PNG.
	var buf bytes.Buffer
	if err := png.Encode(&buf, m); err != nil {
		return nil, fmt.Errorf("failed to encode final PNG: %w", err)
	}
	return buf.Bytes(), nil
}

// parseHexColor parses a hex string (e.g. "#RRGGBB") into color.RGBA.
func parseHexColor(s string) (color.RGBA, error) {
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}
	if len(s) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex color length")
	}
	var r, g, b uint8
	n, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
	if err != nil || n != 3 {
		return color.RGBA{}, fmt.Errorf("invalid hex color format")
	}
	return color.RGBA{r, g, b, 255}, nil
}
