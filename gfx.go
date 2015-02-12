// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

const (
	COLOR_BLACK = iota
	COLOR_WHITE
)

// gfx is the chip8 screen buffer
type gfx [64][32]int

// Return the given pixel
func (g *gfx) GetPixel(x, y int) int {
	return g[x][y]
}

// SetPixel set set the pixel of the given position
func (g *gfx) SetPixel(x, y int) {
	g[x][y] = COLOR_WHITE
}

// ClearPixel set set the pixel of the given position
func (g *gfx) ClearPixel(x, y int) {
	g[x][y] = COLOR_BLACK
}

// Cls Clear the screen
func (g *gfx) Cls() {
	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			g.ClearPixel(i, j)
		}
	}
}
