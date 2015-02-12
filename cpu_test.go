// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

import (
	"testing"

	"github.com/bmizerany/assert"
)

// dumbSound - Only for testing purpuses
type dumbSound struct{}

func (s dumbSound) Beep() {
	// play beep
}

// newCpuAt creates a new cpu and write the ins to the pc memory location
// Only for testing purpuse
func newCpuAt(ins uint16) cpu {
	s := dumbSound{}
	c := newCpu(s)
	c.memory.WriteWord(uint16(c.pc), ins)

	return c
}

// 00E0 - CLS
func Test00E0(t *testing.T) {

	c := newCpuAt(0x00E0)
	c.gfx.SetPixel(0, 0)
	c.Step()

	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			assert.Equal(t, byte(COLOR_BLACK), byte(c.gfx.GetPixel(i, j)))
		}
	}
}

// 00EE - RET
func Test00EE(t *testing.T) {
	c := newCpuAt(0x00EE)

	c.stack.Push(pc(1))
	c.Step()

	assert.Equal(t, pc(1), c.pc)
}

// 1nnn - JP addr
func Test1000(t *testing.T) {
	c := newCpuAt(0x1FFF)
	c.Step()

	assert.Equal(t, pc(0x0FFF), c.pc)
}

// 2nnn - CALL addr
func Test2000(t *testing.T) {
	c := newCpuAt(0x2FFF)
	oldPC := c.pc + 2 // The increment happen before the Step

	c.Step()

	assert.Equal(t, oldPC, c.stack.Pop())
	assert.Equal(t, pc(0x0FFF), c.pc)
}

// 3xkk - SE Vx, byte
func Test3000Equal(t *testing.T) {
	c := newCpuAt(0x31FF)
	c.regs[1] = 0xFF

	oldPC := c.pc + 4
	c.Step()

	assert.Equal(t, oldPC, c.pc)
}

// 3xkk - SE Vx, byte
func Test3000NotEqual(t *testing.T) {
	c := newCpuAt(0x31FF)
	c.regs[1] = 0xAA

	oldPC := c.pc + 2
	c.Step()

	assert.Equal(t, oldPC, c.pc)
}

// 4xkk - SNE Vx, byte
func Test4000Equal(t *testing.T) {
	c := newCpuAt(0x41FF)
	c.regs[1] = 0xFF

	oldPC := c.pc + 2
	c.Step()

	assert.Equal(t, oldPC, c.pc)
}

// 4xkk - SNE Vx, byte
func Test4000NotEqual(t *testing.T) {
	c := newCpuAt(0x41FF)
	c.regs[1] = 0xAA

	oldPC := c.pc + 4
	c.Step()

	assert.Equal(t, oldPC, c.pc)
}