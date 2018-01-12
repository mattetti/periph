package ht16k33

// AdafruitTrellisLayout is the layout for the Adafruit 4x4 trellis.
// https://learn.adafruit.com/adafruit-trellis-diy-open-source-led-keypad/downloads
func AdafruitTrellisLayout(idx int) (byteIDX int, mask byte) {
	if idx > 15 {
		idx = 0
	}

	switch idx {
	case 0:
		return 7, 1 << 2
	case 1:
		return 6, 1 << 7
	case 2:
		return 6, 1 << 5
	case 3:
		return 6, 1 << 4
	case 4:
		return 5, 1 << 0
	case 5:
		return 5, 1 << 1
	case 6:
		return 4, 1 << 3
	case 7:
		return 4, 1 << 4
	case 8:
		return 2, 1 << 6
	case 9:
		return 3, 1 << 3
	case 10:
		return 2, 1 << 1
	case 11:
		return 2, 1 << 0
	case 12:
		return 1, 1 << 6
	case 13:
		return 1, 1 << 5
	case 14:
		return 1, 1 << 4
	case 15:
		return 0, 1 << 2
	default:
		return 0, 0xff
	}
}