package main

import "os"

type MMU struct {
	inbios bool // Flag indicating BIOS is mapped in, BIOS is unmapped with the first instruction above 0x00FF
	// memory regions
	bios [0x100]byte
	rom  [0x8000]byte
	// vram [0x1FFF]byte
	eram [0x2000]byte
	wram [0x3E00]byte
	zram [0x80]byte
}

func (mmu *MMU) readByte(addr uint16) byte {
	switch {
	// BIOS (256b)/ROM0
	case addr <= 0x0FFF:
		if mmu.inbios {
			if addr < 0x100 {
				return mmu.bios[addr]
			} else {
				mmu.inbios = false
			}
		}
		return mmu.rom[addr]

	// ROM0
	case addr >= 0x1000 && addr <= 0x3FFF:
		return mmu.rom[addr]
	// ROM1 (unbanked) (16k)
	case addr >= 0x4000 && addr <= 0x7FFF:
		return mmu.rom[addr]
	// External RAM (8k)
	case addr >= 0xA000 && addr <= 0xBFFF:
		return mmu.eram[addr&0x1FFF]
	// Working RAM (8k)
	case addr >= 0xC000 && addr <= 0xDFFF:
		return mmu.wram[addr&0x1FFF]
	// Working RAM shadow
	case addr >= 0xE000 && addr <= 0xFDFF:
		return mmu.wram[addr&0x1FFF]
	case addr >= 0xFF00 && addr <= 0xFFFF:
		if addr >= 0xFF80 {
			return mmu.zram[addr&0x7F]
		} else {
			return 0
		}
	}

	return 0
}

func (mmu *MMU) readWord(addr uint16) uint16 {
	return uint16(mmu.readByte(addr)) + (uint16(mmu.readByte(addr+1)) << 8)
}

func (mmu *MMU) load(path string) {
	dat, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	copy(mmu.rom[:], dat)
}
