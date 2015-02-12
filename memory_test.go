// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestRead(t *testing.T) {
	var mem memory
	mem[0] = 1
	mem[1] = 1

	assert.Equal(t, mem[0], mem.Read(0))
	assert.Equal(t, (uint16(mem[0])<<8 | uint16(mem[1])), mem.ReadWord(0))
}

func TestWrite(t *testing.T) {
	var mem memory
	mem[0] = 1

	assert.Equal(t, byte(1), mem[0])
	mem.Write(0, 2)

	assert.Equal(t, byte(2), mem[0])

	mem.WriteWord(0, 0x00FF)
	assert.Equal(t, byte(0), mem[0])
	assert.Equal(t, byte(0xFF), mem[1])
}
