package main

import (
	"time"

	"github.com/pterm/pterm"
	"github.com/r523/khandaab/internal/pcf8574"
	"github.com/r523/khandaab/internal/rfid"
	"github.com/r523/khandaab/internal/servo"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

const (
	Gain      = 5
	AllowedID = "0cdb074999"

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

	p, err := pcf8574.New("/dev/i2c-1")
	if err != nil {
		pterm.Error.Printf("cannot open the pcf8574 %s\n", err)

		return
	}

	if err := p.Write(0b0000_0000); err != nil {
		pterm.Error.Printf("cannot set the pcf8574 pins %s\n", err)

		return
	}

	time.Sleep(BuzzTimeout)

	if err := p.Write(0b1111_0000); err != nil {
		pterm.Error.Printf("cannot set the pcf8574 pins %s\n", err)

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
