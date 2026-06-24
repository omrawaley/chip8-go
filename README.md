# chip-8 go

A CHIP-8 emulator written in Go, featuring a TUI interface.

This project was mostly written for myself so that I could learn Go and how to build TUI apps.

## Demo

<video src="chip8-go-demo.mp4" width="50%" controls></video>

## Features

- Blazing fast terminal rendering using only Unicode characters
- Virtual keypad support for both legacy terminals (e.g. `xterm`) and modern terminals (e.g. `kitty`)
- Built-in file browser for selecting ROMs
- Pause functionality
- Interactive help menu
- Customizable display palette

## Installation

Firstly, install [Go](https://go.dev/doc/install) on your machine.

Then, to compile chip8-go, run `go build cmd/tui/main.go`.

Now you can execute the binary to launch the program (e.g. `./chip-8`).

## Custom Frontends

The `internal` folder contains all core CHIP-8 logic and data. You can use the `chip8` package to implement any type of frontend you want (not limited to just TUI).

See `cmd/tui/main.go` to understand how to implement your own frontend.

## Authors

- [@omrawaley](https://www.github.com/omrawaley)

## Related

Here are some related projects of mine:

[Game Boy Emulator (C++)](https://github.com/omrawaley/gameboy-emulator)
[CHIP-8 emulator, debugger, and disassembler (C++)](https://github.com/omrawaley/chip-8-emulator)
[CHIP-8 emulator rewrite (Rust)](https://github.com/omrawaley/chip8-emulator-rust) (yes I rewrote it in Rust lol)

## License

[MIT](https://choosealicense.com/licenses/mit/)
