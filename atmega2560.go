//go:build atmega2560

// This builds for the ATMega2560(which is also inside the Arduino Mega 2560).

package main

import (
	"device/avr"
	"machine"
	"runtime/interrupt"
	"time"
)

// Reads a single byte and returns it. If nothing is available, returns an error.
func ReadByte() (byte, error) {
	return machine.Serial.ReadByte()
}

// Writes a single byte. It should block until it is finished writing.
func WriteByte(b byte) error {
	return machine.Serial.WriteByte(b)
}

// Reset the device, or at least the program.
func Reset() {
	var state = interrupt.Disable()
	avr.Asm("wdr")
	avr.WDTCSR.Set(avr.WDTCSR.Get() | avr.WDTCSR_WDCE | avr.WDTCSR_WDE)
	avr.WDTCSR.Set(avr.WDTCSR_WDE)
	interrupt.Restore(state)
	for { // Watchdog has just been enabled for ~20ms, sleep and do nothing until watchdog force-resets the device
		time.Sleep(time.Millisecond * 50)
	}
}

// Ensures the watchdog is disabled. If the device doesn't have a watchdog, should do nothing and immediately return.
func EnsureWatchdogDisabled() {
	var state = interrupt.Disable()
	avr.Asm("wdr")
	avr.MCUSR.Set(avr.MCUSR.Get() & ^uint8(1<<3))
	avr.WDTCSR.Set(avr.WDTCSR.Get() | avr.WDTCSR_WDCE | avr.WDTCSR_WDE)
	avr.WDTCSR.Set(0)
	interrupt.Restore(state)
}
