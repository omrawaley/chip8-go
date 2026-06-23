package chip8

import (
	"fmt"
	"math/rand/v2"
)

const (
	pcStart              = 0x200
	instructionsPerFrame = 11
)

type CPU struct {
	Stack      [16]uint16
	I          uint16
	SP         uint16
	PC         uint16
	V          [16]byte
	DelayTimer byte
	SoundTimer byte
}

func NewCPU() *CPU {
	return &CPU{
		PC: pcStart,
	}
}

func (c *CPU) Tick(m *Memory, d *Display, k *Keypad) error {
	for i := 0; i < instructionsPerFrame; i++ {
		if err := c.execute(m, d, k); err != nil {
			return fmt.Errorf("failed to execute instruction: %w", err)
		}
	}

	if c.DelayTimer > 0 {
		c.DelayTimer--
	}
	if c.SoundTimer > 0 {
		c.SoundTimer--
	}

	return nil
}

func (c *CPU) fetch(m *Memory) (uint16, error) {
	hi, err := m.Read(c.PC)
	if err != nil {
		return 0, fmt.Errorf("cannot fetch high byte at address 0x%04X: %w", c.PC, err)
	}

	lo, err := m.Read(c.PC)
	if err != nil {
		return 0, fmt.Errorf("cannot fetch low byte at address 0x%04X: %w", c.PC, err)
	}

	c.PC += 2
	return (uint16(hi) << 8) | uint16(lo), nil
}

func (c *CPU) execute(m *Memory, d *Display, k *Keypad) error {
	opcode, err := c.fetch(m)
	if err != nil {
		return fmt.Errorf("failed to fetch opcode: %w", err)
	}

	x := byte((opcode & 0xF00) >> 8)
	y := byte((opcode & 0xF0) >> 4)
	n := byte(opcode & 0xF)
	nn := byte(opcode & 0xFF)
	nnn := uint16(opcode & 0xFFF)

	switch (opcode & 0xF000) >> 12 {
	case 0x0:
		switch n {
		case 0x0: // CLS
			d.Clear()
		case 0xE: // RET
			c.SP--
			c.PC = c.Stack[c.SP]
		}
	case 0x1: // JP addr
		c.PC = nnn
	case 0x2: // CALL addr
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = nnn
	case 0x3: // SE Vx, byte
		if c.V[x] == nn {
			c.PC += 2
		}
	case 0x4: // SNE Vx, byte
		if c.V[x] != nn {
			c.PC += 2
		}
	case 0x5: // SE Vx, Vy
		if c.V[x] == c.V[y] {
			c.PC += 2
		}
	case 0x6: // LD Vx, byte
		c.V[x] = nn
	case 0x7: // ADD Vx, byte
		c.V[x] += nn
	case 0x8:
		switch n {
		case 0x0: // LD Vx, Vy
			c.V[x] = c.V[y]
		case 0x1: // OR Vx, Vy
			c.V[x] |= c.V[y]
		case 0x2: // AND Vx, Vy
			c.V[x] &= c.V[y]
		case 0x3: // XOR Vx, Vy
			c.V[x] ^= c.V[y]
		case 0x4: // ADD Vx, Vy
			sum := uint16(c.V[x]) + uint16(c.V[y])
			c.V[x] = uint8(sum & 0xFF)
			if sum > 0xFF {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
		case 0x5: // SUB Vx, Vy
			notBorrow := uint16(c.V[x]) >= uint16(c.V[y])
			c.V[x] -= c.V[y]
			if notBorrow {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
		case 0x6: // SHR Vx, Vy
			lsb := c.V[x] & 0x1
			c.V[x] >>= 1
			c.V[0xF] = lsb
		case 0x7: // SUBN Vx, Vy
			borrow := uint16(c.V[y]) >= uint16(c.V[x])
			c.V[x] = c.V[y] - c.V[x]
			if borrow {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
		case 0xE: // SHL Vx, Vy
			msb := (c.V[x] & 0x80) >> 7
			c.V[x] <<= 1
			c.V[0xF] = msb
		}
	case 0x9: // SNE Vx, Vy
		if c.V[x] != c.V[y] {
			c.PC += 2
		}
	case 0xA: // LD I, addr
		c.I = nnn
	case 0xB: // JP V0, addr
		c.PC = nnn + uint16(c.V[0])
	case 0xC: // RND Vx, byte
		c.V[x] = byte(rand.Int()) & nn
	case 0xD: // DRW Vx, Vy, byte
		xPos := c.V[x] % DisplayWidth
		yPos := c.V[y] % DisplayHeight

		c.V[0xF] = 0
		for row := 0; row < int(n); row++ {
			if int(yPos)+row >= DisplayWidth {
				break
			}

			spriteByte, err := m.Read(c.I + uint16(row))
			if err != nil {
				return fmt.Errorf("erroneous read from instruction DRW Vx, Vy, byte: %w", err)
			}
			for col := range 8 {
				if int(xPos)+col >= DisplayWidth {
					break
				}

				if spriteByte&(0x80>>col) == 0 {
					continue
				}

				pixelPos := (int(yPos)+row)*DisplayWidth + (int(xPos) + col)

				setFlag, err := d.GetPixelRaw(pixelPos)
				if err != nil {
					return fmt.Errorf("erroneous display read from instruction DRW Vx, Vy, byte: %w", err)
				}
				if setFlag {
					c.V[0xF] = 1
				}

				currPixel, err := d.GetPixelRaw(pixelPos)
				if err != nil {
					return fmt.Errorf("erroneous display read from instruction DRW Vx, Vy, byte: %w", err)
				}
				if err := d.SetPixel(pixelPos, !currPixel); err != nil {
					return fmt.Errorf("erroneous display write from instruction DRW Vx, Vy, byte: %w", err)
				}
			}
		}
	case 0xE:
		switch n {
		case 0xE: // SKP Vx
			if k.GetKey(int(c.V[x])) {
				c.PC += 2
			}
		case 0x1: // SKNP Vx
			if !k.GetKey(int(c.V[x])) {
				c.PC += 2
			}
		}
	case 0xF:
		switch y {
		case 0x0:
			switch n {
			case 0x7: // LD Vx, DT
				c.V[x] = c.DelayTimer
			case 0xA: // LD Vx, K
			}
		case 0x1:
			switch n {
			case 0x5: // LD DT, Vx
				c.DelayTimer = c.V[x]
			case 0x8: // LD ST, Vx
				c.SoundTimer = c.V[x]
			case 0xE: // ADD I, Vx
				sum := uint32(c.I) + uint32(c.V[x])
				c.I = uint16(sum & 0xFFF)
				if sum > 0xFFF {
					c.V[x] = 1
				} else {
					c.V[y] = 0
				}
			}
		case 0x2: // LD F, Vx
			c.I = uint16((c.V[x] & 0xF) * 5)
		case 0x3: // LD B, Vx
			m.Write(c.I+2, c.V[x]%10)
			m.Write(c.I+1, (c.V[x]/10)%10)
			m.Write(c.I, c.V[x]/100)
		case 0x5: // LD [I], Vx
			for i := 0; i < int(x); i++ {
				if err := m.Write((c.I+uint16(i))&0xFFF, c.V[i]); err != nil {
					return fmt.Errorf("erroneous write from LD [i], Vx: %w", err)
				}
			}
		case 0x6: // LD Vx, [I]
			for i := 0; i < int(x); i++ {
				if c.V[i], err = m.Read((c.I + uint16(i)) & 0xFFF); err != nil {
					return fmt.Errorf("erroneous read from instruction LD [i], Vx: %w", err)
				}
			}
		default:
			return fmt.Errorf("unexpected opcode 0x%04X", opcode)
		}
	}

	return nil
}
