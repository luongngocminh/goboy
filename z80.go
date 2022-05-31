package main

const (
	F_ZERO = 0x80 // Last operation resulted in 0
	F_OP   = 0x40 // last operation was a subtraction
	F_HC   = 0x20 // in the result of the last operation, the lower half of the byte overflowed past 15
	F_CA   = 0x10 // the last operation produced a result over 255 (for additions) or under 0 (for subtractions)
)

type Clock struct {
	m byte
	t byte
}

type Register struct {
	a, b, c, d, e, h, l, f byte   // 8-bit registers
	pc, sp                 uint16 // 16-bit registers
	m, t                   byte   // Clock for last instr
}

type IZ80 interface {
	checkZeroFlag(value int16)  // Check should set zero flag
	checkCarryFlag(value int16) // Check should set carry flag
	setSubtractFlag()

	tickMTime(m, t byte)

	ADDre() // Add E to A, leaving result in A (ADD A, E)
	CPrb()  // Compare B to A, setting flags (CP A, B)
	NOP()   // No-operation (NOP)
}

type Z80 struct {
	clock Clock
	r     Register
}

// HELPERS

func (cpu *Z80) checkZeroFlag(value int16) {
	if value&0xff == 0 {
		cpu.r.f |= F_ZERO
	}
}

func (cpu *Z80) checkCarryFlag(value int16) {
	if value > 0xff || value < 0 {
		cpu.r.f |= F_CA
	}
}

func (cpu *Z80) setSubtractFlag() {
	cpu.r.f |= 0x40
}

func (cpu *Z80) tickMTime(m, t byte) {
	cpu.clock.m = m
	cpu.clock.t = t
}

func mapToByte(value int16) int16 {
	return value & 0xff
}

// INSTRUCTIONS

func (cpu *Z80) ADDre() {
	tmp := int16(cpu.r.a) + int16(cpu.r.e)
	cpu.r.f = 0
	cpu.checkZeroFlag(tmp)
	cpu.checkCarryFlag(tmp)
	cpu.r.a = byte(mapToByte(tmp))
	cpu.tickMTime(1, 4)
}

func (cpu *Z80) CPrb() {
	i := int16(cpu.r.a)
	i -= int16(cpu.r.b)
	cpu.setSubtractFlag()
	cpu.checkZeroFlag(i)
	cpu.checkCarryFlag(i) // incase of undeflow
	cpu.tickMTime(1, 4)
}

func (cpu *Z80) NOP() {
	cpu.tickMTime(1, 4)
}
