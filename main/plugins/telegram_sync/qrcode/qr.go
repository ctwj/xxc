// Package qr provides a minimal stub for QR code generation.
// This is a stub implementation to satisfy gotd/td's qrlogin dependency.
package qr

import (
	"errors"
	"image"
)

// Level represents QR error correction level.
type Level int

const (
	L Level = iota
	M
	Q
	H
)

// Encode creates a QR code from the given text.
// This is a stub implementation that returns an error.
func Encode(text string, level Level) (*Code, error) {
	return nil, errors.New("QR code generation not supported - use phone authentication instead")
}

// Code represents a QR code.
type Code struct {
	Text  string
	Level Level
}

// Image returns an image of the QR code.
// This is a stub implementation.
func (c *Code) Image() image.Image {
	return image.NewRGBA(image.Rect(0, 0, 1, 1))
}

// Bitmap returns the QR code bitmap.
func (c *Code) Bitmap() [][]bool {
	return nil
}

// String returns the text encoded in the QR code.
func (c *Code) String() string {
	return c.Text
}
