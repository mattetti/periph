// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package ht16k33_test

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/experimental/devices/ht16k33"
)

func Example() {
	// Open the I²C bus to which the ht16k33 is connected.
	i2cBus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer i2cBus.Close()

	// We will configure the ht16k33 by setting some options, we can start by the
	// defaults.
	opts := ht16k33.DefaultOpts()

	// Open the device so we can detect touch events.
	dev, err := ht16k33.NewI2C(i2cBus, opts)
	if err != nil {
		log.Fatalf("couldn't open ht16k33: %v", err)
	}

	fmt.Println("Turn on all LEDs")
	for i := 0; i < 16; i++ {
		dev.SetLED(i, true)
		if err := dev.WriteDisplay(); err != nil {
			panic(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(500 * time.Millisecond)
	fmt.Println("And turn them back off")
	for i := 15; i >= 0; i-- {
		dev.SetLED(i, false)
		if err := dev.WriteDisplay(); err != nil {
			panic(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Print("\n")
}
