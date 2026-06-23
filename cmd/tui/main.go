package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/omrawaley/chip8-go/internal/chip8"
)

const (
	upperBlock = "▀"
	lowerBlock = "▄"
	fullBlock = "█"
	space = "\u00A0"
)

type emulator struct {
	cpu     *chip8.CPU
	memory  *chip8.Memory
	display *chip8.Display
	keypad  *chip8.Keypad
}

func newEmu() emulator {
	return emulator{
		cpu:     chip8.NewCPU(),
		memory:  chip8.NewMemory(),
		display: chip8.NewDisplay(),
		keypad:  chip8.NewKeypad(),
	}
}

func (e emulator) Init() tea.Cmd {
	return nil
}

func (e emulator) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return e, tea.Quit
		case "h":
			val, _ := e.memory.Read(0)
			e.memory.Write(0, val + 1)
		case "l":
			val, _ := e.memory.Read(0)
			e.memory.Write(0, val - 1)
		}
	}

	return e, nil
}

func (e emulator) View() tea.View {
	var s strings.Builder

	// Increment rows by 2 because 2 rows are represented by one Unicode character
	for row := 0; row < chip8.DisplayHeight; row += 2 {
		for col := 0; col < chip8.DisplayWidth; col++ {
			top, _ := e.display.GetPixelRaw(col + chip8.DisplayWidth*row)
			bottom, _ := e.display.GetPixelRaw(col + chip8.DisplayWidth*(row+1))

			if top && !bottom {
				s.WriteString(upperBlock)
			} else if !top && bottom {
				s.WriteString(lowerBlock)
			} else if top && bottom {
				s.WriteString(fullBlock)
			} else {
				// This is a non-breaking space character meaning the terminal won't strip it
				s.WriteString(space)
			}
		}
		s.WriteString("\n")
	}

	return tea.NewView(s.String())
}

func main() {
	p := tea.NewProgram(newEmu())
	if _, err := p.Run(); err != nil {
		fmt.Printf("an error has occurred: %v", err)
		os.Exit(1)
	}
}
