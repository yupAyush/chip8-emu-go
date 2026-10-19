package chip8

import (
	"fmt"
	"math/rand"
)

func (c *Chip8) Cycle() {
	opcode := uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
	c.PC += 2

	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := (opcode & 0x000F)
	nn := byte(opcode & 0x00FF)
	nnn := opcode & 0x0FFF

	switch opcode & 0xF000 {
	case 0x0000:
		switch nn {
		case 0xE0:
			c.Display = [ScreenWidth * ScreenHeight]bool{}

		case 0xEE: // return from subroutines
			c.SP--
			c.PC = c.Stack[c.SP]

		default:
			c.unknown(opcode)

		}
	case 0x1000: // jump to NNN
		c.PC = nnn

	case 0x2000: // call subroutine
		c.Stack[c.SP] = c.PC // store current pc then go to subrountine
		c.SP++
		c.PC = nnn

	case 0x3000: // skip next if vx==nn
		if c.V[x] == nn {
			c.PC += 2
		}

	case 0x4000: // skip if vx!=nn
		if c.V[x] != nn {
			c.PC += 2
		}

	case 0x5000: // skip next if vx==vy
		if c.V[x] == c.V[y] {
			c.PC += 2
		}

	case 0x6000: // VX=NN
		c.V[x] = nn
	case 0x7000: // 7XNN: VX+=NN
		c.V[x] += nn

	case 0x8000: // for math and logic 0x8XY
		switch n {
		case 0x0:
			c.V[x] = c.V[y]

		case 0x1:
			c.V[x] |= c.V[y]

		case 0x2:
			c.V[x] &= c.V[y]
		case 0x3:
			c.V[x] ^= c.V[y]
		case 0x4: // vx+=vy , vf = carry
			sum := uint16(c.V[x]) + uint16(c.V[y])
			c.V[x] = byte(sum)
			if sum > 0xFF {
				c.V[0xF] = 1

			} else {
				c.V[0xF] = 0

			}
		case 0x5: // vx-=vy , vf=1 if no borrow
			vx, vy := c.V[x], c.V[y] // save first: x or y could be F
			c.V[x] = vx - vy
			if vx >= vy {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}

		case 0x6: // VX >>-1 vf = bit shifted out
			vx := c.V[x]
			c.V[x] = vx >> 1
			c.V[0xF] = vx & 1
		case 0x7: // vx=vy-vx vf =1 if no borrow
			vx, vy := c.V[x], c.V[y]
			c.V[x] = vy - vx
			if vy >= vx {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
		case 0xE: //vx <<=1 vf = bit shiftted out
			vx := c.V[x]
			c.V[x] = vx << 1
			c.V[0xF] = vx >> 7
		default:
			c.unknown(opcode)

		}

	case 0x9000: //skip if next if v !=vy
		if c.V[x] != c.V[y] {
			c.PC += 2
		}

	case 0xA000: // ANNN: I =NNN
		c.I = nnn

	case 0xB000:
		c.PC = nnn + uint16(c.V[0])

	case 0xC000:
		c.V[x] = byte(rand.Intn(256)) & nn

	case 0xD000: // drawing
		c.drawSprite(x, (opcode&0x00F0)>>4, int(opcode&0x000F))

	case 0xE000: // new ex___ keyboard
		switch nn {
		case 0x9E: // skip next if key is pressed
			if c.Keys[c.V[x]&0xF] {
				c.PC += 2

			}
		case 0xA1: // skip if vx is not pressed
			if !c.Keys[c.V[x]&0xF] {
				c.PC += 2

			}
		default:
			c.unknown(opcode)

		}

	case 0xF000: // new fx__ timer , memory font
		switch nn {
		case 0x07: // vx = delay timer
			c.V[x] = c.DelayTimer

		case 0x0A: // wait for key press , store it in vx
			pressed := false
			for i, down := range c.Keys {
				if down {
					c.V[x] = byte(i)
					pressed = true
					break

				}
			}

			if !pressed {
				c.PC -= 2 // no key yet run this same instruction again

			}
		case 0x15: // delaytimer =vx
			c.DelayTimer = c.V[x]

		case 0x18: // sound timer = vx
			c.SounderTimer = c.V[x]

		case 0x1E: // I+=vx
			c.I += uint16(c.V[x])

		case 0x29: // I = address of font character vx( 5 byte eacg)
			c.I = FontAddress + uint16(c.V[x]&0xF)*5

		case 0x33: // BCD: hundreds, tens, ones of VX at I, I+1, I+2
			v := c.V[x]
			c.Memory[c.I] = v / 100
			c.Memory[c.I+1] = (v / 10) % 10
			c.Memory[c.I+2] = v % 10

		case 0x55: // store V0..VX into memory at I
			for i := uint16(0); i <= x; i++ {
				c.Memory[c.I+i] = c.V[i]
			}

		case 0x65: // load V0..VX from memory at I
			for i := uint16(0); i <= x; i++ {
				c.V[i] = c.Memory[c.I+i]
			}

		default:
			c.unknown(opcode)

		}

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
