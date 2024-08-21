package main

import (
	"time"
)

func main() {
	// as recommended by the AVR docs for the ATMega2560(and likely other docs), immediately ensure watchdog is off so that we don't get stuck in a reset loop
	EnsureWatchdogDisabled()
	println("(re)starting...") // I like having text feedback upon reset, it's also helpful for knowing when it's safe to start typing, comment out if you want

	var ctx = CreateContext()

	//InitMethods()

	for {
		c, err := ReadByte()
		if err == nil {
			if c == 13 { // Intercept newline/enter
				print("\n")
			} else if c == 18 { // Intercept ^R, runs the program currently in the context
				ctx.ParseBuffer()
				ctx.Clean()
				println("Running...")
				if ctx.IsInvalid() {
					continue
				}
				Run(ctx.RootNode)
			} else if c == 24 { // Intercept ^X, restarts the entire device/program
				Reset()
			} else if c == 23 { // Intercept ^W, reinitializes the context, should be cleaned up by the garbage collector
				println("Resetting context...")
				ctx = CreateContext()
			} else if c == 127 { // Intercept backspace, unreads the last char
				ctx.Unread()
				WriteByte(c)
			} else if c < 32 {
				// Convert nonprintable control characters to
				// ^A, ^B, etc if they weren't intercepted by the above code
				// Note that I did not write this (imo) jank code! this is from the tinygo serial tutorial!
				WriteByte('^')
				WriteByte(c + '@')
			} else if c >= 127 {
				// Anything equal or above ASCII 127, print ^?.
				WriteByte('^')
				WriteByte('?')
			} else {
				// Echo the printable character back to the
				// host computer.
				// this does not make it repeat characters twice, as by default the characters just vanish into a black hole and aren't written.
				WriteByte(c)
				ctx.ReadChar(c)
			}
		}

		time.Sleep(time.Millisecond * 1) // I found through trial and error that this seems to work well for pasting, adjust as needed
		// You could also replace it with <=100us for pretty much the maximum transfer rate possible, under roughly 100us it doesn't make much of a difference
		// if you are using 9600 baud or most AVR microcontrollers, including MOST arduino boards(particularly the Uno, Mega, Nano, Leonardo, and most except
		// for the more modern ones with ARM microcontrollers, and even many of those use 9600 baud). most others use 115200 and for those use <=8us or so.
	}
}
