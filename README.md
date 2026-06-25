# chip-8 go

<img width="1000" height="" alt="chip8-go" src="https://github.com/user-attachments/assets/680c982f-ae71-4226-8e11-4978a045475c" />

A CHIP-8 emulator written in Go, featuring a TUI interface.

This project was mostly written for myself so that I could learn Go and how to build TUI apps. I also built this for Hack Club Macondo.

## Demo

https://github.com/user-attachments/assets/0aaba8f6-857f-4ec8-8881-251a7e12c83a

## Features

- Blazing fast terminal rendering using only Unicode characters
- Virtual keypad support for both legacy terminals (e.g. `xterm`) and modern terminals (e.g. `kitty`)
- Built-in file browser for selecting ROMs
- Pause functionality
- Interactive help menu
- Customizable display palette

## Building

### Manual

Firstly, install [Go](https://go.dev/doc/install) on your machine.

Then, to compile chip8-go, run `go build ./cmd/tui`.

### Automatic

To automatically compile binaries for all operating systems and architectures, run the `build.bash` script:

```bash
chmod +x ./build.bash
./build.bash
```

A `clean.bash` script is also provided to delete all binaries:

```bash
chmod +x ./clean.bash
./clean.bash
```

## Usage

Now you can execute the binary to launch the program (e.g. `./chip-8`).

It takes an optional argument specifying the path to a ROM.

## Custom Frontends

The `internal` folder contains all core CHIP-8 logic and data. You can use the `chip8` package to implement any type of frontend you want (not limited to just TUI).

See `cmd/tui/main.go` to understand how to implement your own frontend.

## Authors

- [@omrawaley](https://www.github.com/omrawaley)

## Related

Here are some related projects of mine:

- [Game Boy Emulator (C++)](https://github.com/omrawaley/gameboy-emulator)
- [CHIP-8 emulator, debugger, and disassembler (C++)](https://github.com/omrawaley/chip-8-emulator)
- [CHIP-8 emulator rewrite (Rust)](https://github.com/omrawaley/chip8-emulator-rust)

## License

[MIT](https://choosealicense.com/licenses/mit/)
