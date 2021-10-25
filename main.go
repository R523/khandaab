package main

import (
	"encoding/hex"
	"time"

	"github.com/pterm/pterm"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/mfrc522"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

const (
	ReadTimeout = 5 * time.Minute
	Gain        = 5
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

	// get the first available spi port eith empty string.
	port, err := spireg.Open("/dev/spidev0.0")
	if err != nil {
		pterm.Error.Printf("cannot open the spi interface %s\n", err)

		return
	}

	var (
		ResetPin gpio.PinOut = rpi.P1_13
		IRQPin   gpio.PinIn  = rpi.P1_11
	)

	rfid, err := mfrc522.NewSPI(port, ResetPin, IRQPin, mfrc522.WithSync())
	if err != nil {
		pterm.Error.Printf("failed to create mfrc522 device based on spi %s\n", err)

		return
	}

	// setting the antenna signal strength, signal strength from 0 to 7
	if err := rfid.SetAntennaGain(Gain); err != nil {
		pterm.Error.Printf("antenna gain setup failed %s\n", err)

		return
	}

	pterm.Info.Println("Started rfid reader.")

	for {
		// trying to read UID
		data, err := rfid.ReadUID(ReadTimeout)
		if err != nil {
			// here we ignore the reader error because of its many failures.
			// pterm.Error.Printf("cannot read the rfid %s\n", err)
			_ = err
		} else {
			pterm.Info.Println(hex.EncodeToString(data))
		}
	}
}
