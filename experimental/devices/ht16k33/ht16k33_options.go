// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package ht16k33

// Opts is options to pass to the constructor.
type Opts struct {
	// Debug turns on extra logging capabilities.
	Debug bool
	// I2CAddr is the I²C slave address to use. It can only used on creation of
	// an I²C-device. Its default value is 0x70. It can be set to other values
	I2CAddr uint16
	// Fn is the layout function converting a LED index into a message
	// pos/bitmap
	Fn LayoutFunc
}

func (o *Opts) i2cAddr() (uint16, error) {
	if o.I2CAddr == 0 {
		// Default address.
		return 0x70, nil
	}
	return o.I2CAddr, nil
}

// DefaultOpts returns a pointer to a new Opts with the default option values.
func DefaultOpts() *Opts {
	return &Opts{
		I2CAddr: 0x70,
		Fn:      AdafruitTrellisLayout,
	}
}
