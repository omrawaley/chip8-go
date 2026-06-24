package chip8

const NumKeys = 16

const (
	KeyOne = iota
	KeyTwo
	KeyThree
	KeyC
	KeyFour
	KeyFive
	KeySix
	KeyD
	KeySeven
	KeyEight
	KeyNine
	KeyE
	KeyA
	KeyZero
	KeyB
	KeyF
)

type Keypad struct {
	keys [NumKeys]bool
}

func NewKeypad() *Keypad {
	return &Keypad{}
}

func (k *Keypad) SetKey(key int, pressed bool) {
	k.keys[key] = pressed
}

func (k *Keypad) GetKey(key int) bool {
	return k.keys[key]
}
