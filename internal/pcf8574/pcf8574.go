package pcf8574

import (
	"fmt"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

const BaseAddr = 0b0010_0000

// PCF8574 Remote 8-Bit I/O Expander with the following pinout:
// 7 (MSB) | 6  | 5  | 4  | 3  | 2  | 1  | 0 (LSB) |
// P7      | P6 | P5 | P4 | P3 | P2 | P1 | P0      |
// these pins can be written or read.
type PCF8574 struct {
	Addr uint16
	Bus  i2c.BusCloser
}

func New(i2c string) (*PCF8574, error) {
	return NewWithAddr(i2c, 0, 0, 0)
}

func NewWithAddr(i2c string, a0, a1, a2 uint16) (*PCF8574, error) {
	b, err := i2creg.Open(i2c)
	if err != nil {
		return nil, fmt.Errorf("cannot open i2c device %s %w", i2c, err)
	}

	return &PCF8574{
		Addr: BaseAddr | a0 | a1<<1 | a2<<2,
		Bus:  b,
	}, nil
}

// Write writes the given pins on PCF8574.
func (p *PCF8574) Write(pins byte) error {
	if err := p.Bus.Tx(p.Addr, []byte{pins}, nil); err != nil {
		return fmt.Errorf("cannot communicate with i2c device %w", err)
	}

	return nil
}
