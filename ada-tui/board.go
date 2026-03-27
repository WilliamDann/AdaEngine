package main

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

//go:embed pieces/*.svg
var pieceFS embed.FS

const sqSize = 90

var (
	lightSquare = color.RGBA{240, 217, 181, 255}
	darkSquare  = color.RGBA{181, 136, 99, 255}
)

// SVG file names indexed by [color][pieceType].
// color 0 = white (l), color 1 = black (d).
var svgName = [2][7]string{
	{"", "Chess_plt45.svg", "Chess_nlt45.svg", "Chess_blt45.svg", "Chess_rlt45.svg", "Chess_qlt45.svg", "Chess_klt45.svg"},
	{"", "Chess_pdt45.svg", "Chess_ndt45.svg", "Chess_bdt45.svg", "Chess_rdt45.svg", "Chess_qdt45.svg", "Chess_kdt45.svg"},
}

// pieceImages caches rasterized piece images.
var pieceImages [2][7]*image.RGBA

func loadPieces() {
	for c := 0; c < 2; c++ {
		for p := 1; p <= 6; p++ {
			name := svgName[c][p]
			data, err := pieceFS.ReadFile("pieces/" + name)
			if err != nil {
				panic("missing piece SVG: " + name + ": " + err.Error())
			}
			pieceImages[c][p] = rasterSVG(data, sqSize, sqSize)
		}
	}
}

func rasterSVG(data []byte, w, h int) *image.RGBA {
	icon, err := oksvg.ReadIconStream(newByteReadCloser(data))
	if err != nil {
		panic("bad SVG: " + err.Error())
	}
	icon.SetTarget(0, 0, float64(w), float64(h))
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scanner)
	icon.Draw(raster, 1.0)
	return img
}

// renderBoard creates an image of the full board for the given position.
func renderBoard(pos *position.Position) *image.RGBA {
	size := sqSize * 8
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			x := file * sqSize
			y := (7 - rank) * sqSize
			rect := image.Rect(x, y, x+sqSize, y+sqSize)

			// Draw square
			sqColor := lightSquare
			if (rank+file)%2 == 0 {
				sqColor = darkSquare
			}
			draw.Draw(img, rect, &image.Uniform{sqColor}, image.Point{}, draw.Src)

			// Draw piece
			sq := core.NewSquare(rank, file)
			piece := pos.Board.Check(sq)
			if piece != core.None {
				colorIdx := 0
				if piece.Color() == core.Black {
					colorIdx = 1
				}
				pieceImg := pieceImages[colorIdx][piece.Type()]
				if pieceImg != nil {
					draw.Draw(img, rect, pieceImg, image.Point{}, draw.Over)
				}
			}
		}
	}

	return img
}

// displayBoard renders the board and prints it to stdout via Kitty graphics protocol.
func displayBoard(pos *position.Position) {
	img := renderBoard(pos)
	var buf bytes.Buffer
	png.Encode(&buf, img)
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Send in 4096-byte chunks per the Kitty protocol
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
			fmt.Printf("\033_Gf=100,a=T,m=%d;%s\033\\", more, chunk)
		} else {
			fmt.Printf("\033_Gm=%d;%s\033\\", more, chunk)
		}
	}
	fmt.Println()
}

// byteReadCloser wraps a byte slice as an io.ReadCloser.
type byteReadCloser struct {
	data []byte
	pos  int
}

func newByteReadCloser(data []byte) *byteReadCloser {
	return &byteReadCloser{data: data}
}

func (b *byteReadCloser) Read(p []byte) (int, error) {
	if b.pos >= len(b.data) {
		return 0, io.EOF
	}
	n := copy(p, b.data[b.pos:])
	b.pos += n
	return n, nil
}

func (b *byteReadCloser) Close() error {
	return nil
}
