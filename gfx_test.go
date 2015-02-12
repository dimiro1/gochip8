// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestGetPixel(t *testing.T) {
	g := new(gfx)
	g.SetPixel(0, 0)

	assert.Equal(t, COLOR_WHITE, g.GetPixel(0, 0))
}

func TestClearPixel(t *testing.T) {
	g := new(gfx)
	g.SetPixel(0, 0)
	assert.Equal(t, COLOR_WHITE, g.GetPixel(0, 0))
	g.ClearPixel(0, 0)
	assert.Equal(t, COLOR_BLACK, g.GetPixel(0, 0))
}

func TestSetPixel(t *testing.T) {
	g := new(gfx)
	g.SetPixel(0, 0)

	assert.Equal(t, COLOR_WHITE, g[0][0])
}

func TestCls(t *testing.T) {
	g := new(gfx)
	g.SetPixel(0, 0)

	g.Cls()

	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			assert.Equal(t, COLOR_BLACK, g.GetPixel(i, j))
		}
	}
}
