package rfid

import (
	"encoding/hex"
	"fmt"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/mfrc522"
)

func Setup(spiDev string, resetPin gpio.PinOut, irqPin gpio.PinIn, gain int) (*mfrc522.Dev, error) {
	port, err := spireg.Open(spiDev)
	if err != nil {
		return nil, fmt.Errorf("cannot open the spi interface %w", err)
	}

	rfid, err := mfrc522.NewSPI(port, resetPin, irqPin, mfrc522.WithSync())
	if err != nil {
		return nil, fmt.Errorf("failed to create mfrc522 device based on spi %w", err)
	}

	// setting the antenna signal strength, signal strength from 0 to 7
	if err := rfid.SetAntennaGain(gain); err != nil {
		return nil, fmt.Errorf("antenna gain setup failed %w", err)
	}

	return rfid, nil
}

func ReadRFIDWithRetries(rfid *mfrc522.Dev) string {
	for {
		// trying to read UID
		data, err := rfid.ReadUID(-1)
		if err != nil {
			// here we ignore the reader error because of its many failures.
			// pterm.Error.Printf("cannot read the rfid %s\n", err)
			_ = err
		} else {
			return hex.EncodeToString(data)
		}
	}
}
