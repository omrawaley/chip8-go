package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/filepicker"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/omrawaley/chip8-go/internal/chip8"
)

const (
	upperBlock = "▀"
	lowerBlock = "▄"
	fullBlock  = "█"
	space      = "\u00A0"
)

type styles struct {
	display lipgloss.Style
}

func newStyles() (s styles) {
	s.display = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#000000"))

	return s
}

type TickMsg time.Time

type emulator struct {
	isRunning    bool
	cpu          *chip8.CPU
	memory       *chip8.Memory
	display      *chip8.Display
	keypad       *chip8.Keypad
	styles       styles
	filepicker   filepicker.Model
	selectedFile string
}

func newEmu() emulator {
	e := emulator{
		cpu:     chip8.NewCPU(),
		memory:  chip8.NewMemory(),
		display: chip8.NewDisplay(),
		keypad:  chip8.NewKeypad(),
		styles:  newStyles(),
	}

	e.filepicker = filepicker.New()
	e.filepicker.AllowedTypes = []string{".bin", ".ch8", ".c8"}
	e.filepicker.CurrentDirectory, _ = os.UserHomeDir()

	return e
}

func (e *emulator) reset() {
	e.isRunning = false
	e.cpu = chip8.NewCPU()
	e.memory = chip8.NewMemory()
	e.display = chip8.NewDisplay()
	e.keypad = chip8.NewKeypad()
	e.selectedFile = ""
}

func (e *emulator) loadProgram(name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return fmt.Errorf("failed to read file %v", name)
	}

	for i := range len(data) {
		e.memory.Write(uint16(i+chip8.PCStart), data[i])
	}

	return nil
}

func poll() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (e emulator) Init() tea.Cmd {
	return tea.Batch(poll(), e.filepicker.Init())
}

func (e *emulator) chip8KeyFor(s string) (int, bool) {
	switch s {
	case "1":
		return chip8.KeyOne, true
	case "2":
		return chip8.KeyTwo, true
	case "3":
		return chip8.KeyThree, true
	case "4":
		return chip8.KeyC, true
	case "q":
		return chip8.KeyFour, true
	case "w":
		return chip8.KeyFive, true
	case "e":
		return chip8.KeySix, true
	case "r":
		return chip8.KeyD, true
	case "a":
		return chip8.KeySeven, true
	case "s":
		return chip8.KeyEight, true
	case "d":
		return chip8.KeyNine, true
	case "f":
		return chip8.KeyE, true
	case "z":
		return chip8.KeyA, true
	case "x":
		return chip8.KeyZero, true
	case "c":
		return chip8.KeyB, true
	case "v":
		return chip8.KeyF, true
	}

	return 0, false
}

func (e emulator) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return e, tea.Quit
		case "k":
			e.styles.display = e.styles.display.
				Background(lipgloss.Color("#ff0000"))
		case "m":
			e.reset()
		case "1", "2", "3", "4", "q", "w", "e", "r", "a", "s", "d", "f", "z", "x", "c", "v":
			key, ok := e.chip8KeyFor(msg.String())
			if !ok {
				return e, nil
			} else {
				e.keypad.SetKey(key, true)
			}
		}

	case tea.KeyReleaseMsg:
		switch msg.String() {
		case "1", "2", "3", "4", "q", "w", "e", "r", "a", "s", "d", "f", "z", "x", "c", "v":
			key, ok := e.chip8KeyFor(msg.String())
			if !ok {
				return e, nil
			} else {
				e.keypad.SetKey(key, false)
			}
		}

	case TickMsg:
		if e.isRunning {
			err := e.cpu.Tick(e.memory, e.display, e.keypad)
			if err != nil {
				fmt.Println(err)
				return e, tea.Quit
			}
		}
		return e, poll()
	}

	var cmd tea.Cmd
	if !e.isRunning {
		e.filepicker, cmd = e.filepicker.Update(msg)
		if didSelect, path := e.filepicker.DidSelectFile(msg); didSelect {
			e.selectedFile = path
		}
	}

	if !e.isRunning && e.selectedFile != "" {
		e.loadProgram(e.selectedFile)
		e.isRunning = true
	}

	return e, cmd
}

func (e emulator) View() tea.View {
	var v tea.View
	v.KeyboardEnhancements.ReportEventTypes = true
	v.AltScreen = true

	var s strings.Builder
	if !e.isRunning {
		if e.selectedFile == "" {
			s.WriteString("Please pick a ROM file")
		} else {
			s.WriteString("Selected file: " + e.filepicker.Styles.Selected.Render(e.selectedFile))
		}
		s.WriteString("\n\n" + e.filepicker.View() + "\n")
	} else {
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

	}

	content := e.styles.display.Render(s.String())
	v.SetContent(content)
	return v
}

func main() {
	e := newEmu()

	args := os.Args
	if len(args) >= 2 {
		e.loadProgram(args[1])
	}

	p := tea.NewProgram(e)
	if _, err := p.Run(); err != nil {
		fmt.Printf("an error has occurred: %v", err)
		os.Exit(1)
	}
}
