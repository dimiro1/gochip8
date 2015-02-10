// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

// reg represents a processor 8 bit register
type reg byte

// reg16 represents a processor 16 bit register
type reg16 uint16

type cpu struct {
	regs  [16]reg    // Registers v0 - vF
	stack [16]uint16 // The Stack
	i     reg16      // Index
	pc    reg16      // Program Counter
	sp    reg16      // Stack Pointer

	opcode byte          // The current opcode
	gfx    [64 * 32]byte // Graphics
	memory [4096]byte    // Memory - 4K
	keys   [16]byte      // Key State
}

func newCpu() cpu {
	c := cpu{
		pc:     0x200,
		opcode: 0,
		i:      0,
		sp:     0,
	}

	// Load font
	for i := 0; i < 80; i++ {
		c.memory[i] = 0 // fontset
	}

	return c
}

func (c *cpu) Step() {
	c.opcode = c.memory[c.pc]<<8 | c.memory[c.pc+1]

	switch c.opcode {
	}
}
