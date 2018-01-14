// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package ht16k33 controls a Microchip ht16k33 device over I²C.
//
// The ht16k33 RAM Mapping 16*8 LED Controller Driver with keyscan.
//
// Datasheet
//
// The official data sheet can be found here:
//
// http://www.holtek.com.tw/productdetail/-/vg/HT16K33
package ht16k33

import (
	"fmt"
	"log"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/devices"
)

// LayoutFunc is a function setting up the board layout converting
// the led index to a message position and bit mask.
type LayoutFunc func(idx int) (byteIDX int, mask byte)

// Dev is a handle to a ht16k33.
type Dev struct {
	c              conn.Conn
	opts           Opts
	displayBuffer  [8]byte
	lastKeysBuffer [6]byte
	keysBuffer     [6]byte
	Layout         LayoutFunc
}

func (d *Dev) String() string {
	return fmt.Sprintf("ht16k33{%s}", d.c)
}

// Halt is a noop for the ht16k33.
func (d *Dev) Halt() error {
	return nil
}

// NewI2C returns a new device that communicates over I²C to ht16k33.
//
// Use default options if nil is used.
func NewI2C(b i2c.Bus, opts *Opts) (*Dev, error) {
	if opts == nil {
		opts = DefaultOpts()
	}
	addr, err := opts.i2cAddr()
	if err != nil {
		return nil, wrapf("%v", err)
	}
	d, err := makeDev(&i2c.Dev{Bus: b, Addr: addr}, opts)
	if err != nil {
		return nil, err
	}
	if opts.Debug {
		log.Printf("ht16k33: Connecting to address: %#x\n", addr)
	}
	return d, d.ClearAll()
}

func makeDev(c conn.Conn, opts *Opts) (*Dev, error) {
	d := &Dev{
		c:    c,
		opts: *opts,
	}
	if d.opts.LEDMap == nil {
		d.opts.LEDMap = AdafruitTrellisLEDLayout
	}
	if d.opts.ButtonMap == nil {
		d.opts.ButtonMap = AdafruitTrellisButtonLayout
	}

	// Turn on the oscillator
	if err := d.c.Tx([]byte{0x21}, nil); err != nil {
		return nil, wrapf("failed to turn on the oscillator: %v", err)
	}

	// Turn on the display but set the blinking off
	if err := d.c.Tx([]byte{0x80 | 0x01 | 0x00}, nil); err != nil {
		return nil, wrapf("failed to turn on the blinking: %v", err)
	}

	// Set brightness to the maximum
	if err := d.c.Tx([]byte{0xE0 | 15}, nil); err != nil {
		return nil, wrapf("failed to set the brightness to the max: %v", err)
	}

	// // Turn on interrupt
	if err := d.c.Tx([]byte{0xA1}, nil); err != nil {
		return nil, wrapf("failed to turn on interrupt: %v", err)
	}

	return d, nil
}

// SetLED sets one of the 128n LEDs as on or off.
// WriteDisplay needs to be called for this change to make effect.
func (d *Dev) SetLED(idx int, on bool) {
	pos, mask := d.opts.LEDMap(idx)
	if on {
		d.displayBuffer[pos] |= mask & 0xff
	} else {
		d.displayBuffer[pos] &= ^mask
	}
}

// ClearAll resets turns off all the LEDs
func (d *Dev) ClearAll() error {
	d.displayBuffer = [8]byte{}
	return d.WriteDisplay()
}

// WriteDisplay applies the LED changes.
func (d *Dev) WriteDisplay() error {
	msg := append([]byte{0x00}, d.displayBuffer[:]...)

	if err := d.c.Tx(msg, nil); err != nil {
		return wrapf("failed to transmit display state: %v", err)
	}
	return nil
}

// IsPressed returns true if the button at the passed index was just pressed (since last read)
func (d *Dev) IsPressed(idx int) bool {
	if idx > 15 || d.opts.ButtonMap == nil {
		return false
	}
	pos, mask := d.opts.ButtonMap(idx)
	return d.keysBuffer[pos]&mask > 0
}

// WasPressed returns true if the button at the passed index was pressed.
func (d *Dev) WasPressed(idx int) bool {
	if idx > 15 || d.opts.ButtonMap == nil {
		return false
	}
	pos, mask := d.opts.ButtonMap(idx)
	return d.lastKeysBuffer[pos]&mask > 0
}

// WasJustReleased returns true if the button at the passed index was pressed and is
// now released.
func (d *Dev) WasJustReleased(idx int) bool {
	return d.WasPressed(idx) && !d.IsPressed(idx)
}

// ReadKeys gets the raw data for the keys buffer (and resets it).
func (d *Dev) ReadKeys() error {
	copy(d.lastKeysBuffer[:], d.keysBuffer[:])
	return d.c.Tx([]byte{0x40}, d.keysBuffer[:])
}

func wrapf(format string, a ...interface{}) error {
	return fmt.Errorf("ht16k33: "+format, a...)
}

var _ devices.Device = &Dev{}
