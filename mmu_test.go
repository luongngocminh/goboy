package main

import (
	"testing"
)

func TestMMU_readByte(t *testing.T) {
	type args struct {
		addr uint16
	}

	mmu := MMU{inbios: true}
	for i := 0; i < int(len(mmu.bios)); i++ {
		mmu.bios[i] = 0
	}
	for i := 0; i < int(len(mmu.rom)); i++ {
		mmu.rom[i] = 1
	}
	for i := 0; i < int(len(mmu.eram)); i++ {
		mmu.eram[i] = 2
	}
	for i := 0; i < int(len(mmu.wram)); i++ {
		mmu.wram[i] = 3
	}
	for i := 0; i < int(len(mmu.zram)); i++ {
		mmu.zram[i] = 4
	}

	tests := []struct {
		name string
		mmu  *MMU
		args args
		want byte
	}{
		{"MMU: Test BIOS Memory address at 0x0000", &mmu, args{0x0000}, 0},
		{"MMU: Test BIOS Memory address at 0x00FF", &mmu, args{0x00FF}, 0},
		{"MMU: Test ROM Memory address at 0x0100", &mmu, args{0x0100}, 1}, // Trigger BIOS removal
		{"MMU: Test ROM Memory address at 0x0000", &mmu, args{0x0000}, 1},
		{"MMU: Test ROM Memory address at 0x7FFF", &mmu, args{0x7FFF}, 1},
		{"MMU: Test ERAM address at 0xA000", &mmu, args{0xA000}, 2},
		{"MMU: Test ERAM address at 0xBFFF", &mmu, args{0xBFFF}, 2},
		{"MMU: Test WRAM address at 0xC000", &mmu, args{0xC000}, 3},
		{"MMU: Test WRAM address at 0xDFFF", &mmu, args{0xDFFF}, 3},
		{"MMU: Test WRAM address at 0xC000", &mmu, args{0xE000}, 3},
		{"MMU: Test WRAM address at 0xDFFF", &mmu, args{0xFDFF}, 3},
		{"MMU: Test ZRAM address at 0xFF00", &mmu, args{0xFF00}, 0},
		{"MMU: Test ZRAM address at 0xFF80", &mmu, args{0xFF80}, 4},
		{"MMU: Test ZRAM address at 0xFFF0", &mmu, args{0xFFF0}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mmu.readByte(tt.args.addr); got != tt.want {
				t.Errorf("MMU.readByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMMU_readWord(t *testing.T) {
	type args struct {
		addr uint16
	}
	tests := []struct {
		name string
		mmu  *MMU
		args args
		want uint16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mmu.readWord(tt.args.addr); got != tt.want {
				t.Errorf("MMU.readWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMMU_load(t *testing.T) {
	type args struct {
		path string
	}
	mmu := MMU{inbios: true}

	tests := []struct {
		name string
		mmu  *MMU
		args args
	}{
		{"MMU: Test load ROM", &mmu, args{"./tetris.gb"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mmu.load(tt.args.path)
		})
	}
}
