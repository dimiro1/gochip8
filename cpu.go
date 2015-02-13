// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

import (
	"fmt"
	"math/rand"
	"time"
)

// pc is the program counter
type pc uint16

// Increments the counter
func (p *pc) Increment() {
	(*p) += 2
}

// The stack
type stack struct {
	data [16]pc
	sp   uint16 // Stack Pointer
}

// Push the addr to the stack
func (s *stack) Push(data pc) {
	s.data[s.sp] = data
	s.sp++
}

// Pop a value from stack
func (s *stack) Pop() pc {
	s.sp--
	return s.data[s.sp]
}

// cpu is the chip8 main cpu
type cpu struct {
	regs  [16]byte // Registers v0 - vF
	stack stack    // The Stack
	i     uint16   // Index
	pc    pc       // Program Counter

	sound  sound    // Sound Card
	gfx    gfx      // Graphics
	memory memory   // Memory - 4K
	keys   [16]bool // Key State
	dt     byte     // Delay Timer
	st     byte     // Sound Timer
}

// Chip8 Fonts
var font = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

// newCpu creates a new cpu and initialize it
func newCpu(s sound) cpu {
	rand.Seed(time.Now().UTC().UnixNano())

	c := cpu{
		pc:    0x200,
		i:     0,
		dt:    0,
		st:    0,
		stack: stack{sp: 0},
		sound: s,
	}

	// Load font into memory
	for i, e := range font {
		c.memory.Write(uint16(i), e)
	}

	return c
}

// Step the cpu one instruction at a time
func (c *cpu) Step() {
	opcode := c.memory.ReadWord(uint16(c.pc))
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4

	c.pc.Increment()

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0: // 00E0 - CLS
			c.gfx.Cls()
		case 0x00EE: // 00EE - RET
			c.pc = c.stack.Pop()
		default:
			panic(fmt.Sprintf("Unknown opcode %X", opcode))
		}
	case 0x1000: // 1nnn - JP addr
		c.pc = pc(opcode & 0x0FFF)
	case 0x2000: // 2nnn - CALL addr
		c.stack.Push(c.pc)
		c.pc = pc(opcode & 0x0FFF)
	case 0x3000: // 3xkk - SE Vx, byte
		kk := opcode & 0x00FF

		if c.regs[x] == byte(kk) {
			c.pc.Increment()
		}
	case 0x4000: // 4xkk - SNE Vx, byte
		kk := opcode & 0x00FF

		if c.regs[x] != byte(kk) {
			c.pc.Increment()
		}
	case 0x5000: // 5xy0 - SE Vx, Vy
		if c.regs[x] == c.regs[y] {
			c.pc.Increment()
		}
	case 0x6000: // 6xkk - LD Vx, byte
		kk := opcode & 0x00FF

		c.regs[x] = byte(kk)
	case 0x7000: // 7xkk - ADD Vx, byte
		kk := opcode & 0x00FF

		c.regs[x] += byte(kk)
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000: // 8xy0 - LD Vx, Vy
			c.regs[x] = c.regs[y]
		case 0x0001: // 8xy1 - OR Vx, Vy
			c.regs[x] |= c.regs[y]
		case 0x0002: // 8xy2 - AND Vx, Vy
			c.regs[x] &= c.regs[y]
		case 0x0003: // 8xy3 - XOR Vx, Vy
			c.regs[x] ^= c.regs[y]
		case 0x0004: // 8xy4 - ADD Vx, Vy
			r := uint16(c.regs[x]) + uint16(c.regs[y])

			if r > 0xFF {
				c.regs[0xF] = 1
			} else {
				c.regs[0xF] = 0
			}

			c.regs[x] += c.regs[y]
		case 0x0005: // 8xy5 - SUB Vx, Vy
			if c.regs[x] > c.regs[y] {
				c.regs[0xF] = 1
			} else {
				c.regs[0xF] = 0
			}

			c.regs[x] = c.regs[x] - c.regs[y]
		case 0x0006: // 8xy6 - SHR Vx {, Vy}
			c.regs[0xF] = c.regs[x] & 0x0001
			c.regs[x] = c.regs[x] / 2
		case 0x0007: // 8xy7 - SUBN Vx, Vy
			if c.regs[y] > c.regs[x] {
				c.regs[0xF] = 1
			} else {
				c.regs[0xF] = 0
			}
			c.regs[x] = c.regs[y] - c.regs[x]
		case 0x000E: // 8xyE - SHL Vx {, Vy}
			c.regs[0xF] = c.regs[x] & 0x0001
			c.regs[x] = c.regs[x] * 2
		default:
			panic(fmt.Sprintf("Unknown opcode %X", opcode))
		}
	case 0x9000: // 9xy0 - SNE Vx, Vy
		if c.regs[x] != c.regs[y] {
			c.pc.Increment()
		}
	case 0xA000: // Annn - LD I, addr
		c.i = opcode & 0x0FFF
	case 0xB000: // Bnnn - JP V0, addr
		addr := opcode & 0x0FFF
		c.pc = pc(addr + uint16(c.regs[0]))
	case 0xC000: // Cxkk - RND Vx, byte
		kk := opcode & 0x00FF
		c.regs[x] = byte(rand.Intn(256)) & byte(kk)
	case 0xD000: // Dxyn - DRW Vx, Vy, nibble
		// TODO:  Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
	case 0xE000:
		switch opcode & 0x000F {
		case 0x000E: // Ex9E - SKP Vx
			if c.keys[x] {
				c.pc.Increment()
			}
		case 0x0001: // ExA1 - SKNP Vx
			if !c.keys[x] {
				c.pc.Increment()
			}
		default:
			panic(fmt.Sprintf("Unknown opcode %X", opcode))
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007: // Fx07 - LD Vx, DT
			c.regs[x] = c.dt
		case 0x000A: // Fx0A - LD Vx, K
			// TODO: Wait for a key press, store the value of the key in Vx.
			// All execution stops until a key is pressed, then the value of that key is stored in Vx.
		case 0x0015: // Fx15 - LD DT, Vx
			c.dt = c.regs[x]
		case 0x0018: // Fx18 - LD ST, Vx
			c.st = c.regs[x]
		case 0x001E: // Fx1E - ADD I, Vx
			c.i += uint16(c.regs[x])
		case 0x0029: // Fx29 - LD F, Vx
			c.i = uint16(c.regs[x] * 5) // 5 is the number of rows per character.
		case 0x0033: // Fx33 - LD B, Vx
			num := c.regs[x]

			for i := uint16(3); i > 0; i-- {
				c.memory.Write(c.i+i-1, num%10)
				num /= 10
			}
		case 0x0055: // Fx55 - LD [I], Vx
			for i := uint16(0); i <= x; i++ {
				c.memory.Write(c.i+i, c.regs[i])
			}
		case 0x0065: // Fx65 - LD Vx, [I]
			for i := uint16(0); i <= x; i++ {
				c.regs[i] = c.memory.Read(c.i + i)
			}
		}
	default:
		panic(fmt.Sprintf("Unknown opcode %X", opcode))
	}

	// Timers
	if c.dt > 0 {
		c.dt--
	}

	if c.st > 0 {
		if c.st == 1 {
			c.sound.Beep()
		}
		c.st--
	}
}
