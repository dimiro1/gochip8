// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestCls(t *testing.T) {
	g := new(gfx)
	g[0] = 1

	g.Cls()

	for _, e := range g {
		assert.Equal(t, byte(0), e)
	}
}
