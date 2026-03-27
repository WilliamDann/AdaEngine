package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// KittyImage is a tview primitive that displays an image via the Kitty graphics protocol.
type KittyImage struct {
	*tview.Box
	img     image.Image
	dirty   bool
	lastCol int
	lastRow int
}

func NewKittyImage() *KittyImage {
	return &KittyImage{
		Box: tview.NewBox(),
	}
}

func (k *KittyImage) SetImage(img image.Image) {
	k.img = img
	k.dirty = true
}

func (k *KittyImage) Draw(screen tcell.Screen) {
	k.Box.DrawForSubclass(screen, k)

	if k.img == nil {
		return
	}

	x, y, w, h := k.GetInnerRect()

	// Only redraw if position changed or image is dirty
	if !k.dirty && x == k.lastCol && y == k.lastRow {
		return
	}
	k.dirty = false
	k.lastCol = x
	k.lastRow = y

	// Encode image as PNG
	var buf bytes.Buffer
	png.Encode(&buf, k.img)
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Calculate display size in pixels based on cell size
	// Use columns/rows to set display size so the image fits the region
	imgW := k.img.Bounds().Dx()
	imgH := k.img.Bounds().Dy()
	_ = w
	_ = h

	// Write Kitty graphics escape directly to the terminal.
	// Use virtual placement: transmit then display at cursor position.
	// First, delete any previous image at this ID.
	fmt.Fprintf(os.Stdout, "\033_Ga=d,d=I,i=1\033\\")

	// Move cursor to the board position
	fmt.Fprintf(os.Stdout, "\033[%d;%dH", y+1, x+1)

	// Transmit and display in chunks
	const chunkSize = 4096
	for i := 0; i < len(b64); i += chunkSize {
		end := i + chunkSize
		more := 1
		if end >= len(b64) {
			end = len(b64)
			more = 0
		}
		chunk := b64[i:end]
		if i == 0 {
			// f=100 (PNG), a=T (transmit+display), i=1 (image ID)
			// c/r = columns/rows to display in
			fmt.Fprintf(os.Stdout, "\033_Gf=100,a=T,i=1,s=%d,v=%d,m=%d;%s\033\\",
				imgW, imgH, more, chunk)
		} else {
			fmt.Fprintf(os.Stdout, "\033_Gm=%d;%s\033\\", more, chunk)
		}
	}
}
