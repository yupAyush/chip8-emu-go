# chip8-go

A CHIP-8 emulator written in Go, with graphics and input by [Ebitengine](https://ebitengine.org).

<p align="center">
  <img src="https://media2.giphy.com/media/v1.Y2lkPTc5MGI3NjExNGhhaHZuZjEydnRmMjUybWk5c2pjOXlyODJlMG1wcXM5ZXdvbm9ldiZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/VLzbEtlbwJUFljcRbf/giphy.gif" alt="demo" width="480">
</p>

## Specs

| Part | Spec |
|---|---|
| Memory | 4096 bytes (`0x000`-`0xFFF`), ROMs load at `0x200` |
| Registers | 16 x 8-bit (`V0`-`VF`), `VF` is the flag register |
| Index register `I` | 16-bit |
| Program counter `PC` | 16-bit |
| Stack | 16 levels x 16-bit |
| Timers | Delay and sound, 8-bit, count down at 60 Hz |
| Display | 64 x 32, monochrome, sprites drawn with XOR |
| Input | 16-key hex keypad |
| Font | 16 characters (0-F), 5 bytes each, stored at `0x050` |
| Instructions | 35 opcodes, 2 bytes each, big-endian |
| CPU speed | ~600 instructions/second (10 per frame at 60 fps) |

## Run

```bash
go mod 
go run ./cmd/chip8 roms/your_rom.ch8
```

## Keys

```
CHIP-8 keypad      Keyboard
1 2 3 C            1 2 3 4
4 5 6 D            Q W E R
7 8 9 E            A S D F
A 0 B F            Z X C V
```
