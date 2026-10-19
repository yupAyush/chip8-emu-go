package main

import (
	"chip8-go-emu/internal/chip8"
	"chip8-go-emu/internal/platform"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage chip8 rom file")
		os.Exit(1)
	}

	rom, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Print("cannot read rom:", err)
		os.Exit(1)
	}

	c := chip8.New()
	if err := c.LoadROM(rom); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := platform.Run(c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
