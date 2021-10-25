package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/pterm/pterm"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/mfrc522"
	"periph.io/x/host/v3"
)

// mfrc522 rfid device
var rfid *mfrc522.Dev

// spi port
var port spi.PortCloser

// pins used for rest and irq
const (
	resetPin = "P1_22" // GPIO 25
	irqPin   = "P1_18" // GPIO 24
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

	// get GPIO rest pin from its name
	var gpioResetPin gpio.PinOut = gpioreg.ByName(resetPin)
	if gpioResetPin == nil {
		log.Fatalf("Failed to find %v", resetPin)
	}

	// get GPIO irq pin from its name
	var gpioIRQPin gpio.PinIn = gpioreg.ByName(irqPin)
	if gpioIRQPin == nil {
		log.Fatalf("Failed to find %v", irqPin)
	}

	rfid, err = mfrc522.NewSPI(port, gpioResetPin, gpioIRQPin, mfrc522.WithSync())
	if err != nil {
		log.Fatal(err)
	}

	// setting the antenna signal strength, signal strength from 0 to 7
	rfid.SetAntennaGain(5)

	fmt.Println("Started rfid reader.")

	// trying to read UID
	data, err := rfid.ReadUID(5 * time.Second)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(hex.EncodeToString(data))
	}
}
