package platform

import (
	chip8 "chip8-go-emu/internal/chip8"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	cyclesPerFrame = 10 // 600 instructions per second at 60 fps
	scale          = 10
)

//	CHIP-8 keypad      Your keyboard
//	1 2 3 C            1 2 3 4
//	4 5 6 D            Q W E R
//	7 8 9 E            A S D F
//	A 0 B F            Z X C V

var keyMap = [16]ebiten.Key{

	ebiten.KeyX, // 0
	ebiten.Key1, // 1
	ebiten.Key2, // 2
	ebiten.Key3, // 3
	ebiten.KeyQ, // 4
	ebiten.KeyW, // 5
	ebiten.KeyE, // 6
	ebiten.KeyA, // 7
	ebiten.KeyS, // 8
	ebiten.KeyD, // 9
	ebiten.KeyZ, // A
	ebiten.KeyC, // B
	ebiten.Key4, // C
	ebiten.KeyR, // D
	ebiten.KeyF, // E
	ebiten.KeyV, // F
}

type Game struct {
	c      *chip8.Chip8
	pixels []byte // RGBA biffer 4 bytes per pixel

}

func (g *Game) Update() error {
	for i, key := range keyMap {
		g.c.Keys[i] = ebiten.IsKeyPressed(key)
	}

	for i := 0; i < cyclesPerFrame; i++ {
		g.c.Cycle()
	}

	g.c.TickTimers()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, on := range g.c.Display {
		var r, gr, b byte // black is default
		if on {
			r, gr, b = 0, 255, 100
		}
		p := i * 4
		g.pixels[p+0] = r
		g.pixels[p+1] = gr
		g.pixels[p+2] = b
		g.pixels[p+3] = 255 // alpha for opacity

	}
	screen.WritePixels(g.pixels)

}

func (g *Game) Layout(outW, outH int) (int, int) {
	return chip8.ScreenWidth, chip8.ScreenHeight
}

func Run(c *chip8.Chip8) error {
	g := &Game{
		c:      c,
		pixels: make([]byte, chip8.ScreenWidth*chip8.ScreenHeight*4),
	}

	ebiten.SetWindowSize(chip8.ScreenWidth*scale, chip8.ScreenHeight*scale)
	ebiten.SetWindowTitle("chip-8-go-emu")

	return ebiten.RunGame(g)
}
