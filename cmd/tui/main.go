package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/omrawaley/chip8-go/internal/chip8"
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
		}
	}

	return e, nil
}

func (e emulator) View() tea.View {
	s := "Hello, CHIP-8!"
	return tea.NewView(s)
}

func main() {
	p := tea.NewProgram(newEmu())
	if _, err := p.Run(); err != nil {
		fmt.Printf("an error has occurred: %v", err)
		os.Exit(1)
	}
}
