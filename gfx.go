// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

const (
	COLOR_BLACK = iota
	COLOR_WHITE
)

// gfx is the chip8 screen buffer
type gfx [64 * 32]byte

// Cls Clear the screen
func (g *gfx) Cls() {
	for i, _ := range g {
		g[i] = COLOR_BLACK
	}
}
