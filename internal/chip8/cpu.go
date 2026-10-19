package chip8

import "fmt"

func (c *Chip8) Cycle() {
	opcode := uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
	c.PC += 2

	x := (opcode & 0x0F00) >> 8
	nn := byte(opcode & 0x00FF)
	nnn := opcode & 0x0FFF

	switch opcode & 0xF000 {
	case 0x000:
		switch nn {
		case 0xE0:
			c.Display = [ScreenWidth * ScreenHeight]bool{}

		default:
			c.unknown(opcode)

		}
	case 0x1000: // jump to NNN
		c.PC = nnn

	case 0x6000: // VX=NN
		c.V[x] = nn
	case 0x7000: // 7XNN: VX+=NN
		c.V[x] += nn

	case 0xA000: // ANNN: I =NNN
		c.I = nnn
	case 0xD000: // drawing
		c.drawSprite(x, (opcode&0x00F0)>>4, int(opcode&0x000F))
	default:
		c.unknown(opcode)
	}

}

func (c *Chip8) drawSprite(x, y uint16, n int) { //draws n rows sprite from (v[x],v[y])
	px := int(c.V[x]) % ScreenWidth
	py := int(c.V[y]) % ScreenHeight
	c.V[0xF] = 0

	for row := 0; row < n; row++ {
		spriteByte := c.Memory[int(c.I)+row]

		for col := 0; col < 8; col++ {
			if px+col >= ScreenWidth || py+row >= ScreenHeight {
				continue
			}

			bit := (spriteByte >> (7 - col)) & 1 // building bit by bit
			if bit == 1 {
				idx := (py+row)*ScreenWidth + (px + col)
				if c.Display[idx] {
					c.V[0xF] = 1
				}
				c.Display[idx] = !c.Display[idx] // flipping the bits
			}

		}
	}
}

func (c *Chip8) unknown(opcode uint16) {
	fmt.Printf("unknown opcode:%04X (PC=%03X) ", opcode, c.PC-2)
}
