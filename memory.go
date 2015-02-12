// Copyright 2015 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gochip8

type memory [4096]byte

// Read reads a byte from memory
func (m *memory) Read(addr uint16) byte {
	return m[addr]
}

// ReadWord reads a 16 bit unsigned from memory
func (m *memory) ReadWord(addr uint16) uint16 {
	return uint16(m.Read(addr))<<8 | uint16(m.Read(addr+1))
}

// Write writes data to address specified in addr
func (m *memory) Write(addr uint16, data byte) {
	m[addr] = data
}

// WriteWord writes two bytes in memory
func (m *memory) WriteWord(addr, data uint16) {
	m.Write(addr, byte(data>>8))
	m.Write(addr+1, byte(data&0xFF))
}
