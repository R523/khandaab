package servo

import (
	"fmt"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
)

type Servo struct {
	DutyNumerator   gpio.Duty
	DutyDenominator gpio.Duty
	Frequency       physic.Frequency

	Pin gpio.PinOut
}

func New(pin gpio.PinOut, dutyNumerator, dutyDenominator gpio.Duty, freq physic.Frequency) Servo {
	return Servo{
		DutyNumerator:   dutyNumerator,
		DutyDenominator: dutyDenominator,

		Pin: pin,

		Frequency: freq,
	}
}

func (s Servo) Start() error {
	if err := s.Pin.PWM((gpio.DutyMax*s.DutyNumerator)/s.DutyDenominator, s.Frequency); err != nil {
		return fmt.Errorf("cannot setup pwm for pin %s %w", s.Pin, err)
	}

	return nil
}

func (s Servo) Stop() error {
	return fmt.Errorf("cannot halt the pin %s %w", s.Pin, s.Pin.Halt())
}
