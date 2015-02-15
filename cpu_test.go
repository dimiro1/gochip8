// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

import (
	"testing"

	"math/rand"

	"github.com/bmizerany/assert"
)

// dumbSound - Only for testing purpuses
type dumbSound struct {
	BeepCalled bool
}

func (s *dumbSound) Beep() {
	(*s).BeepCalled = true
}

// newCpuAt creates a new cpu and write the ins to the pc memory location
// Only for testing purpuse
func newCpuAt(ins uint16) cpu {
	s := &dumbSound{}
	c := newCpu(s)
	c.memory.WriteWord(uint16(c.pc), ins)

	return c
}

func TestTimers(t *testing.T) {
	s := &dumbSound{}
	c := newCpu(s)

	// With Beep
	c.memory.WriteWord(uint16(c.pc), 0x0E0)

	c.dt = 1
	c.st = 1

	assert.Equal(t, s.BeepCalled, false)

	c.Step()

	assert.Equal(t, c.dt, byte(0))
	assert.Equal(t, c.st, byte(0))
	assert.Equal(t, s.BeepCalled, true)

	// No Beep
	c.memory.WriteWord(uint16(c.pc), 0x0E0)
	s.BeepCalled = false

	c.st = 2
	c.Step()

	assert.Equal(t, s.BeepCalled, false)
	assert.Equal(t, c.st, byte(1))
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

// 5xy0 - SE Vx, Vy
func Test5000Equal(t *testing.T) {
	c := newCpuAt(0x512F)
	c.regs[1] = 0xFF
	c.regs[2] = 0xFF

	oldPC := c.pc + 4
	c.Step()

	assert.Equal(t, oldPC, c.pc)
}

// 5xy0 - SE Vx, Vy
func Test5000NotEqual(t *testing.T) {
	c := newCpuAt(0x512F)
	c.regs[1] = 0xFF
	c.regs[2] = 0xAA

	oldPC := c.pc + 2
	c.Step()

	assert.Equal(t, oldPC, c.pc)
}

// 6xkk - LD Vx, byte
func Test6000(t *testing.T) {
	c := newCpuAt(0x61FF)
	c.Step()
	assert.Equal(t, c.regs[1], byte(0xFF))
}

// 7xkk - ADD Vx, byte
func Test7000(t *testing.T) {
	c := newCpuAt(0x7101)
	c.regs[1] = 1
	c.Step()
	assert.Equal(t, c.regs[1], byte(2))
}

// 8xy0 - LD Vx, Vy
func Test8xy0(t *testing.T) {
	c := newCpuAt(0x8120)
	c.regs[1] = 1
	c.regs[2] = 0
	assert.Equal(t, c.regs[1], byte(1))
	assert.Equal(t, c.regs[2], byte(0))
	c.Step()
	assert.Equal(t, c.regs[1], byte(0))
	assert.Equal(t, c.regs[2], byte(0))
}

// 8xy1 - OR Vx, Vy
func Test8xy1(t *testing.T) {
	c := newCpuAt(0x8121)
	c.regs[1] = 1
	c.regs[2] = 2
	assert.Equal(t, c.regs[1], byte(1))
	assert.Equal(t, c.regs[2], byte(2))
	c.Step()
	assert.Equal(t, c.regs[1], byte(1|2))
	assert.Equal(t, c.regs[2], byte(2))
}

// 8xy2 - OR Vx, Vy
func Test8xy2(t *testing.T) {
	c := newCpuAt(0x8122)
	c.regs[1] = 1
	c.regs[2] = 2
	assert.Equal(t, c.regs[1], byte(1))
	assert.Equal(t, c.regs[2], byte(2))
	c.Step()
	assert.Equal(t, c.regs[1], byte(1&2))
	assert.Equal(t, c.regs[2], byte(2))
}

// 8xy3 - OR Vx, Vy
func Test8xy3(t *testing.T) {
	c := newCpuAt(0x8123)
	c.regs[1] = 1
	c.regs[2] = 2
	assert.Equal(t, c.regs[1], byte(1))
	assert.Equal(t, c.regs[2], byte(2))
	c.Step()
	assert.Equal(t, c.regs[1], byte(1^2))
	assert.Equal(t, c.regs[2], byte(2))
}

// 8xy4 - ADD Vx, Vy
func Test8xy4Overflow(t *testing.T) {
	c := newCpuAt(0x8124)
	c.regs[1] = 0xFF
	c.regs[2] = 1

	assert.Equal(t, c.regs[0xF], byte(0))

	c.Step()

	assert.Equal(t, c.regs[1], byte(0))
	assert.Equal(t, c.regs[0xF], byte(1))
}

// 8xy4 - ADD Vx, Vy
func Test8xy4ONotverflow(t *testing.T) {
	c := newCpuAt(0x8124)
	c.regs[1] = 0xF0
	c.regs[2] = 1

	assert.Equal(t, c.regs[0xF], byte(0))

	c.Step()

	assert.Equal(t, c.regs[1], byte(0xF0+1))
	assert.Equal(t, c.regs[0xF], byte(0))
}

// 8xy5 - SUB Vx, Vy
func Test8xy5Greater(t *testing.T) {
	c := newCpuAt(0x8125)
	c.regs[1] = 2
	c.regs[2] = 1

	assert.Equal(t, c.regs[0xF], byte(0))

	c.Step()

	assert.Equal(t, c.regs[1], byte(1))
	assert.Equal(t, c.regs[0xF], byte(1))
}

// 8xy5 - SUB Vx, Vy
func Test8xy5NotGreater(t *testing.T) {
	c := newCpuAt(0x8125)
	c.regs[1] = 1
	c.regs[2] = 2

	assert.Equal(t, c.regs[0xF], byte(0))

	c.Step()

	assert.Equal(t, c.regs[1], byte(0xFF))
	assert.Equal(t, c.regs[0xF], byte(0))
}

// 8xy6 - SHR Vx {, Vy}
func Test8xy6(t *testing.T) {
	c := newCpuAt(0x8126)
	c.regs[1] = 1

	c.Step()

	assert.Equal(t, c.regs[0xF], byte(1&0x0001))
	assert.Equal(t, c.regs[1], byte(1/2))
}

// 8xy7 - SUBN Vx, Vy
func Test8xy7Greater(t *testing.T) {
	c := newCpuAt(0x8127)
	c.regs[1] = 1
	c.regs[2] = 2

	assert.Equal(t, c.regs[0xF], byte(0))

	c.Step()

	assert.Equal(t, c.regs[1], byte(1))
	assert.Equal(t, c.regs[0xF], byte(1))
}

// 8xy7 - SUBN Vx, Vy
func Test8xy7NotGreater(t *testing.T) {
	c := newCpuAt(0x8127)
	c.regs[1] = 2
	c.regs[2] = 1

	assert.Equal(t, c.regs[0xF], byte(0))

	c.Step()

	assert.Equal(t, c.regs[1], byte(0xFF))
	assert.Equal(t, c.regs[0xF], byte(0))
}

// 8xyE - SHL Vx {, Vy}
func Test8xyE(t *testing.T) {
	c := newCpuAt(0x812E)
	c.regs[1] = 1

	c.Step()

	assert.Equal(t, c.regs[0xF], byte(1&0x0001))
	assert.Equal(t, c.regs[1], byte(1*2))
}

// 9xy0 - SNE Vx, Vy
func Test9xy0NotEqual(t *testing.T) {
	c := newCpuAt(0x9120)
	c.regs[1] = 1
	c.regs[2] = 2

	oldPC := c.pc

	c.Step()

	assert.Equal(t, c.pc, oldPC+4)
}

// 9xy0 - SNE Vx, Vy
func Test9xy0Equal(t *testing.T) {
	c := newCpuAt(0x9120)
	c.regs[1] = 1
	c.regs[2] = 1

	oldPC := c.pc

	c.Step()

	assert.Equal(t, c.pc, oldPC+2)
}

// Annn - LD I, addr
func TestA000(t *testing.T) {
	c := newCpuAt(0xAFFF)

	c.Step()

	assert.Equal(t, c.i, uint16(0xFFF))
}

// Bnnn - JP V0, addr
func TestB000(t *testing.T) {
	c := newCpuAt(0xBFFF)
	c.regs[0] = 1

	c.Step()

	assert.Equal(t, c.pc, pc(0xFFF+1))
}

// Cxkk - RND Vx, byte
func TestC000(t *testing.T) {
	c := newCpuAt(0xC1FF)

	rand.Seed(1)         // Seed with a Known value
	expectedRandom := 33 // With the expected seed, 33 is the first random.

	c.Step()

	assert.Equal(t, c.regs[1], byte(expectedRandom&0xFF))
}

// TODO: Dxyn - DRW Vx, Vy, nibble

// Ex9E - SKP Vx
func TestEx9E_True(t *testing.T) {
	c := newCpuAt(0xE19E)
	c.keys[1] = true

	oldPC := c.pc + 4

	c.Step()

	assert.Equal(t, pc(oldPC), c.pc)
}

// Ex9E - SKP Vx
func TestEx9E_False(t *testing.T) {
	c := newCpuAt(0xE19E)
	c.keys[1] = false

	oldPC := c.pc + 2

	c.Step()

	assert.Equal(t, pc(oldPC), c.pc)
}

// ExA1 - SKNP Vx
func TestExA1_True(t *testing.T) {
	c := newCpuAt(0xE1A1)
	c.keys[1] = true

	oldPC := c.pc + 2

	c.Step()

	assert.Equal(t, pc(oldPC), c.pc)
}

// ExA1 - SKNP Vx
func TestExA1_False(t *testing.T) {
	c := newCpuAt(0xE1A1)
	c.keys[1] = false

	oldPC := c.pc + 4

	c.Step()

	assert.Equal(t, pc(oldPC), c.pc)
}
