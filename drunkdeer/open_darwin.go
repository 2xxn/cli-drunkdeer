package main

import (
	"github.com/sstallion/go-hid"
)

// configureOpenMode disables exclusive (seize) opens on macOS.
//
// hidapi defaults to kIOHIDOptionsTypeSeizeDevice for backward compatibility
// (see hid_init in hid_darwin.c). Seizing takes the device away from the system
// driver, which IOKit only permits for root, so IOHIDDeviceOpen fails with
// kIOReturnNotPrivileged (0xE00002C1) for a normal user.
//
// On macOS every top-level collection of the keyboard, including the
// vendor-defined 0xFF00 one this driver talks to, resolves to a single
// IOHIDDevice, so there is no non-keyboard interface to open instead. We only
// send configuration packets and read an identity response, so a shared open is
// enough.
//
// Must be called after hid.Init: hid_init resets this to exclusive.
func configureOpenMode() {
	hid.SetOpenExclusive(false)
}
