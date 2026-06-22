package chip8

import "fmt"

const (
	DisplayWidth  = 64
	DisplayHeight = 32
)

type Display struct {
	data [DisplayWidth * DisplayHeight]bool
}

func NewDisplay() *Display {
	return &Display{}
}

func (d *Display) Clear() {
	for i := range d.data {
		d.data[i] = false
	}
}

func (d *Display) GetPixel(index int) (bool, error) {
	if index >= DisplayWidth*DisplayHeight {
		return false, fmt.Errorf("display get pixel index %v out of bounds", index)
	}

	return d.data[index], nil
}

func (d *Display) SetPixel(index int, val bool) error {
	if index >= DisplayWidth*DisplayHeight {
		return fmt.Errorf("display set pixel index %v out of bounds", index)
	}

	d.data[index] = val
	return nil
}

func (d *Display) GetData() *[DisplayWidth * DisplayHeight]bool {
	return &d.data
}
