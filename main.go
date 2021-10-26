package main

import (
	"time"

	"github.com/pterm/pterm"
	"github.com/r523/khandaab/internal/rfid"
	"github.com/r523/khandaab/internal/servo"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

const (
	Gain      = 5
	AllowedID = "0cdb074999"

	I2CAddr = 0x20

	BuzzTimeout = 1 * time.Second

	ServoDutyNumerator   gpio.Duty = 2
	ServoDutyDenominator gpio.Duty = 5
	ServoFreq                      = 10 * physic.Hertz
	ServoTimeout                   = 10 * time.Second
)

func main() {
	if err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("Khan", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("daab", pterm.NewStyle(pterm.FgLightRed)),
	).Render(); err != nil {
		_ = err
	}

	if _, err := host.Init(); err != nil {
		pterm.Error.Printf("host initiation failed %s\n", err)

		return
	}

	var (
		ResetPin gpio.PinOut = rpi.P1_13
		IRQPin   gpio.PinIn  = rpi.P1_11
	)

	rid, err := rfid.Setup("/dev/spidev0.0", ResetPin, IRQPin, Gain)
	if err != nil {
		pterm.Error.Printf("cannot create rfid device %s\n", err)

		return
	}

	pterm.Info.Println("Started rfid reader.")

	id := rfid.ReadRFIDWithRetries(rid)

	pterm.Info.Println(id)

	if id != AllowedID {
		pterm.Error.Printf("you cannot have access %s\n", id)

		return
	}

	b, err := i2creg.Open("/dev/i2c-1")
	if err != nil {
		pterm.Error.Printf("cannot open i2c device %s\n", err)

		return
	}
	defer b.Close()

	// PCF8574 Remote 8-Bit I/O Expander
	// 7 (MSB) | 6  | 5  | 4  | 3  | 2  | 1  | 0 (LSB) |
	// P7      | P6 | P5 | P4 | P3 | P2 | P1 | P0      |
	// P7 is attached to buzzer
	if err := b.Tx(I2CAddr, []byte{0x00}, nil); err != nil {
		pterm.Error.Printf("cannot communicate with i2c device %s\n", err)

		return
	}

	time.Sleep(BuzzTimeout)

	if err := b.Tx(I2CAddr, []byte{0xF0}, nil); err != nil {
		pterm.Error.Printf("cannot communicate with i2c device %s\n", err)

		return
	}

	s := servo.New(rpi.P1_33, ServoDutyNumerator, ServoDutyDenominator, ServoFreq)

	if err := s.Start(); err != nil {
		pterm.Error.Printf("cannot start the servo %s", err)

		return
	}

	time.Sleep(ServoTimeout)

	_ = s.Stop()
}
