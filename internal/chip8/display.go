//
// Copyright (c) 2026 Om Rawaley. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for details.
//

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

func (d *Display) GetPixelRaw(index int) (bool, error) {
	if index >= DisplayHeight*DisplayWidth {
		return false, fmt.Errorf("display get pixel raw index %v out of bounds", index)
	}

	return d.data[index], nil
}

func (d *Display) GetPixel(x int, y int) (bool, error) {
	if x > DisplayWidth {
		return false, fmt.Errorf("display get pixel x %v out of bounds", x)
	}
	if y > DisplayHeight {
		return false, fmt.Errorf("display get pixel y %v out of bounds", y)
	}

	return d.data[(y * DisplayWidth) + x], nil
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
