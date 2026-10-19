package chip8

import "fmt"

const (
	MemorySize = 4096

	StartAddress = 0x200

	FontAddress = 0x50

	ScreenWidth = 64

	ScreenHeight = 32
)

type Chip8 struct {
	Memory [MemorySize]byte
	V      [16]byte
	I      uint16
	PC     uint16

	Stack [16]uint16
	SP    uint16

	DelayTimer   byte
	SounderTimer byte

	Display [ScreenWidth * ScreenHeight]bool

	Keys [16]bool
}

func New() *Chip8 {
	c := &Chip8{PC: StartAddress}
	copy(c.Memory[FontAddress:], fontset[:])
	return c
}

func (c *Chip8) LoadROM(data []byte) error {
	if len(data) > MemorySize-StartAddress {
		return fmt.Errorf("Rom is too large; %d bytes", len(data))
	}
	copy(c.Memory[StartAddress:], data)
	return nil
}

func (c *Chip8) TickTimers() {
	if c.DelayTimer > 0 {
		c.DelayTimer--
	}

	if c.SounderTimer > 0 {
		c.SounderTimer--
	}
}
